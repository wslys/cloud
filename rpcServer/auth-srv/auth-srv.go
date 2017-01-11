package main

import (
	"beaconCloud/rpcServer/auth-srv/db"
	"beaconCloud/rpcServer/auth-srv/db/mysql"
	"beaconCloud/rpcServer/auth-srv/handler"
	account "beaconCloud/rpcServer/auth-srv/proto/account"
	oauth2 "beaconCloud/rpcServer/auth-srv/proto/oauth2"
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
	mysql.Url = username + ":" + password + "@tcp("+ hostname +":3306)/auth?charset=utf8"
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.auth"),
		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root:root@tcp(127.0.0.1:3306)/auth",
			},
		),
		micro.Version("v1.0"),
		// micro.Server(s)
		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				mysql.Url = c.String("database_url")
			}
		}),
	)

	// initialise service
	service.Init()

	// register account handler
	account.RegisterAccountHandler(service.Server(), new(handler.Account))
	// register oauth2 handler
	oauth2.RegisterOauth2Handler(service.Server(), new(handler.Oauth2))
	// initialise database
	if err := db.Init(); err != nil {
		log.Println("db init err message ", err)
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
