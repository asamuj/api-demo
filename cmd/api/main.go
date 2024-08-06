package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/asamuj/api-demo/api/server"
	"github.com/asamuj/api-demo/api/service"
	"github.com/asamuj/api-demo/cmd/runtime/version"
	"github.com/asamuj/api-demo/config"
	"github.com/asamuj/api-demo/database/mysql"
)

var (
	// configPathFlag specifies the api config file path.
	configPathFlag = &cli.StringFlag{
		Name:     "config-file",
		Usage:    "The filepath to a json file, flag is required",
		Required: true,
	}
)

func main() {
	app := cli.App{
		Name:    "api-demo",
		Usage:   "this is an api demo",
		Action:  exec,
		Version: version.Get(),
		Flags: []cli.Flag{
			configPathFlag,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("running api application failed")
	}
}

func exec(ctx *cli.Context) error {
	cfg := &Config{}
	if err := config.Load(ctx.String(configPathFlag.Name), cfg); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("reading api config failed")
	}

	db, err := mysql.NewMySQLDB(cfg.MySQL)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalln("initialize mysql db error")
	}

	log.Info("Starting explorer api server...")

	server.New(
		cfg.Port,
		service.New(db),
	).Run()
	return nil
}

// Config defines the config for api service.
type Config struct {
	Port  int          `yaml:"port"`
	MySQL mysql.Config `yaml:"mysql"`
}
