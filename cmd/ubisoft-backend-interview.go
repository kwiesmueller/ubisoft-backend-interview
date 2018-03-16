package main

import (
	"flag"
	"fmt"
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
)

func main() {
	flag.Parse()

	if *versionInfo {
		fmt.Printf("-- %s --\n", appName)
		version.PrintFull()
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
	err := db.Open(`host=localhost port=5432 user=db password=db dbname=db sslmode=disable`)
	if err != nil {
		return err
	}

	svc := feedback.New(log, db)

	err = svc.Add(feedback.Entry{SessionID: "1", UserID: "a", Rating: 1, Comment: ""})
	if err != nil {
		log.Error("add error", zap.Error(err))
	}
	err = svc.Add(feedback.Entry{SessionID: "1", UserID: "b", Rating: 2, Comment: "abc"})
	if err != nil {
		log.Error("add error", zap.Error(err))
	}
	err = svc.Add(feedback.Entry{SessionID: "2", UserID: "a", Rating: 2, Comment: "abc"})
	if err != nil {
		log.Error("add error", zap.Error(err))
	}
	err = svc.Add(feedback.Entry{SessionID: "2", UserID: "b", Rating: 2, Comment: "abc"})
	if err != nil {
		log.Error("add error", zap.Error(err))
	}

	fmt.Println(svc.GetLatest(1))

	return nil
}
