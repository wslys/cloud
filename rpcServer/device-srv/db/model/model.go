package model

import (
	"beaconCloud/rpcServer/device-srv/db"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DatabaseUrl = "root:root@tcp(127.0.0.1:3306)/beacon_cloud?charset=utf8"
	o           orm.Ormer
)

func init() {
	judgeDBLink()

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", DatabaseUrl)

	// 需要在init中注册定义的model
	orm.RegisterModel(new(Gateway), new(BeaconSetting), new(Beacon))
	o = orm.NewOrm()
}

func UseingDB(url string) {
	if len(url) == 0 {
		err := o.Using("default")
		fmt.Print(err)
		return
	}

	DatabaseUrl = url
	if err := judgeDBLink(); err == nil {
		url1 := strings.Split(url, "/")
		databaseName := strings.Split(url1[1], "?")

		if err := o.Using(databaseName[0]); err != nil {
			if err := orm.RegisterDataBase(databaseName[0], "mysql", url); err != nil {
				fmt.Println(err)
			}else {
				o.Using(databaseName[0])
			}
		}

	} else {
		fmt.Println(err)
	}
}

func judgeDBLink() error {
	var d *sql.DB
	var err error

	parts := strings.Split(DatabaseUrl, "/")
	if len(parts) != 2 {
		panic("Invalid database url")
	}

	if len(parts[1]) == 0 {
		panic("Invalid database name")
	}

	url := parts[0]
	database := strings.Split(parts[1], "?")

	if d, err = sql.Open("mysql", url+"/"); err != nil {
		fmt.Print("database link fail.")
		return err
	}
	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + database[0] + " DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci"); err != nil {
		fmt.Print("database create fail.")
		return err
	}
	d.Close()
	if d, err = sql.Open("mysql", DatabaseUrl); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(db.GatewaySchema); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(db.BeaconSchema); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(db.BeaconSetting); err != nil {
		log.Fatal(err)
	}
	d.Close()
	return nil
}
