package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghaperf/pkg/cli"
	"github.com/suzuki-shunsuke/ghaperf/pkg/controller"
	"github.com/suzuki-shunsuke/ghaperf/pkg/log"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

var version = ""

func main() {
	if code := core(); code != 0 {
		os.Exit(code)
	}
}

func core() int {
	logger, logLevel := log.New(os.Stderr, version)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	home, err := os.UserHomeDir()
	if err != nil {
		slogerr.WithError(logger, err).Error("failed to get user home dir")
		return 1
	}
	if err := cli.Run(ctx, logger, logLevel, &controller.Arg{
		Getenv:  os.Getenv,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Fs:      afero.NewOsFs(),
		Args:    os.Args,
		Home:    home,
		Version: version,
	}); err != nil {
		slogerr.WithError(logger, err).Error("ghaperf failed")
		return 1
	}
	return 0
}
