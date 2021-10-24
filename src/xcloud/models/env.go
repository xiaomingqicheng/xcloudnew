package models

import (
	_ "github.com/go-sql-driver/mysql"
)




type Env struct {
	Id             int
	Name           string
    Clusters       []*Cluster `orm:"rel(m2m)"`
	Registry      []*Registry `orm:"reverse(many)"`
}


