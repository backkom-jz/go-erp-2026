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

	testCases := []struct {
		name       string
		firstStep  func(orderNo string) error
		secondStep func(orderNo string) error
	}{
		{
			name: "pay_then_timeout_cancel",
			firstStep: func(orderNo string) error {
				return paymentService.Callback(context.Background(), dtopayment.CallbackRequest{
					OrderNo:   orderNo,
					PaymentNo: "p-first-" + orderNo,
					Channel:   "mock",
					Status:    "paid",
				})
			},
			secondStep: func(orderNo string) error {
				return orderService.CancelIfTimeout(context.Background(), orderNo)
			},
		},
		{
			name: "timeout_cancel_then_pay",
			firstStep: func(orderNo string) error {
				return orderService.CancelIfTimeout(context.Background(), orderNo)
			},
			secondStep: func(orderNo string) error {
				return paymentService.Callback(context.Background(), dtopayment.CallbackRequest{
					OrderNo:   orderNo,
					PaymentNo: "p-second-" + orderNo,
					Channel:   "mock",
					Status:    "paid",
				})
			},
		},
	}

	for i, tc := range testCases {
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
		if err := tc.firstStep(orderNo); err != nil {
			t.Fatalf("%s first step failed: %v", tc.name, err)
		}
		if err := tc.secondStep(orderNo); err != nil {
			t.Fatalf("%s second step failed: %v", tc.name, err)
		}
		var latest domainorder.Order
		if err := db.Where("order_no = ?", orderNo).First(&latest).Error; err != nil {
			t.Fatalf("query order failed: %v", err)
		}
		if latest.Status != domainorder.StatusPaid {
			t.Fatalf("%s should end paid, got %s", tc.name, latest.Status)
		}
	}
}
