package models

import (
	_ "github.com/go-sql-driver/mysql"
)


type Cluster struct {
	Id             int
	Name           string
	Apiserver_ip   string
	Apiserver_port string
	Cacrt          string
	Publickey      string
	Privitekey     string
	Env            []*Env `orm:"reverse(many)"`
	Registry       []*Registry `orm:"reverse(many)"`
}


