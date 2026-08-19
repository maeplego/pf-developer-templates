package catalog_test

import (
	"context"
	"testing"
	"time"

	"{{MODULE}}/api/internal/catalog"
	"{{MODULE}}/api/internal/store/memory"
)

func TestCreateAndGetBySKU(t *testing.T) {
	ctx := context.Background()
	st := memory.New()
	svc := catalog.NewService(memory.Catalog{Store: st}, func() time.Time { return time.Date(2026, 8, 19, 4, 0, 0, 0, time.UTC) })

	p, err := svc.Create(ctx, catalog.CreateInput{
		SKU: "mug-1", Name: "Demo Mug", PriceMinor: 1200, Currency: "JPY",
		ImageURL: "https://example.invalid/mug.png",
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := svc.GetBySKU(ctx, "MUG-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != p.ID || got.Price.Minor != 1200 || got.ImageURL == "" {
		t.Fatalf("%+v", got)
	}
	if _, err := svc.Create(ctx, catalog.CreateInput{SKU: "MUG-1", Name: "Dup", PriceMinor: 1, Currency: "JPY"}); err != catalog.ErrConflict {
		t.Fatalf("dup sku: %v", err)
	}
}

func TestSeedIdempotent(t *testing.T) {
	ctx := context.Background()
	st := memory.New()
	svc := catalog.NewService(memory.Catalog{Store: st}, time.Now)
	if err := catalog.Seed(ctx, svc); err != nil {
		t.Fatal(err)
	}
	if err := catalog.Seed(ctx, svc); err != nil {
		t.Fatal(err)
	}
	list, err := svc.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 {
		t.Fatalf("%d", len(list))
	}
}
