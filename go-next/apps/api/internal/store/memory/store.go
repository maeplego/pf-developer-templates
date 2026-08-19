package memory

import (
	"context"
	"sort"
	"sync"

	"{{MODULE}}/api/internal/catalog"
)

type Store struct {
	mu            sync.Mutex
	products      map[string]catalog.Product
	productsBySKU map[string]string
}

func New() *Store {
	return &Store{
		products:      map[string]catalog.Product{},
		productsBySKU: map[string]string{},
	}
}

func (s *Store) Ping(context.Context) error { return nil }

func (s *Store) Create(_ context.Context, p catalog.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.productsBySKU[p.SKU]; ok {
		return catalog.ErrConflict
	}
	s.products[p.ID] = p
	s.productsBySKU[p.SKU] = p.ID
	return nil
}

func (s *Store) Get(_ context.Context, id string) (catalog.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[id]
	if !ok {
		return catalog.Product{}, catalog.ErrNotFound
	}
	return p, nil
}

func (s *Store) GetBySKU(_ context.Context, sku string) (catalog.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.productsBySKU[sku]
	if !ok {
		return catalog.Product{}, catalog.ErrNotFound
	}
	return s.products[id], nil
}

func (s *Store) List(context.Context) ([]catalog.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]catalog.Product, 0, len(s.products))
	for _, p := range s.products {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SKU < out[j].SKU })
	return out, nil
}

type Catalog struct{ *Store }

func (c Catalog) Create(ctx context.Context, p catalog.Product) error { return c.Store.Create(ctx, p) }
func (c Catalog) Get(ctx context.Context, id string) (catalog.Product, error) {
	return c.Store.Get(ctx, id)
}
func (c Catalog) GetBySKU(ctx context.Context, sku string) (catalog.Product, error) {
	return c.Store.GetBySKU(ctx, sku)
}
func (c Catalog) List(ctx context.Context) ([]catalog.Product, error) { return c.Store.List(ctx) }
