package inventory

import (
	"context"
	"errors"
	domaininventory "go-erp/internal/domain/inventory"
	dtoinventory "go-erp/internal/dto/inventory"
	inventoryrepo "go-erp/internal/repository/inventory"
	"go-erp/pkg/dbtx"
	"go-erp/pkg/errs"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const stockDeductLua = `
local current = redis.call("GET", KEYS[1])
if not current then
  return -2
end
current = tonumber(current)
local deduct = tonumber(ARGV[1])
if current < deduct then
  return -1
end
return redis.call("DECRBY", KEYS[1], deduct)
`

type Service struct {
	db   *gorm.DB
	repo inventoryrepo.Repository
	rdb  *redis.Client
}

// NewService 创建库存服务。
func NewService(db *gorm.DB, repo inventoryrepo.Repository, rdb *redis.Client) *Service {
	return &Service{db: db, repo: repo, rdb: rdb}
}

// Deduct 扣减库存（自动开启事务）。
func (s *Service) Deduct(ctx context.Context, req dtoinventory.DeductRequest) error {
	return dbtx.WithTransaction(ctx, s.db, func(tx *gorm.DB) error {
		return s.DeductWithTx(ctx, tx, req)
	})
}

// DeductWithTx 在给定事务中扣减库存。
// 备注：先走 Redis Lua 原子预扣，再落库，失败会尝试回补 Redis。
func (s *Service) DeductWithTx(ctx context.Context, tx *gorm.DB, req dtoinventory.DeductRequest) error {
	cacheKey := "stock:sku:" + strconv.FormatUint(uint64(req.SKUID), 10)
	if s.rdb != nil {
		stock, err := s.rdb.Eval(ctx, stockDeductLua, []string{cacheKey}, req.Qty).Int64()
		if err != nil {
			return errs.Wrap(errs.CodeInternal, "redis_deduct_failed", err)
		}
		if stock == -1 {
			return errs.New(errs.CodeInsufficientSKU, "insufficient_stock")
		}
		if stock == -2 {
			dbInv, err := s.repo.GetBySKUID(ctx, req.SKUID)
			if err != nil {
				return errs.New(errs.CodeNotFound, "inventory_not_found")
			}
			if dbInv.Stock < req.Qty {
				return errs.New(errs.CodeInsufficientSKU, "insufficient_stock")
			}
			if err := s.rdb.Set(ctx, cacheKey, strconv.FormatInt(dbInv.Stock, 10), 0).Err(); err != nil {
				return errs.Wrap(errs.CodeInternal, "redis_warmup_failed", err)
			}
			stock, err = s.rdb.Eval(ctx, stockDeductLua, []string{cacheKey}, req.Qty).Int64()
			if err != nil {
				return errs.Wrap(errs.CodeInternal, "redis_deduct_failed", err)
			}
			if stock < 0 {
				return errs.New(errs.CodeInsufficientSKU, "insufficient_stock")
			}
		}
	}
	deductErr := s.repo.Deduct(ctx, tx, req.SKUID, req.Qty, req.BusinessNo)
	if deductErr == nil {
		return nil
	}
	if s.rdb != nil {
		_ = s.rdb.IncrBy(ctx, cacheKey, req.Qty).Err()
	}
	if errors.Is(deductErr, gorm.ErrRecordNotFound) {
		return errs.New(errs.CodeNotFound, "inventory_not_found")
	}
	return errs.Wrap(errs.CodeInternal, "db_deduct_failed", deductErr)
}

// InitStock 初始化或更新 SKU 库存。
func (s *Service) InitStock(ctx context.Context, skuID uint, qty int64) error {
	row := &domaininventory.Inventory{
		SKUID: skuID,
		Stock: qty,
	}
	return s.db.WithContext(ctx).Where("sku_id = ?", skuID).Assign(row).FirstOrCreate(row).Error
}
