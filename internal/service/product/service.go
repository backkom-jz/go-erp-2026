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

// NewService 创建商品服务。
func NewService(repo productrepo.Repository) *Service {
	return &Service{repo: repo}
}

// CreateSPU 创建 SPU。
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

// ListSPU 查询 SPU 列表。
func (s *Service) ListSPU(ctx context.Context, limit int) ([]product.SPU, error) {
	return s.repo.ListSPU(ctx, limit)
}

// CreateSKU 创建 SKU。
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
