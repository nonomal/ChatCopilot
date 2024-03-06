package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lw396/WeComCopilot/crontab"
	"github.com/lw396/WeComCopilot/internal/repository/gorm"
	"github.com/lw396/WeComCopilot/pkg/cache"
	"github.com/lw396/WeComCopilot/service"
	"github.com/urfave/cli/v3"
)

var scheduleCmd = &cli.Command{
	Name:  "crontab",
	Usage: "启动定时服务",
	Flags: []cli.Flag{
		&cli.UintFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   6977,
			Usage:   "端口号",
		},
	},
	Before: func(c context.Context, cmd *cli.Command) (err error) {
		ctx, err = buildContext(cmd, "app")
		if err != nil {
			return err
		}
		return nil
	},
	Action: func(c context.Context, cmd *cli.Command) error {
		db, err := ctx.buildDB()
		if err != nil {
			return err
		}

		service := service.New(
			service.WithRepository(gorm.New(db)),
			service.WithLogger(ctx.buildLogger("CRONTAB")),
			service.WithSQLite(ctx.buildSQLite()),
			service.WithCache(cache.DefaultStore()),
		)

		s := crontab.NewServer(service)

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sig
			s.Stop()
			os.Exit(0)
		}()

		defer func() {
			s.Stop()
		}()
		return s.Start(c)
	},
}
