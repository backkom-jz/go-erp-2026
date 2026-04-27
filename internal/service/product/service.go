package product

import (
	"context"
	"go-erp/internal/domain/product"
	dtoproduct "go-erp/internal/dto/product"
	productrepo "go-erp/internal/repository/product"
	"go-erp/pkg/errs"
)

type Service struct {
	repo productrepo.Repository
}

func NewService(repo productrepo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSPU(ctx context.Context, req dtoproduct.CreateSPURequest) (*product.SPU, error) {
	row := &product.SPU{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Brand:      req.Brand,
	}
	if err := s.repo.CreateSPU(ctx, row); err != nil {
		return nil, errs.Wrap(errs.CodeInternal, "create_spu_failed", err)
	}
	return row, nil
}

func (s *Service) ListSPU(ctx context.Context, limit int) ([]product.SPU, error) {
	return s.repo.ListSPU(ctx, limit)
}

func (s *Service) CreateSKU(ctx context.Context, req dtoproduct.CreateSKURequest) (*product.SKU, error) {
	row := &product.SKU{
		SPUID:      req.SPUID,
		Code:       req.Code,
		Name:       req.Name,
		PriceCents: req.PriceCents,
	}
	if err := s.repo.CreateSKU(ctx, row); err != nil {
		return nil, errs.Wrap(errs.CodeInternal, "create_sku_failed", err)
	}
	return row, nil
}
