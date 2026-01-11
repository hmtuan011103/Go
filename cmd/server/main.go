package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/gostructure/app/internal/adapter/handler/http"
	"github.com/gostructure/app/internal/adapter/storage/mysql"
	"github.com/gostructure/app/internal/config"
)

func main() {
	// =====================================================
	// 1. Load App / Server / JWT config
	// =====================================================
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	// Set timezone
	if loc, err := time.LoadLocation(cfg.App.Timezone); err == nil {
		time.Local = loc
	}

	// =====================================================
	// 2. Load Database config (SEPARATE)
	// =====================================================
	dbCfg, err := config.LoadDatabaseOnly()
	if err != nil {
		log.Fatalf("load database config failed: %v", err)
	}

	// =====================================================
	// 3. Connect Database
	// =====================================================
	db, err := mysql.NewMySQLConnection(dbCfg, cfg.App.Timezone)
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}
	defer db.Close()

	// =====================================================
	// 4. Init HTTP Server (ALL wiring inside)
	// =====================================================
	server, err := httpadapter.NewServer(cfg, db)
	if err != nil {
		log.Fatalf("init http server failed: %v", err)
	}

	httpServer := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      server.Handler(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// =====================================================
	// 5. Start HTTP server
	// =====================================================
	go func() {
		log.Printf("HTTP server listening on %s", cfg.Server.Address)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// =====================================================
	// 6. Graceful shutdown
	// =====================================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped cleanly")
}
