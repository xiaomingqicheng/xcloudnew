package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:Root@123@tcp(192.168.0.103:3306)/xcloud?charset=utf8")
	orm.RegisterModel(new(Env),new(Cluster),new(Cert),new(Registry))
	orm.RunSyncdb("default", false, true)
}
