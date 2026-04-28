package bootstrap

import (
	"context"
	aihandler "go-erp/internal/handler/http/ai"
	authhandler "go-erp/internal/handler/http/auth"
	inventoryhandler "go-erp/internal/handler/http/inventory"
	orderhandler "go-erp/internal/handler/http/order"
	paymenthandler "go-erp/internal/handler/http/payment"
	producthandler "go-erp/internal/handler/http/product"
	userhandler "go-erp/internal/handler/http/user"
	inventoryrepo "go-erp/internal/repository/inventory"
	orderrepo "go-erp/internal/repository/order"
	paymentrepo "go-erp/internal/repository/payment"
	productrepo "go-erp/internal/repository/product"
	userrepo "go-erp/internal/repository/user"
	authsvc "go-erp/internal/service/auth"
	inventorysvc "go-erp/internal/service/inventory"
	ordersvc "go-erp/internal/service/order"
	paymentsvc "go-erp/internal/service/payment"
	productsvc "go-erp/internal/service/product"
	usersvc "go-erp/internal/service/user"
	aisvc "go-erp/internal/service/ai"
	"go-erp/pkg/auth/jwt"
	"go-erp/pkg/idempotency"
	"go-erp/pkg/mq"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	AIHandler        *aihandler.Handler
	AuthHandler      *authhandler.Handler
	UserHandler      *userhandler.Handler
	ProductHandler   *producthandler.Handler
	InventoryHandler *inventoryhandler.Handler
	OrderHandler     *orderhandler.Handler
	PaymentHandler   *paymenthandler.Handler

	JWTManager       *jwt.Manager
	IdempotencyStore *idempotency.Store
	OrderService     *ordersvc.Service
	OutboxDispatcher *ordersvc.OutboxDispatcher
}

func BuildApp(cfg *Config, db *gorm.DB, redisClient *redis.Client, publisher mq.Publisher) *App {
	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTLMinutes, cfg.JWT.RefreshTTLMinutes)
	var idemStore *idempotency.Store
	if redisClient != nil {
		idemStore = idempotency.NewStore(redisClient)
	}

	userRepository := userrepo.NewRepository(db)
	productRepository := productrepo.NewRepository(db)
	inventoryRepository := inventoryrepo.NewRepository(db)
	orderRepository := orderrepo.NewRepository(db)
	paymentRepository := paymentrepo.NewRepository(db)

	userService := usersvc.NewService(userRepository)
	aiService := aisvc.NewService(aisvc.Config{
		Enabled:        cfg.AI.Enabled,
		BaseURL:        cfg.AI.BaseURL,
		APIKey:         cfg.AI.APIKey,
		Model:          cfg.AI.Model,
		TimeoutSeconds: cfg.AI.TimeoutSeconds,
		Temperature:    cfg.AI.Temperature,
		MaxTokens:      cfg.AI.MaxTokens,
	})
	authService := authsvc.NewService(userRepository, jwtManager)
	productService := productsvc.NewService(productRepository)
	inventoryService := inventorysvc.NewService(db, inventoryRepository, redisClient)
	orderService := ordersvc.NewService(db, orderRepository, orderRepository, inventoryService)
	paymentService := paymentsvc.NewService(paymentRepository, orderService, idemStore)
	outboxDispatcher := ordersvc.NewOutboxDispatcher(
		orderRepository,
		publisher,
		2*time.Second,
		cfg.MQ.OutboxMaxRetry,
		time.Duration(cfg.MQ.OutboxBaseDelaySeconds)*time.Second,
	)

	return &App{
		AIHandler:        aihandler.NewHandler(aiService),
		AuthHandler:      authhandler.NewHandler(authService),
		UserHandler:      userhandler.NewHandler(userService),
		ProductHandler:   producthandler.NewHandler(productService),
		InventoryHandler: inventoryhandler.NewHandler(inventoryService),
		OrderHandler:     orderhandler.NewHandler(orderService),
		PaymentHandler:   paymenthandler.NewHandler(paymentService),
		JWTManager:       jwtManager,
		IdempotencyStore: idemStore,
		OrderService:     orderService,
		OutboxDispatcher: outboxDispatcher,
	}
}

func (a *App) StartBackgroundWorkers(ctx context.Context, mqClient *MQClient, logger *zap.Logger) {
	if a.OutboxDispatcher != nil {
		go a.OutboxDispatcher.Run(ctx)
	}
	if mqClient == nil || a.OrderService == nil {
		return
	}

	msgs, err := mqClient.Channel.Consume(
		orderTimeoutProcessQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("start timeout consumer failed", zap.Error(err))
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				if err := a.OrderService.HandleTimeoutMessage(ctx, msg.Body); err != nil {
					logger.Error("handle timeout message failed", zap.Error(err))
					_ = msg.Nack(false, true)
					continue
				}
				_ = msg.Ack(false)
			}
		}
	}()
}
