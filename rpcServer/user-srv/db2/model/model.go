package model

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DatabaseUrl = "root:root@tcp(127.0.0.1:3306)/beacon_cloud?charset=utf8"
	o           orm.Ormer
)

func init() {

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", DatabaseUrl)

	// 需要在init中注册定义的model
	orm.RegisterModel(new(Account))
	o = orm.NewOrm()
}

