package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type Cert struct {
	Id             int
	Remark         string
    Crt            string
	Key            string
	Registry       []*Registry `orm:"reverse(many)"`
}


