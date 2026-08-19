package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{MODULE}}/api/internal/auth"
	"{{MODULE}}/api/internal/catalog"
	"{{MODULE}}/api/internal/clock"
	"{{MODULE}}/api/internal/config"
	"{{MODULE}}/api/internal/store/memory"
	"{{MODULE}}/api/internal/web"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.OTelEndpoint != "" {
		log.Printf("otel endpoint reserved: %s (service=%s)", cfg.OTelEndpoint, cfg.ServiceName)
	}

	ctx := context.Background()
	clk := clock.System{}
	st := memory.New()
	cat := catalog.NewService(memory.Catalog{Store: st}, clk.Now)
	if err := catalog.Seed(ctx, cat); err != nil {
		log.Fatal(err)
	}

	mw := auth.New(cfg.DevAuth, cfg.OIDCIssuer, cfg.OIDCInternalBase, cfg.OIDCAudience)
	handler := web.New(cat, cfg.CORSOrigin, mw, func() error { return st.Ping(context.Background()) }).Routes()
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("%s api listening on %s (devAuth=%v)", cfg.ServiceName, cfg.HTTPAddr, cfg.DevAuth)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down")
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutCtx)
}
