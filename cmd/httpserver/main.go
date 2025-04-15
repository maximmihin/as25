package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	slogfiber "github.com/samber/slog-fiber"

	handlers "github.com/maximmihin/as25/internal/controllers/http"
	"github.com/maximmihin/as25/internal/controllers/http/middleware"
	product_repo "github.com/maximmihin/as25/internal/dal/repos/products"
	pvz_repo "github.com/maximmihin/as25/internal/dal/repos/pvz"
	receptions_repo "github.com/maximmihin/as25/internal/dal/repos/receptions"
)

type Config struct {
	DbConnString  string
	JwtPrivateKey string
	jwtPublicKey  string

	Host string
	Port string

	LogLevel string
}

func main() {

	cfg := Config{
		DbConnString:  os.Getenv("POSTGRES_CONN_STRING"),
		JwtPrivateKey: os.Getenv("JWT_SECRET_KEY"),
		jwtPublicKey:  os.Getenv("JWT_SECRET_KEY"),

		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),

		LogLevel: os.Getenv("LOG_LEVEL"),
	}

	var logLevel slog.Level

	if err := logLevel.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).
			Error("invalid logLevel in config",
				slog.String("error", err.Error()))
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	errC, shutDownFunc := Run(cfg, logger)
	defer shutDownFunc(context.TODO())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM) // TODO another signal?

	// TODO обудмать вариант с более красивым gs
	select {
	case err := <-errC:
		logger.Error("fail on Run",
			slog.Any("error", err))
	case <-sig:
	}

}

type ShutDownFunc func(ctx context.Context)

var globalOnListenFunc fiber.OnListenHandler

func Run(cfg Config, logger *slog.Logger) (<-chan error, ShutDownFunc) {

	errC := make(chan error, 1)
	ctx := context.TODO()

	// TODO GS
	dbpool, err := pgxpool.New(ctx, cfg.DbConnString)
	if err != nil {
		errC <- err
	}

	server := handlers.Server{
		JwtSecret:     cfg.JwtPrivateKey,
		PvzRepo:       pvz_repo.New(dbpool),
		ReceptionRepo: receptions_repo.New(dbpool),
		ProductRepo:   product_repo.New(dbpool),
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:          middleware.NewErrorHandler(),
		DisableStartupMessage: true,
	})

	var once sync.Once
	var shutDownFunc = ShutDownFunc(func(ctx context.Context) {
		once.Do(func() {
			err := app.ShutdownWithContext(ctx)
			if err != nil { // todo???
				logger.Error("fail to shutDown",
					slog.String("error", err.Error()))
			}

			dbpool.Close()
		})
	})

	if globalOnListenFunc != nil {
		app.Hooks().OnListen(globalOnListenFunc)
	}

	HealthRegisterHandlers(app, server, logger, cfg.jwtPublicKey)

	go func() {
		errC <- app.Listen(net.JoinHostPort(cfg.Host, cfg.Port))
	}()

	// TODO add gs

	return errC, shutDownFunc
}

func HealthRegisterHandlers(router fiber.Router, si handlers.ServerInterface, logger *slog.Logger, jwtPublicKey string) {

	// common middlewares
	router.Use(slogfiber.NewWithConfig(logger, slogfiber.Config{
		DefaultLevel:      slog.LevelDebug,
		ClientErrorLevel:  slog.LevelInfo,
		ServerErrorLevel:  slog.LevelError,
		WithRequestID:     true,
		WithRequestBody:   true,
		WithRequestHeader: true,
		WithResponseBody:  true}))
	router.Use(recover.New())

	// generated wrapper (boilerplate query params  extractor, usually)
	wrapper := handlers.ServerInterfaceWrapper{
		Handler: si,
	}

	authOnlyModer := middleware.NewJWT(jwtPublicKey, middleware.AccessModerator)
	authOnlyEmployee := middleware.NewJWT(jwtPublicKey, middleware.AccessEmployee)
	authModerAndEmployee := middleware.NewJWT(jwtPublicKey, middleware.AccessModerator|middleware.AccessEmployee)

	// registration
	router.
		// no need auth
		Post("/dummyLogin", wrapper.PostDummyLogin).
		Post("/register", wrapper.PostRegister).
		Post("/login", wrapper.PostLogin).
		// moder only
		Post("/pvz", authOnlyModer, wrapper.PostPvz).
		// employee only
		Post("/pvz/:pvzId/close_last_reception", authOnlyEmployee, wrapper.PostPvzPvzIdCloseLastReception).
		Post("/pvz/:pvzId/delete_last_product", authOnlyEmployee, wrapper.PostPvzPvzIdDeleteLastProduct).
		Post("/receptions", authOnlyEmployee, wrapper.PostReceptions).
		Post("/products", authOnlyEmployee, wrapper.PostProducts).
		// moder and employee
		Get("/pvz", authModerAndEmployee, wrapper.GetPvz)

}
