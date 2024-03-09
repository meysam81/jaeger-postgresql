package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc/shared"
	"github.com/jaegertracing/jaeger/storage/dependencystore"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robbert229/fxslog"
	_ "github.com/robbert229/fxslog"
	"github.com/robbert229/jaeger-postgresql/internal/sql"
	"github.com/robbert229/jaeger-postgresql/internal/store"
	"google.golang.org/grpc"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// ProvideLogger returns a function that provides a logger
func ProvideLogger() any {
	return func() (*slog.Logger, error) {
		levelFn := func() (slog.Level, error) {
			if loglevelFlag == nil {
				return slog.LevelWarn, nil
			}

			switch *loglevelFlag {
			case "info":
				return slog.LevelInfo, nil
			case "warn":
				return slog.LevelWarn, nil
			case "error":
				return slog.LevelError, nil
			case "debug":
				return slog.LevelDebug, nil
			default:
				return 0, fmt.Errorf("invalid log level: %s", *loglevelFlag)
			}
		}
		level, err := levelFn()
		if err != nil {
			return nil, fmt.Errorf("failed to build logger: %w", err)
		}

		return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})), nil
	}
}

// ProvidePgxPool returns a function that provides a pgx pool
func ProvidePgxPool() any {
	return func(logger *slog.Logger, lc fx.Lifecycle) (*pgxpool.Pool, error) {
		if databaseURLFlag == nil {
			return nil, fmt.Errorf("invalid database url")
		}

		databaseURL := *databaseURLFlag
		if databaseURL == "" {
			return nil, fmt.Errorf("invalid database url")
		}

		err := sql.Migrate(logger, databaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}

		pgxconfig, err := pgxpool.ParseConfig(databaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse database url")
		}

		var maxConns int32
		if databaseMaxConnsFlag == nil {
			maxConns = 20
		} else {
			maxConns = int32(*databaseMaxConnsFlag)
		}

		pgxconfig.MaxConns = maxConns

		pool, err := pgxpool.NewWithConfig(context.Background(), pgxconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to the postgres database: %w", err)
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				pool.Close()
				return nil
			},
		})

		return pool, nil
	}
}

// ProvideSpanStoreReader returns a function that provides a spanstore reader.
func ProvideSpanStoreReader() any {
	return func(pool *pgxpool.Pool, logger *slog.Logger) spanstore.Reader {
		q := sql.New(pool)
		return store.NewInstrumentedReader(store.NewReader(q, logger), logger)
	}
}

// ProvideSpanStoreWriter returns a function that provides a spanstore writer
func ProvideSpanStoreWriter() any {
	return func(pool *pgxpool.Pool, logger *slog.Logger) spanstore.Writer {
		q := sql.New(pool)
		return store.NewInstrumentedWriter(store.NewWriter(q, logger), logger)
	}
}

// ProvideDependencyStoreReader provides a dependencystore reader
func ProvideDependencyStoreReader() any {
	return func(pool *pgxpool.Pool, logger *slog.Logger) dependencystore.Reader {
		q := sql.New(pool)
		return store.NewReader(q, logger)
	}
}

// ProvideHandler provides a grpc handler.
func ProvideHandler() any {
	return func(reader spanstore.Reader, writer spanstore.Writer, dependencyReader dependencystore.Reader) *shared.GRPCHandler {
		handler := shared.NewGRPCHandler(&shared.GRPCHandlerStorageImpl{
			SpanReader:          func() spanstore.Reader { return reader },
			SpanWriter:          func() spanstore.Writer { return writer },
			DependencyReader:    func() dependencystore.Reader { return dependencyReader },
			ArchiveSpanReader:   func() spanstore.Reader { return nil },
			ArchiveSpanWriter:   func() spanstore.Writer { return nil },
			StreamingSpanWriter: func() spanstore.Writer { return nil },
		})

		return handler
	}
}

// ProvideGRPCServer provides a grpc server.
func ProvideGRPCServer() any {
	return func(lc fx.Lifecycle) (*grpc.Server, error) {
		srv := grpc.NewServer()

		if grpcHostPort == nil {
			return nil, fmt.Errorf("invalid grpc-server.host-port")
		}

		lis, err := net.Listen("tcp", *grpcHostPort)
		if err != nil {
			return nil, fmt.Errorf("failed to listen: %w", err)
		}

		lc.Append(fx.StartStopHook(
			func(ctx context.Context) error {
				go srv.Serve(lis)
				return nil
			},

			func(ctx context.Context) error {
				srv.GracefulStop()
				return lis.Close()
			},
		))

		return srv, nil
	}
}

// ProvideAdminServer provides the admin http server.
func ProvideAdminServer() any {
	return func(lc fx.Lifecycle) (*http.ServeMux, error) {
		mux := http.NewServeMux()

		srv := http.Server{
			Handler: mux,
		}

		lis, err := net.Listen("tcp", *adminHttpHostPort)
		if err != nil {
			return nil, fmt.Errorf("failed to listen: %w", err)
		}

		lc.Append(fx.StartStopHook(
			func(ctx context.Context) error {
				go srv.Serve(lis)
				return nil
			},

			func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		))

		return mux, nil
	}
}

var (
	databaseURLFlag      = flag.String("database.url", "", "the postgres connection url to use to connect to the database")
	databaseMaxConnsFlag = flag.Int("database.max-conns", 20, "Max number of database connections of which the plugin will try to maintain at any given time")
	loglevelFlag         = flag.String("log-level", "warn", "Minimal allowed log level")
	grpcHostPort         = flag.String("grpc-server.host-port", ":12345", "the host:port (eg 127.0.0.1:12345 or :12345) of the storage provider's gRPC server")
	adminHttpHostPort    = flag.String("admin.http.host-port", ":12346", "The host:port (e.g. 127.0.0.1:12346 or :12346) for the admin server, including health check, /metrics, etc.")
)

func main() {
	flag.Parse()

	fx.New(
		fx.WithLogger(
			func(logger *slog.Logger) fxevent.Logger {
				return &fxslog.SlogLogger{Logger: logger.With("component", "uber/fx")}
			},
		),
		fx.Provide(
			ProvideLogger(),
			ProvidePgxPool(),
			ProvideSpanStoreReader(),
			ProvideSpanStoreWriter(),
			ProvideDependencyStoreReader(),
			ProvideHandler(),
			ProvideGRPCServer(),
			ProvideAdminServer(),
		),
		fx.Invoke(func(srv *grpc.Server, handler *shared.GRPCHandler) error {
			return handler.Register(srv)
		}),
		fx.Invoke(func(mux *http.ServeMux, conn *pgxpool.Pool) {
			mux.Handle("/metrics", promhttp.Handler())
			mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx, cancelFn := context.WithTimeout(r.Context(), time.Second*5)
				defer cancelFn()

				err := conn.Ping(ctx)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}

				w.WriteHeader(http.StatusOK)
			}))
		}),
	).Run()
}