package product

import (
	"context"
	"errors"

	"github.com/vincentconace/api-gin/internal/domain"
)

type Service interface {
	Get(ctx context.Context) ([]domain.Product, error)
	GetById(ctx context.Context, id int) (domain.Product, error)
	Create(ctx context.Context, p domain.Product) (domain.Product, error)
	Update(ctx context.Context, id int, p domain.Product) (domain.Product, error)
	Delete(ctx context.Context, id int) error
}

var (
	EmptyProduct          = domain.Product{}
	ErrNotFound           = errors.New("product not found")
	ErrInternal           = errors.New("internal error")
	ErrProductAlredyExist = errors.New("product already exists")
	ErrCreatedProduct     = errors.New("error creating product")
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context) ([]domain.Product, error) {
	return s.repo.Get(ctx)
}

func (s *service) GetById(ctx context.Context, id int) (domain.Product, error) {
	return s.repo.GetById(ctx, id)
}

func (s *service) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	result := s.repo.Exists(ctx, *p.ProductCode)
	if result {
		return EmptyProduct, ErrProductAlredyExist
	}
	id, err := s.repo.Save(ctx, p)
	if err != nil {
		return EmptyProduct, ErrCreatedProduct
	}
	p.ID = &id
	return p, nil
}

func (s *service) Update(ctx context.Context, id int, p domain.Product) (domain.Product, error) {
	result := s.repo.Exists(ctx, *p.ProductCode)
	if result {
		return EmptyProduct, ErrProductAlredyExist
	}
	persistendProduct, err := s.repo.GetById(ctx, id)
	if err != nil {
		return EmptyProduct, err
	}
	err = s.repo.Update(ctx, id, p)
	if err != nil {
		return EmptyProduct, err
	}

	persistendProduct.ProductCode = p.ProductCode
	persistendProduct.Name = p.Name
	persistendProduct.Description = p.Description
	persistendProduct.Price = p.Price
	persistendProduct.Stock = p.Stock

	return persistendProduct, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
