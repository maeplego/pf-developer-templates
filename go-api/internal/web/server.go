package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"{{MODULE}}/internal/auth"
	"{{MODULE}}/internal/catalog"
	"{{MODULE}}/internal/money"
)

type Server struct {
	cat   *catalog.Service
	cors  string
	auth  *auth.Middleware
	ready func() error
}

func New(cat *catalog.Service, cors string, mw *auth.Middleware, ready func() error) *Server {
	if ready == nil {
		ready = func() error { return nil }
	}
	return &Server{cat: cat, cors: cors, auth: mw, ready: ready}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, _ *http.Request) {
		if err := s.ready(); err != nil {
			writeError(w, http.StatusServiceUnavailable, "not_ready", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})
	mux.HandleFunc("GET /v1/products", s.listProducts)
	mux.HandleFunc("GET /v1/products/{id}", s.getProduct)
	mux.Handle("POST /v1/ops/products", s.auth.Handler(http.HandlerFunc(s.opsCreateProduct)))
	return s.withCORS(mux)
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cors != "" {
			w.Header().Set("Access-Control-Allow-Origin", s.cors)
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Dev-User-Sub")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type productJSON struct {
	ID          string `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceMinor  int64  `json:"priceMinor"`
	Currency    string `json:"currency"`
	ImageURL    string `json:"imageUrl"`
	Active      bool   `json:"active"`
}

func toProductJSON(p catalog.Product) productJSON {
	return productJSON{
		ID: p.ID, SKU: p.SKU, Name: p.Name, Description: p.Description,
		PriceMinor: p.Price.Minor, Currency: p.Price.Currency, ImageURL: p.ImageURL,
		Active: p.Active,
	}
}

func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	list, err := s.cat.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	out := make([]productJSON, 0, len(list))
	for _, p := range list {
		out = append(out, toProductJSON(p))
	}
	writeJSON(w, http.StatusOK, map[string]any{"products": out})
}

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
	p, err := s.cat.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProductJSON(p))
}

func (s *Server) opsCreateProduct(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.UserFrom(r.Context()); !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}
	var body struct {
		SKU         string `json:"sku"`
		Name        string `json:"name"`
		Description string `json:"description"`
		PriceMinor  int64  `json:"priceMinor"`
		Currency    string `json:"currency"`
		ImageURL    string `json:"imageUrl"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid", "invalid json")
		return
	}
	p, err := s.cat.Create(r.Context(), catalog.CreateInput{
		SKU: body.SKU, Name: body.Name, Description: body.Description,
		PriceMinor: body.PriceMinor, Currency: body.Currency, ImageURL: body.ImageURL,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toProductJSON(p))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, catalog.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "not found")
	case errors.Is(err, catalog.ErrConflict):
		writeError(w, http.StatusConflict, "conflict", err.Error())
	case errors.Is(err, catalog.ErrInvalid), errors.Is(err, money.ErrInvalid), errors.Is(err, money.ErrCurrency):
		writeError(w, http.StatusBadRequest, "invalid", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "error", "internal error")
	}
}
