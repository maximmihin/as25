package main

import (
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	. "github.com/maximmihin/as25/cmd/httpserver/testclient"
	"github.com/maximmihin/as25/internal/dal"
	"github.com/maximmihin/as25/internal/services/accesscontrol"
)

const (
	DbUser = "pvz_db_user"
	DbPass = "pvz_db_pass"
	DbName = "pvz_db_name"

	jwtPublicKey  = "access_secret"
	jwtPrivateKey = jwtPublicKey // TODO make async

	logLevel = "DEBUG"
)

func TestE2E(t *testing.T) {

	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	connString := RunPostgres(t)
	PostgresPrepare(t, connString)
	serverHost, serverPort := RunApp(t, connString)

	tcl := NewTestClient(t, serverHost, serverPort)

	// authorized client to reuse between tests
	var authModer *TClient
	var authEmployee *TClient
	{
		employeeToken, _ := accesscontrol.NewDummyJWT(jwtPublicKey, "employee")
		moderatorToken, _ := accesscontrol.NewDummyJWT(jwtPublicKey, "moderator")
		authModer = tcl.WithBearer(moderatorToken)
		authEmployee = tcl.WithBearer(employeeToken)
	}

	t.Run("DummyLogin", func(t *testing.T) {

		t.Run("green path", func(t *testing.T) {
			res := tcl.DummyLogin(t, PostDummyLoginJSONBody{
				Role: "moderator",
			})
			require.Equal(t, 200, res.StatusCode())

			res2 := tcl.DummyLogin(t, PostDummyLoginJSONBody{
				Role: "employee",
			})
			require.Equal(t, 200, res2.StatusCode())
		})

		t.Run("invalid role", func(t *testing.T) {
			res := tcl.DummyLogin(t, PostDummyLoginJSONBody{
				Role: "terminator",
			})

			require.Equal(t, 400, res.StatusCode())
		})

	})

	t.Run("Create PVZ", func(t *testing.T) {

		t.Run("green path", func(t *testing.T) {
			res2 := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res2.StatusCode())
		})

		t.Run("permission denied", func(t *testing.T) {
			res2 := authEmployee.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 403, res2.StatusCode())
		})

		t.Run("unavailable city", func(t *testing.T) {
			res2 := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Novosibirsk",
			})
			require.Equal(t, 400, res2.StatusCode())
		})

	})

	t.Run("Create reception", func(t *testing.T) {

		t.Run("green path", func(t *testing.T) {

			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 201, res2.StatusCode())
		})

		t.Run("nonexistent Pvz Id", func(t *testing.T) {
			res := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: uuid.New(),
			})
			require.Equal(t, 400, res.StatusCode())
		})

		t.Run("previous reception in_progress", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 201, res2.StatusCode())

			res3 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 400, res3.StatusCode())
		})

	})

	t.Run("Create product", func(t *testing.T) {

		t.Run("green path", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 201, res2.StatusCode())

			res3 := authEmployee.CreateProduct(t, PostProductsJSONBody{
				PvzId: res.JSON201.Id,
				Type:  Обувь,
			})
			require.Equal(t, 201, res3.StatusCode())
		})

		t.Run("invalid pvz", func(t *testing.T) {
			res := authEmployee.CreateProduct(t, PostProductsJSONBody{
				PvzId: uuid.New(),
				Type:  Обувь,
			})
			require.Equal(t, 400, res.StatusCode())
		})

		t.Run("no reception in pvz", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateProduct(t, PostProductsJSONBody{
				PvzId: res.JSON201.Id,
				Type:  Обувь,
			})
			require.Equal(t, 400, res2.StatusCode())
		})

		t.Run("reception already close", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 201, res2.StatusCode())

			res3 := authEmployee.CloseLastReception(t, res.JSON201.Id)
			require.Equal(t, 200, res3.StatusCode())
			require.Equal(t, Close, res3.JSON200.Status)

			res4 := authEmployee.CreateProduct(t, PostProductsJSONBody{
				PvzId: res.JSON201.Id,
				Type:  Обувь,
			})
			require.Equal(t, 400, res4.StatusCode())
		})
	})

	t.Run("Close last reception", func(t *testing.T) {

		t.Run("green path", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 201, res2.StatusCode())

			res3 := authEmployee.CloseLastReception(t, res.JSON201.Id)
			require.Equal(t, 200, res3.StatusCode())
			require.Equal(t, Close, res3.JSON200.Status)
		})

		t.Run("nonexistent Pvz Id", func(t *testing.T) {
			res := authEmployee.CloseLastReception(t, uuid.New())
			require.Equal(t, 400, res.StatusCode())
		})

		t.Run("pvz without reception", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CloseLastReception(t, res.JSON201.Id)
			require.Equal(t, 400, res2.StatusCode())
		})

		t.Run("close after close", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 201, res2.StatusCode())

			res3 := authEmployee.CloseLastReception(t, res.JSON201.Id)
			require.Equal(t, 200, res3.StatusCode())
			require.Equal(t, Close, res3.JSON200.Status)

			res4 := authEmployee.CloseLastReception(t, res.JSON201.Id)
			require.Equal(t, 400, res4.StatusCode())
		})
	})

	t.Run("Delete last product", func(t *testing.T) {

		t.Run("green path", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 201, res2.StatusCode())

			res3 := authEmployee.CreateProduct(t, PostProductsJSONBody{
				PvzId: res.JSON201.Id,
				Type:  Обувь,
			})
			require.Equal(t, 201, res3.StatusCode())

			res4 := authEmployee.DeleteLastProduct(t, res.JSON201.Id)
			require.Equal(t, 204, res4.StatusCode())
		})

		t.Run("delete product from empty reception", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 201, res2.StatusCode())

			res3 := authEmployee.CreateProduct(t, PostProductsJSONBody{
				PvzId: res.JSON201.Id,
				Type:  Обувь,
			})
			require.Equal(t, 201, res3.StatusCode())

			res4 := authEmployee.DeleteLastProduct(t, res.JSON201.Id)
			require.Equal(t, 204, res4.StatusCode())

			res5 := authEmployee.DeleteLastProduct(t, res.JSON201.Id)
			require.Equal(t, 400, res5.StatusCode())
		})

		t.Run("delete product from closed reception", func(t *testing.T) {
			res := authModer.CreatePVZ(t, PostPvzJSONRequestBody{
				City: "Казань",
			})
			require.Equal(t, 201, res.StatusCode())

			res2 := authEmployee.CreateReception(t, PostReceptionsJSONBody{
				PvzId: res.JSON201.Id,
			})
			require.Equal(t, 201, res2.StatusCode())

			res3 := authEmployee.CreateProduct(t, PostProductsJSONBody{
				PvzId: res.JSON201.Id,
				Type:  Обувь,
			})
			require.Equal(t, 201, res3.StatusCode())

			res4 := authEmployee.CloseLastReception(t, res.JSON201.Id)
			require.Equal(t, 200, res4.StatusCode())

			res5 := authEmployee.DeleteLastProduct(t, res.JSON201.Id)
			require.Equal(t, 400, res5.StatusCode())
		})
	})

	//t.Run("Get Pvz", func(t *testing.T) {
	//
	//	t.Run("green path", func(t *testing.T) {
	//
	//		var createdPvz *PVZ
	//
	//		// moderator create PVZ
	//		{
	//			res := tcl.DummyLogin(t, PostDummyLoginJSONBody{
	//				Role: "moderator",
	//			})
	//			require.Equal(t, 200, res.StatusCode())
	//
	//			authModerator := tcl.WithBearer(*res.JSON200)
	//
	//			res2 := authModerator.CreatePVZ(t, PVZ{
	//				City: "Казань",
	//			})
	//			require.Equal(t, 201, res2.StatusCode())
	//			createdPvz = res2.JSON201
	//		}
	//
	//		// employee get created pvz
	//		{
	//			res := tcl.DummyLogin(t, PostDummyLoginJSONBody{
	//				Role: "employee",
	//			})
	//			require.Equal(t, 200, res.StatusCode())
	//
	//			authEmployee := tcl.WithBearer(*res.JSON200)
	//
	//			res2 := authEmployee.GetPVZ(t, GetPvzParams{
	//				StartDate: createdPvz.RegistrationDate,
	//				EndDate:   ptr(time.Now()),
	//			})
	//			require.Equal(t, 201, res2.StatusCode())
	//		}
	//
	//	})

}

