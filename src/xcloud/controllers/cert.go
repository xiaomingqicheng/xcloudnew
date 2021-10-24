package controllers

import (
	"encoding/json"
	"xcloud/models"
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/beego/beego/v2/client/orm"
)

// Operations about object
type CertController struct {
	beego.Controller
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (this *CertController) Post() {
	var cert map[string]interface{}
	json.Unmarshal([]byte(this.Ctx.Input.RequestBody), &cert)
	beego.Info(cert,"===============================")
	o := orm.NewOrm()
	ct := new(models.Cert)
	ct.Remark = cert["remark"].(string)
	ct.Crt = cert["crt"].(string)
	ct.Key = cert["key"].(string)
	_, err := o.Insert(ct)
	if err != nil{
		beego.Info(err, "656564546546546546565654654")
	}
	this.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (this *CertController) GetAll() {
	o := orm.NewOrm()
	var cert []*models.Cert
	o.QueryTable("Cert").RelatedSel().All(&cert)
	//for _,v := range cert {
	//	o.LoadRelated(v, "Clusters")
	//}

	this.Data["json"] = cert
	this.ServeJSON()
}
