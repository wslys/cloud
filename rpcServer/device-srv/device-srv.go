package main

import (
	"beaconCloud/rpcServer/device-srv/db/model"
	"beaconCloud/rpcServer/device-srv/handler"
	"beaconCloud/rpcServer/device-srv/proto/device"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"log"
	"github.com/widuu/goini"
	"github.com/micro/go-micro/errors"
	"os"
)

var (
	conf *goini.Config
)

func init() {
	conf = goini.SetConfig("/etc/beacon-cloud/conf/conf.ini")
	username := conf.GetValue("database", "username")
	password := conf.GetValue("database", "password")
	hostname := conf.GetValue("database", "hostname")

	if len(username) == 0 || len(password) == 0 || len(hostname) == 0 {
		log.Fatal(errors.New("Device-SRV Init", "Config Database Username or Password or hostname IS NULL.", 500))
		os.Exit(1)
	}
	if username == "no value" || password == "no value" || hostname == "no value" {
		log.Fatal(errors.New("Device-SRV Init", "Config Database Username or Password or hostname IS NULL.", 500))
		os.Exit(1)
	}
	model.DatabaseUrl = username + ":" + password + "@tcp("+ hostname +":3306)/beacon_cloud?charset=utf8"
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.device-srv"),
		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/gateway",
			},
		),

		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				model.DatabaseUrl = c.String("database_url")
			}
		}),
	)

	service.Init()
	// db.Init()

	device.RegisterDeviceHandler(service.Server(), new(handler.Device))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
