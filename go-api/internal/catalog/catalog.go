package catalog

import (
	"context"
	"errors"
	"strings"
	"time"

	"{{MODULE}}/internal/id"
	"{{MODULE}}/internal/money"
)

var (
	ErrNotFound = errors.New("product not found")
	ErrInvalid  = errors.New("invalid product")
	ErrConflict = errors.New("product conflict")
)

type Product struct {
	ID          string
	SKU         string
	Name        string
	Description string
	Price       money.Amount
	ImageURL    string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Repository interface {
	Create(ctx context.Context, p Product) error
	Get(ctx context.Context, id string) (Product, error)
	GetBySKU(ctx context.Context, sku string) (Product, error)
	List(ctx context.Context) ([]Product, error)
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository, now func() time.Time) *Service {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Service{repo: repo, now: now}
}

type CreateInput struct {
	SKU         string
	Name        string
	Description string
	PriceMinor  int64
	Currency    string
	ImageURL    string
}

func (s *Service) Create(ctx context.Context, in CreateInput) (Product, error) {
	sku := strings.ToUpper(strings.TrimSpace(in.SKU))
	name := strings.TrimSpace(in.Name)
	if sku == "" || name == "" {
		return Product{}, ErrInvalid
	}
	price, err := money.New(in.PriceMinor, in.Currency)
	if err != nil {
		return Product{}, err
	}
	if _, err := s.repo.GetBySKU(ctx, sku); err == nil {
		return Product{}, ErrConflict
	} else if err != ErrNotFound {
		return Product{}, err
	}
	now := s.now()
	p := Product{
		ID:          id.New(),
		SKU:         sku,
		Name:        name,
		Description: strings.TrimSpace(in.Description),
		Price:       price,
		ImageURL:    strings.TrimSpace(in.ImageURL),
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return Product{}, err
	}
	return p, nil
}

func (s *Service) Get(ctx context.Context, productID string) (Product, error) {
	if err := id.Parse(productID); err != nil {
		return Product{}, ErrInvalid
	}
	return s.repo.Get(ctx, productID)
}

func (s *Service) GetBySKU(ctx context.Context, sku string) (Product, error) {
	sku = strings.ToUpper(strings.TrimSpace(sku))
	if sku == "" {
		return Product{}, ErrInvalid
	}
	return s.repo.GetBySKU(ctx, sku)
}

func (s *Service) List(ctx context.Context) ([]Product, error) {
	return s.repo.List(ctx)
}

func Seed(ctx context.Context, cat *Service) error {
	items := []CreateInput{
		{SKU: "MUG-1", Name: "Demo Mug", Description: "Adapted from P06 seed.", PriceMinor: 1200, Currency: "JPY", ImageURL: "https://placehold.co/400x400?text=Mug"},
		{SKU: "TEE-1", Name: "Demo Tee", Description: "Plenty in stock in the full commerce app.", PriceMinor: 3500, Currency: "JPY", ImageURL: "https://placehold.co/400x400?text=Tee"},
		{SKU: "STK-1", Name: "Demo Sticker", Description: "Integer money, no floats.", PriceMinor: 300, Currency: "JPY", ImageURL: "https://placehold.co/400x400?text=Sticker"},
	}
	for _, in := range items {
		if _, err := cat.GetBySKU(ctx, in.SKU); err == nil {
			continue
		} else if err != ErrNotFound {
			return err
		}
		if _, err := cat.Create(ctx, in); err != nil {
			return err
		}
	}
	return nil
}
