package controllers

import (
	"context"
	//restclient "k8s.io/client-go/rest"
	//"github.com/beego/beego/v2/client/orm"
	//"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"xcloud/models"
	//"context"
	//"encoding/json"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//apiv1 "k8s.io/api/core/v1"
	//"log"
	////"encoding/json"
	//"xcloud/models"
	"xcloud/util"
	//"fmt"
	"github.com/astaxie/beego"



	//"k8s.io/client-go/util/retry"
)

// Operations about object
type ClusterresdetailController struct {
	beego.Controller
}

// @Title Get
// @Description get cluster by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Cluster
// @Failure 403 :uid is empty
// @router / [get]
func (this *ClusterresdetailController) Get() {
	label := this.GetString("label")
	name := this.GetString("clustername")
	//beego.Info(name,"--------------------------------")
	//o := orm.NewOrm()
	//var cluster models.Cluster
	//o.QueryTable("Cluster").Filter("name", name).One(&cluster)
	//config := restclient.Config{}
	//config.CAData=[]byte(cluster.Cacrt)
	//tlsCfg := restclient.TLSClientConfig{
	//	Insecure: false,
	//	CAData :[]byte(cluster.Cacrt),
	//	KeyData:[]byte(cluster.Privitekey),
	//	CertData:[]byte(cluster.Publickey),
	//}
	//config.TLSClientConfig = tlsCfg
	//var port string = cluster.Apiserver_port
	//var master string = cluster.Apiserver_ip
	//config.Host = "https://" + master + ":" + port

	clientset := util.Getclient(name)
	if label == "NODES" {
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			arry := [](map[string]interface{}){}
			for _, o := range nodes.Items {
				json := make(map[string]interface{})
				json["name"] = o.Name
				//json["Conditions"] = o.Status
				json["ip"] = o.Status.Addresses[0].Address
				for _, status := range o.Status.Conditions {
					//beego.Info(status.Type,"88888888888888888888888888888")
					if status.Type == "Ready" {
						json["status"] = status.Status
					}
				}
				json["createtime"] = o.CreationTimestamp
				json["AllocatableMemory"] = o.Status.Allocatable.Memory().String()
				json["AllocatableCpu"] = o.Status.Allocatable.Cpu().String()
				arry = append(arry, json)
			}
			this.Data["json"] = arry
			this.ServeJSON()
		}
	}
	if label == "PODS" {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			arry := [](map[string]interface{}){}

			for _, o := range pods.Items {
				//beego.Info(o,"===============================888")
				json := make(map[string]interface{})
				json["name"] = o.Name
				json["createtime"] = o.CreationTimestamp
				json["Labels"] = o.Labels
				json["Namespace"] = o.Namespace
				json["HostIP"] = o.Status.HostIP
				json["Conditions"] = o.Status.Conditions
				json["PodIP"] = o.Status.PodIP
				json["StartTime"] = o.Status.StartTime
				//json["RestartCount"] = o.Status.ContainerStatuses[0].RestartCount
				arry = append(arry, json)
			}
			this.Data["json"] = arry
			this.ServeJSON()
		}
	}
	if label == "SVC" {
		services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})

		if err == nil {
			arry := [](map[string]interface{}){}
			for _, o := range services.Items {
				json := make(map[string]interface{})
				json["name"] = o.Name
				json["createtime"] = o.CreationTimestamp
				json["Namespace"] = o.Namespace
				json["ClusterIP"] = o.Spec.ClusterIP
				//json["Ports"] = o.Spec.Ports
				arry = append(arry, json)
			}
			beego.Info(arry, "===============================888")
			this.Data["json"] = arry
			this.ServeJSON()
		}
	}
	if label == "DEPLOYMENTS" {
		deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		//beego.Info(deployments, "000000000000000000000000000")
		if err == nil {
			arry := [](map[string]interface{}){}
			beego.Info(deployments.Items[0])
			for _, o := range deployments.Items {
				json := make(map[string]interface{})
				json["name"] = o.Name
				json["replicasets"]=o.Status.Replicas
				json["availablereplicasets"]=o.Status.AvailableReplicas
				json["Namespace"] = o.Namespace
				json["createtime"] = o.CreationTimestamp
				//json["ClusterIP"] = o.Spec.ClusterIP
				//json["Ports"] = o.Spec.Ports
				arry = append(arry, json)
			}
			this.Data["json"] = arry
			this.ServeJSON()
		}
	}
	if label == "REPLICASETS" {
		deployments, err := clientset.AppsV1().ReplicaSets("").List(context.TODO(), metav1.ListOptions{})
		//beego.Info(deployments, "000000000000000000000000000")
		if err == nil {
			arry := [](map[string]interface{}){}
			for _, o := range deployments.Items {
				json := make(map[string]interface{})
				json["name"] = o.Name
				json["Namespace"] = o.Namespace
				json["createtime"] = o.CreationTimestamp
				//json["ClusterIP"] = o.Spec.ClusterIP
				//json["Ports"] = o.Spec.Ports
				arry = append(arry, json)
			}
			this.Data["json"] = arry
			this.ServeJSON()
		}
	}
}