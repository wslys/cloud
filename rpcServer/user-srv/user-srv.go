package main

import (
	"log"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"beaconCloud/rpcServer/user-srv/db"
	"beaconCloud/rpcServer/user-srv/handler"
	proto "beaconCloud/rpcServer/user-srv/proto/account"
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
	db.Url = username + ":" + password + "@tcp("+ hostname +":3306)/user?charset=utf8"
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/user",
			},
		),

		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				db.Url = c.String("database_url")
			}
		}),
	)

	service.Init()
	db.Init()

	proto.RegisterAccountHandler(service.Server(), new(handler.Account))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
