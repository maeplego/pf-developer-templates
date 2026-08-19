package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"{{MODULE}}/api/internal/auth"
	"{{MODULE}}/api/internal/catalog"
	"{{MODULE}}/api/internal/clock"
	"{{MODULE}}/api/internal/store/memory"
	"{{MODULE}}/api/internal/web"
)

func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	ctx := context.Background()
	st := memory.New()
	clk := &clock.Fixed{T: time.Date(2026, 8, 19, 4, 0, 0, 0, time.UTC)}
	cat := catalog.NewService(memory.Catalog{Store: st}, clk.Now)
	if err := catalog.Seed(ctx, cat); err != nil {
		t.Fatal(err)
	}
	mw := auth.New(true, "", "", "")
	srv := web.New(cat, "", mw, func() error { return st.Ping(context.Background()) })
	return httptest.NewServer(srv.Routes())
}

func TestHealthReady(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	for _, path := range []string{"/health", "/ready"} {
		res, err := http.Get(ts.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 200 {
			t.Fatalf("%s %d", path, res.StatusCode)
		}
		_ = res.Body.Close()
	}
}

func TestListProductsPublic(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/products")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("%d", res.StatusCode)
	}
	var body struct {
		Products []struct {
			SKU        string `json:"sku"`
			PriceMinor int64  `json:"priceMinor"`
		} `json:"products"`
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Products) != 3 {
		t.Fatalf("%d products", len(body.Products))
	}
	for _, p := range body.Products {
		if p.PriceMinor < 0 {
			t.Fatal("negative money")
		}
	}
}

func TestOpsCreateRequiresAuth(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(map[string]any{
		"sku": "NEW-1", "name": "New", "priceMinor": 100, "currency": "JPY",
	})
	res, err := http.Post(ts.URL+"/v1/ops/products", "application/json", &buf)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", res.StatusCode)
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/ops/products", bytes.NewBufferString(`{"sku":"NEW-1","name":"New","priceMinor":100,"currency":"JPY"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Dev-User-Sub", "ops-1")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("want 201 got %d", res.StatusCode)
	}
}