func ptr[T any](v T) *T {
	return &v
}

func RunApp(t *testing.T, connString string) (serverHost, serverPort string) {

	cfg := Config{
		DbConnString:  connString,
		JwtPrivateKey: jwtPrivateKey,
		jwtPublicKey:  jwtPublicKey,

		Host: "localhost",
		Port: "0",

		LogLevel: logLevel,
	}
	t.Logf("server config: %#v", cfg)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	startListen := make(chan struct{})
	globalOnListenFunc = func(listenData fiber.ListenData) error {

		serverHost = listenData.Host
		serverPort = listenData.Port

		startListen <- struct{}{}
		return nil
	}

	errC, shutDownFunc := Run(cfg, logger)
	t.Cleanup(func() {
		shutDownFunc(t.Context())
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errC:
		t.Fatalf("fail on Run: %v", err)
	case <-startListen: // TODO log about start
		t.Logf("server started listen on %s:%s", serverHost, serverPort)
	case s := <-sig:
		t.Log("stopped by sig: ", s)
	}
	return
}

func RunPostgres(t *testing.T) (ConnStr string) {
	ctx := t.Context()

	//pgCtrName := "pvz_Postgres" + CurrentTimeToHMSString()
	pgCtrName := "pvz_Postgres"

	defer func() {
		handleFknPanicTestContainers(t, recover())
	}()
	t.Logf("starting %s...", pgCtrName)
	pgContainer, err := postgres.Run(ctx,
		"postgres:17.2-alpine3.21",
		withName(pgCtrName),
		withReuse(),
		postgres.WithDatabase(DbName),
		postgres.WithUsername(DbUser),
		postgres.WithPassword(DbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second)),
	)
	// TODO add cleanup terminate container
	require.NoError(t, err)

	ConnStr, err = pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	t.Logf("postgres conn string: %s", ConnStr)

	// print the port on which postgress is running
	{
		pgPort, err := pgContainer.MappedPort(t.Context(), "5432/tcp")
		require.NoError(t, err)
		t.Logf("%s started on %s", pgCtrName, pgPort)
	}

	return ConnStr
}

