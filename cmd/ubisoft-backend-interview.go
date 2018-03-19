package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"

	"github.com/kwiesmueller/ubisoft-backend-interview/pkg/feedback"

	"github.com/kwiesmueller/ubisoft-backend-interview/pkg/database"

	"github.com/golang/glog"
	"github.com/kolide/kit/version"
	"github.com/playnet-public/libs/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	appName = "ubisoft-backend-interview"
	appKey  = "ubisoft-backend-interview"
)

var (
	maxprocs    = flag.Int("maxprocs", runtime.NumCPU(), "max go procs")
	dbg         = flag.Bool("debug", false, "enable debug mode")
	versionInfo = flag.Bool("version", true, "show version info")
	sentryDsn   = flag.String("sentryDsn", "", "sentry dsn key")

	dbHost     = flag.String("dbHost", "127.0.0.1", "database hostname")
	dbPort     = flag.Int("dbPort", 5432, "database port")
	dbUsername = flag.String("dbUsername", "db", "database username")
	dbName     = flag.String("dbName", "db", "database name")
	dbPassword = flag.String("dbPassword", "db", "database password")
)

func main() {
	flag.Parse()

	if *versionInfo {
		v := version.Version()
		fmt.Printf("-- %s --\n", appName)
		fmt.Printf(" - version: %s\n", v.Version)
		fmt.Printf("   branch: \t%s\n", v.Branch)
		fmt.Printf("   revision: \t%s\n", v.Revision)
		fmt.Printf("   build date: \t%s\n", v.BuildDate)
		fmt.Printf("   build user: \t%s\n", v.BuildUser)
		fmt.Printf("   go version: \t%s\n", v.GoVersion)
	}
	runtime.GOMAXPROCS(*maxprocs)

	defer glog.Flush()
	glog.CopyStandardLogTo("info")

	var zapFields []zapcore.Field
	if !*dbg {
		zapFields = []zapcore.Field{
			zap.String("app", appKey),
			zap.String("version", version.Version().Version),
		}
	}

	log := log.New(appKey, *sentryDsn, *dbg).WithFields(zapFields...)
	defer log.Sync()
	log.Info("starting")

	if err := do(log); err != nil {
		log.Fatal("terminating", zap.Error(err))
	}
}

func do(log *log.Logger) error {
	db := database.New(log)
	err := db.Open(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		*dbHost,
		*dbPort,
		*dbUsername,
		*dbPassword,
		*dbName,
	))
	if err != nil {
		return err
	}

	svc := feedback.New(log, db)

	m := http.NewServeMux()
	m.Handle("/", svc.Handler())

	log.Info("listening", zap.String("addr", ":8080"))
	err = http.ListenAndServe(":8080", m)
	if err != nil {
		log.Error("server error", zap.Error(err))
		return err
	}
	return nil
}
