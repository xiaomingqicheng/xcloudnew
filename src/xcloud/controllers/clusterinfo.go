package controllers

import (
	"encoding/json"
	"xcloud/models"
	//"fmt"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/beego/beego/v2/client/orm"
)

// Operations about object
type ClusterController struct {
	beego.Controller
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (this *ClusterController) Post() {
	data := this.Ctx.Input.RequestBody
	var mapResult map[string]interface{}
	json.Unmarshal([]byte(data), &mapResult)
	beego.Info(mapResult,"===============================")
	o := orm.NewOrm()
	o.Begin()
	cluster := new(models.Cluster)
	cluster.Name = mapResult["name"].(string)
	cluster.Apiserver_ip = mapResult["apiserver_ip"].(string)
	cluster.Apiserver_port = mapResult["apiserver_port"].(string)
	cluster.Cacrt = mapResult["cacrt"].(string)
	cluster.Publickey = mapResult["publickey"].(string)
	cluster.Privitekey = mapResult["privitekey"].(string)

	env_id := int(mapResult["env"].(float64))
	_, err := o.Insert(cluster)
	if err != nil{
		beego.Info(err, "656564546546546546565654654")
	}
	o1 := orm.NewOrm()
	var env []*models.Env
	o1.QueryTable("Env").Filter("Id__in",env_id).All(&env)
	m2m := o.QueryM2M(cluster, "Env")
	for _, envv := range env{
		num, err := m2m.Add(envv)
		if err == nil {
			beego.Info("Added nums: %v", num)
		}
	}
	this.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (this *ClusterController) Get() {
	uid := this.GetString(":uid")
	beego.Info(uid,"--------------------------------")
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (this *ClusterController) GetAll() {
	o := orm.NewOrm()
	var clusters []*models.Cluster
	o.QueryTable("Cluster").All(&clusters)
	for _,v := range clusters {
		o.LoadRelated(v, "Env")
	}
	this.Data["json"] = clusters
	this.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Cluster
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Cluster	true		"body for Cluster content"
// @Success 200 {object} models.Cluster
// @Failure 403 :id is not int
// @router /:id [put]
func (this *ClusterController) Put() {
	id := this.Ctx.Input.Param(":id")
	data := this.Ctx.Input.RequestBody
	var mapResult map[string]interface{}
	json.Unmarshal([]byte(data), &mapResult)
	beego.Info(mapResult,"===============================")
	o := orm.NewOrm()
	o.Begin()
	cluster := new(models.Cluster)
	cluster.Name = mapResult["name"].(string)
	cluster.Apiserver_ip = mapResult["apiserver_ip"].(string)
	cluster.Apiserver_port = mapResult["apiserver_port"].(string)
	cluster.Cacrt = mapResult["cacrt"].(string)
	cluster.Publickey = mapResult["publickey"].(string)
	cluster.Privitekey = mapResult["privitekey"].(string)
	cluster.Id,_ = strconv.Atoi(id)
	_, err := o.Update(cluster)
	if err != nil{
		beego.Info(err, "656564546546546546565654654")
	}
	env_id := int(mapResult["env"].(float64))
	o1 := orm.NewOrm()
	var env []*models.Env
	o1.QueryTable("Env").Filter("Id__in",env_id).All(&env)
	m2m := o.QueryM2M(cluster, "Env")
	for _, envv := range env{
		m2m.Clear()
		num, err := m2m.Add(envv)
		if err == nil {
			beego.Info("Added nums: %v", num)
		}
	}
	this.ServeJSON()
}

