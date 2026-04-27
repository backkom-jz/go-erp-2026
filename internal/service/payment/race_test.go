package payment

import (
	"context"
	"fmt"
	domainorder "go-erp/internal/domain/order"
	domainpayment "go-erp/internal/domain/payment"
	dtopayment "go-erp/internal/dto/payment"
	orderrepo "go-erp/internal/repository/order"
	paymentrepo "go-erp/internal/repository/payment"
	ordersvc "go-erp/internal/service/order"
	"sync"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPaymentCallbackRaceWithTimeoutCancel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&domainorder.Order{}, &domainorder.OrderItem{}, &domainpayment.Record{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	orderRepository := orderrepo.NewRepository(db)
	paymentRepository := paymentrepo.NewRepository(db)
	orderService := ordersvc.NewService(db, orderRepository, orderRepository, nil)
	paymentService := NewService(paymentRepository, orderService, nil)

	for i := 0; i < 50; i++ {
		orderNo := fmt.Sprintf("o-race-%d", i)
		row := domainorder.Order{
			OrderNo:    orderNo,
			UserID:     1,
			TenantID:   "t1",
			Status:     domainorder.StatusPending,
			TotalCents: 100,
		}
		if err := db.Create(&row).Error; err != nil {
			t.Fatalf("create order failed: %v", err)
		}

		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(2)

		go func(orderNo string) {
			defer wg.Done()
			<-start
			_ = orderService.CancelIfTimeout(context.Background(), orderNo)
		}(orderNo)

		go func(orderNo string, idx int) {
			defer wg.Done()
			<-start
			_ = paymentService.Callback(context.Background(), dtopayment.CallbackRequest{
				OrderNo:   orderNo,
				PaymentNo: fmt.Sprintf("p-race-%d", idx),
				Channel:   "mock",
				Status:    "paid",
			})
		}(orderNo, i)

		close(start)
		wg.Wait()

		var latest domainorder.Order
		if err := db.Where("order_no = ?", orderNo).First(&latest).Error; err != nil {
			t.Fatalf("query order failed: %v", err)
		}
		if latest.Status != domainorder.StatusPaid {
			t.Fatalf("order should be paid after race, got %s", latest.Status)
		}
	}
}
