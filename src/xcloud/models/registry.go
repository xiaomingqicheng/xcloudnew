package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type Registry struct {
	Id             int
	Name           string
    Env            *Env `orm:"rel(fk)"`
	Cluster        *Cluster `orm:"rel(fk)"`
	Domain         string
	Hostpath       string
	Cert           *Cert `orm:"null;rel(fk)"`
	IsSsl          bool
}


