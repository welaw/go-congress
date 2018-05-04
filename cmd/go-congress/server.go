package gocongress

import (
	"net"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/welaw/go-congress/backend/filesystem"
	"github.com/welaw/go-congress/services"
	"google.golang.org/grpc"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/welaw/go-congress/backend/database"
	"github.com/welaw/go-congress/endpoints"
	"github.com/welaw/go-congress/instrumentation"
	"github.com/welaw/go-congress/logging"
	"github.com/welaw/go-congress/proto"
	"github.com/welaw/go-congress/transport"
	"github.com/welaw/welaw/client"
	welawproto "github.com/welaw/welaw/proto"
)

func makeServeCmd(_ *string, upstreamName *string) *cobra.Command {

	var grpcAddr string
	var environment string
	var debugAddr string

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "run the server",
		Long:  `Run the HTTP and gRPC servers.`,
		Run: func(cmd *cobra.Command, args []string) {
			runServeCmd(environment, upstreamName, grpcAddr)
		},
	}

	serveCmd.Flags().StringVar(&grpcAddr, "gocongress.grpc.addr", ":8091", "GoUSGPO gRPC (HTTP) listen address")
	serveCmd.Flags().StringVar(&debugAddr, "debug.addr", ":8092", "Debug and metrics listen address")

	return serveCmd
}

func runServeCmd(environment string, upstreamName *string, goUsaGrpcAddr string) {
	// dotenv
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// logging middleware
	var (
		logger = log.NewLogfmtLogger(os.Stderr)
	)
	logger = log.With(logger, "legacy", true)

	// tracing
	var (
		tracer    = stdopentracing.GlobalTracer()
		fieldKeys = []string{"method", "error"}
	)

	// database
	var (
		connStr = os.Getenv("POSTGRES_CONNECTION")
	)

	// Errors
	var (
		errc = make(chan error, 2)
	)

	var (
		welawUrl      = os.Getenv("WELAW_URL")
		fsPath        = os.Getenv("FILESYSTEM_PATH")
		welawUsername = os.Getenv("WELAW_USERNAME")
		welawPassword = os.Getenv("WELAW_PASSWORD")
	)

	client := client.NewClientBasicAuth(welawUrl, welawUsername, welawPassword)

	var (
		service services.Service
	)
	{
		db, err := database.NewDatabase(connStr)
		if err != nil {
			panic(err)
		}
		fs := filesystem.NewFilesystem(fsPath)
		service = services.NewService(
			db,
			fs,
			logger,
			welawproto.Upstream{
				Ident: *upstreamName,
				Name:  *upstreamName,
			},
			client,
		)
		service = logging.NewLoggingMiddleware(logger, service)
		service = instrumentation.NewInstrumentatingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "go_congress",
				Subsystem: "service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "go_congress",
				Subsystem: "service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			service,
		)

		grpcEndpoints := endpoints.Endpoints{
			SendVoteEndpoint: endpoints.MakeSendVoteEndpoint(service),
			SendLawEndpoint:  endpoints.MakeSendLawEndpoint(service),
			StatusEndpoint:   endpoints.MakeStatusEndpoint(service),
		}

		// ServeGRPC
		go func() {
			logger := log.With(logger, "transport", "gRPC")
			ln, err := net.Listen("tcp", goUsaGrpcAddr)
			if err != nil {
				errc <- err
				return
			}
			srv := transport.MakeGrpcServer(grpcEndpoints, tracer, logger)
			s := grpc.NewServer()
			proto.RegisterGoCongressServiceServer(s, srv)
			logger.Log("addr", goUsaGrpcAddr)
			errc <- s.Serve(ln)
		}()

		logger.Log("terminated", <-errc)
		db.Close()
	}

}
