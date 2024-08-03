package main

import (
	"context"
	"errors"
	"fmt"
	"journey/internal/api"
	"journey/internal/api/spec"
	"journey/internal/mailer/mailpit"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phenpessoa/gutils/netutils/httputils"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("goodbye :)")
}

func run(ctx context.Context) error {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	logger = logger.Named("journey_app")
	defer logger.Sync()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("JOURNEY_DATABASE_USER"),
		os.Getenv("JOURNEY_DATABASE_PASSWORD"),
		os.Getenv("JOURNEY_DATABASE_HOST"),
		os.Getenv("JOURNEY_DATABASE_PORT"),
		os.Getenv("JOURNEY_DATABASE_NAME"),
	))
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	mailer := mailpit.NewMailpit(pool)

	si := api.NewAPI(pool, logger, mailer)

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, httputils.ChiLogger(logger))
	r.Mount("/", spec.Handler(si))

	// Setup Swagger UI
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "../../internal/api/spec/journey.spec.json")
    })
	r.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("/swagger.json"), // The url pointing to API definition
    ))
	// Setup Scalar docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "../../internal/api/spec/journey.spec.json", 
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})


	srv := &http.Server{
		Addr: ":8080",
		Handler: r,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,	
	}

	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown server", zap.Error(err))
		}
	}()

	errChan := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
		
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}