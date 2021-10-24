package controllers


import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"xcloud/models"
	"xcloud/util"
	"github.com/astaxie/beego"
	"github.com/beego/beego/v2/client/orm"
)


type ImageController struct {
	beego.Controller
}

// @router / [get]
func (this *ImageController) Get() {
	clusterId := this.GetString("clusterId")
	beego.Info(clusterId,"9999999999999")
	cluster := models.Cluster{}
	o := orm.NewOrm()
	o.QueryTable("Cluster").Filter("Id", clusterId).RelatedSel().One(&cluster)
	o.LoadRelated(&cluster, "Registry")
	beego.Info(cluster.Registry[0].IsSsl,"00000000000000000")
	if (cluster.Registry[0].IsSsl == false){
		clientset := util.Getclient(cluster.Name)
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		registryIpList := []string{}
		if err == nil {
			for _, o := range nodes.Items {
				registryIpList = append(registryIpList, o.Status.Addresses[0].Address)
			}
			beego.Info(registryIpList, "00000000000009999999999999")
		}
		for _, ip := range registryIpList {
			resp, err := http.Get("http://" + ip + ":31000/v2/_catalog" )
			if (err == nil) {
				body, err:=ioutil.ReadAll(resp.Body)
				if err == nil {
					var images map[string]interface{}
					err := json.Unmarshal(body, &images)
					if (err == nil) {
						beego.Info(images,"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii")
						this.Data["json"] = images
						beego.Info(images)
						this.ServeJSON()
					}
				}
			}
		}
	}
}

// @router /imagetaglist/ [get]
func (this *ImageController) GetImageTag() {
	clusterId := this.GetString("clusterId")
	image := this.GetString("image")
	beego.Info(clusterId,"9999999999999")
	cluster := models.Cluster{}
	o := orm.NewOrm()
	o.QueryTable("Cluster").Filter("Id", clusterId).RelatedSel().One(&cluster)
	o.LoadRelated(&cluster, "Registry")
	beego.Info(cluster.Registry[0].IsSsl,"00000000000000000")
	if (cluster.Registry[0].IsSsl == false){
		clientset := util.Getclient(cluster.Name)
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		registryIpList := []string{}
		if err == nil {
			for _, o := range nodes.Items {
				registryIpList = append(registryIpList, o.Status.Addresses[0].Address)
			}
			beego.Info(registryIpList, "00000000000009999999999999")
		}
		for _, ip := range registryIpList {
			resp, err := http.Get("http://" + ip + ":31000/v2/" + image + "/tags/list" )
			if (err == nil) {
				body, err:=ioutil.ReadAll(resp.Body)
				var tags map[string]interface{}
				if err == nil {
					err := json.Unmarshal(body, &tags)
					if err == nil {
						this.Data["json"] = tags
						this.ServeJSON()
					}
				}
			}
		}
	}
}