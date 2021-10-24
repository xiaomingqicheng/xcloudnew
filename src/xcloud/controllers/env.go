package controllers

import (
	"encoding/json"
	"strconv"
	"xcloud/models"
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/beego/beego/v2/client/orm"
)

// Operations about Env
type EnvController struct {
	beego.Controller
}

// @Title Create
// @Description create Env
// @Param	body		body 	models.Env	true		"The env content"
// @Success 200 {string} models.Env.Id
// @Failure 403 body is empty
// @router / [post]
func (this *EnvController) Post() {
	var cl map[string]interface{}
	json.Unmarshal([]byte(this.Ctx.Input.RequestBody), &cl)
	var env models.Env
	json.Unmarshal(this.Ctx.Input.RequestBody, &env)
	//var mapResult map[string]interface{}
	//json.Unmarshal([]byte(data), &mapResult)
	//beego.Info(cl,"===============================")
	o := orm.NewOrm()
	o.Begin()
	//env := new(models.Env)
	//env.Name = mapResult["name"].(string)
	//env.Clusters = clusters.([...]*models.Cluster)

	o.Insert(&env)
	o1 := orm.NewOrm()
	var clusters []*models.Cluster
	o1.QueryTable("Cluster").Filter("Id__in",cl["clusters"]).All(&clusters)
	m2m := o.QueryM2M(&env, "Clusters")
	for _, clu := range clusters{
	num, err := m2m.Add(clu)
		if err == nil {
			beego.Info("Added nums: %v", num)
		}
	}
	beego.Info(env, "==----------------")
	this.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (this *EnvController) GetAll() {
	o := orm.NewOrm()
	var env []*models.Env
	o.QueryTable("Env").RelatedSel().All(&env)
	for _,v := range env {
		o.LoadRelated(v, "Clusters")
	}

	this.Data["json"] = env
	this.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:id [get]
func (this *EnvController) Get() {
	id,err := strconv.Atoi(this.GetString(":id"))
	beego.Info(err,"-----------rrrrrrr---------------------")
	env := models.Env{}
	o := orm.NewOrm()
	o.QueryTable("Env").Filter("Id", id).RelatedSel().One(&env)
	o.LoadRelated(&env, "Clusters")
	this.Data["json"] = env
	this.ServeJSON()
}