package controllers

import (
	"encoding/json"
	//"strconv"
	//"xcloud/models"
	//"fmt"
	//apiv1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/astaxie/beego"
	//"github.com/beego/beego/v2/client/orm"
)

type ConfigmapController struct {
	beego.Controller
}


// @router / [post]
func (this *ConfigmapController) Post() {
	var cm map[string]interface{}
	json.Unmarshal([]byte(this.Ctx.Input.RequestBody), &cm)
	beego.Info(cm)
	dataMap := make(map[string]interface{})
	for _,v := range cm["keyvalue"].([]interface{}) {
		dataMap[v.(map[string]interface{})["key"].(string)] = v.(map[string]interface{})["value"].(string)
	}
	beego.Info(dataMap)
	configmap_yaml := &apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind: "configmap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: cm["configmapname"].(string),
		},
		Data: dataMap,
	}
	_, err3 := clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), configmap_yaml, metav1.CreateOptions{})
	beego.Info(err3)
	if err3 != nil {
		this.CustomAbort(500 ,err3.Error())
	}
}