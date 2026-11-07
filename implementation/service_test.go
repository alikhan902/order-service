package implementation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/log"

	"example.com/order-service/implementation"
	"example.com/order-service/order"
	"example.com/order-service/repository"
)

func TestOrderLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := implementation.NewService(repository.NewMemoryRepository(), log.NewNopLogger())
	id, err := svc.Create(ctx, order.Order{CustomerID: "client-7"})
	if err != nil || id == "" {
		t.Fatalf("create: id=%q err=%v", id, err)
	}
	saved, err := svc.GetByID(ctx, id)
	if err != nil || saved.ID != id || saved.Status != "Pending" || saved.CustomerID != "client-7" || saved.CreatedOn == 0 {
		t.Fatalf("read: %+v err=%v", saved, err)
	}
	if err := svc.ChangeStatus(ctx, id, "Ready"); err != nil {
		t.Fatal(err)
	}
	saved, err = svc.GetByID(ctx, id)
	if err != nil || saved.Status != "Ready" {
		t.Fatalf("updated: %+v err=%v", saved, err)
	}
}

func TestMissingOrder(t *testing.T) {
	svc := implementation.NewService(repository.NewMemoryRepository(), log.NewNopLogger())
	ctx := context.Background()
	if _, err := svc.GetByID(ctx, "missing"); !errors.Is(err, order.ErrOrderNotFound) {
		t.Fatalf("read: %v", err)
	}
	if err := svc.ChangeStatus(ctx, "missing", "Ready"); !errors.Is(err, order.ErrCmdRepository) {
		t.Fatalf("update: %v", err)
	}
}

func TestChangeStatusWithEmptyID(t *testing.T) {
	svc := implementation.NewService(repository.NewMemoryRepository(), log.NewNopLogger())
	svc.ChangeStatus(context.Background(), "", "Ready")
}
