package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bluesky-social/indigo/api"
	"github.com/bluesky-social/indigo/bgs"
	"github.com/bluesky-social/indigo/carstore"
	cliutil "github.com/bluesky-social/indigo/cmd/gosky/util"
	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/indexer"
	"github.com/bluesky-social/indigo/notifs"
	"github.com/bluesky-social/indigo/plc"
	"github.com/bluesky-social/indigo/repomgr"

	_ "net/http/pprof"

	logging "github.com/ipfs/go-log"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/plugin/opentelemetry/tracing"
)

var log = logging.Logger("bigsky")

func init() {
	//logging.SetAllLoggers(logging.LevelDebug)
	logging.SetAllLoggers(logging.LevelWarn)
}

func main() {

	// only try dotenv if it exists
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	run(os.Args)
}

func run(args []string) {

	app := cli.App{
		Name:  "bigsky",
		Usage: "atproto BGS/firehose daemon",
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name: "jaeger",
		},
		&cli.StringFlag{
			Name:    "db",
			Value:   "sqlite=data/bigsky/bgs.sqlite",
			EnvVars: []string{"DATABASE_URL"},
		},
		&cli.StringFlag{
			Name:    "carstoredb",
			Value:   "sqlite=data/bigsky/carstore.sqlite",
			EnvVars: []string{"CARSTORE_DATABASE_URL"},
		},
		&cli.StringFlag{
			Name:    "carstore",
			Value:   "data/bigsky/carstore",
			EnvVars: []string{"CARSTORE_DIRECTORY"},
		},
		&cli.BoolFlag{
			Name: "dbtracing",
		},
		&cli.StringFlag{
			Name:    "plc",
			Usage:   "hostname of the plc server (including https:// prefix)",
			Value:   "https://plc.directory",
			EnvVars: []string{"ATP_PLC_HOST"},
		},
		&cli.BoolFlag{
			Name: "ssl-events",
		},
		&cli.BoolFlag{
			Name:  "aggregation",
			Value: true,
		},
	}

	app.Action = func(cctx *cli.Context) error {

		if cctx.Bool("jaeger") {
			url := "http://localhost:14268/api/traces"
			exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
			if err != nil {
				return err
			}
			tp := tracesdk.NewTracerProvider(
				// Always be sure to batch in production.
				tracesdk.WithBatcher(exp),
				// Record information about this application in a Resource.
				tracesdk.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String("bgs"),
					attribute.String("environment", "test"),
					attribute.Int64("ID", 1),
				)),
			)

			otel.SetTracerProvider(tp)
		}

		// ensure data directory exists; won't error if it does
		os.MkdirAll("data/bigsky/", os.ModePerm)

		dbstr := cctx.String("db")
		db, err := cliutil.SetupDatabase(dbstr)
		if err != nil {
			return err
		}

		if cctx.Bool("dbtracing") {
			if err := db.Use(tracing.NewPlugin()); err != nil {
				return err
			}
		}

		carstoredbstr := cctx.String("carstoredb")
		cardb, err := cliutil.SetupDatabase(carstoredbstr)
		if err != nil {
			return err
		}

		csdir := cctx.String("carstore")
		os.MkdirAll(filepath.Dir(csdir), os.ModePerm)
		cstore, err := carstore.NewCarStore(cardb, csdir)
		if err != nil {
			return err
		}

		didr := &api.PLCServer{Host: cctx.String("plc")}
		cachedidr := plc.NewCachingDidResolver(didr, time.Minute*5, 1000)

		kmgr := indexer.NewKeyManager(cachedidr, nil)

		repoman := repomgr.NewRepoManager(db, cstore, kmgr)

		dbp, err := events.NewDbPersistence(db, cstore)
		if err != nil {
			return fmt.Errorf("setting up db event persistence: %w", err)
		}

		evtman := events.NewEventManager(dbp)

		go evtman.Run()

		notifman := &notifs.NullNotifs{}

		ix, err := indexer.NewIndexer(db, notifman, evtman, cachedidr, repoman, true, cctx.Bool("aggregation"))
		if err != nil {
			return err
		}

		repoman.SetEventHandler(func(ctx context.Context, evt *repomgr.RepoEvent) {
			if err := ix.HandleRepoEvent(ctx, evt); err != nil {
				log.Errorw("failed to handle repo event", "err", err)
			}
		})

		bgs := bgs.NewBGS(db, ix, repoman, evtman, cachedidr, cctx.Bool("ssl-events"))

		// set up pprof endpoint
		go func() {
			if err := bgs.StartDebug("localhost:2471"); err != nil {
				panic(err)
			}
		}()

		return bgs.Start(":2470")
	}

	app.RunAndExitOnError()
}