func withName(ctrName string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Name = ctrName
		return nil
	}
}

func withReuse() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Reuse = true
		return nil
	}
}

func handleFknPanicTestContainers(t *testing.T, some any) {
	if some == nil {
		return
	}
	err, ok := some.(error)
	if !ok {
		t.Fatal("testcontainers doesnt run for some reason: ", some)
	}
	if err.Error() == "rootless Docker not found" {
		t.Fatal("need running docker on host for run e2e tests - testcontainers will create and launch the necessary containers automatically")
	}
	t.Fatal("testcontainers fail to connect to docker on host: " + err.Error())
}

func PostgresPrepare(t *testing.T, pgConnString string) {

	db, err := sql.Open("pgx", pgConnString)
	require.NoError(t, err)

	goose.SetBaseFS(dal.EmbedMigrations)
	goose.SetLogger(gooseLogeAdapter{t})

	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	err = goose.DownToContext(t.Context(), db, "migrations", 0)
	require.NoError(t, err)
	t.Log("postgres migration down")

	err = goose.UpContext(t.Context(), db, "migrations")
	require.NoError(t, err)
	t.Log("postgres migration up")

}

type gooseLogeAdapter struct{ *testing.T }

func (g gooseLogeAdapter) Fatalf(format string, v ...interface{}) { g.T.Logf(format, v...) }
func (g gooseLogeAdapter) Printf(format string, v ...interface{}) { g.T.Logf(format, v...) }
