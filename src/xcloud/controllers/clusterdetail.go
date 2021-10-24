package controllers

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//restclient "k8s.io/client-go/rest"
	"log"
	//"encoding/json"
	//"xcloud/models"
	"xcloud/util"
	//"fmt"
	"github.com/astaxie/beego"
	//"github.com/beego/beego/v2/client/orm"


	//"k8s.io/client-go/util/retry"
)

// Operations about object
type ClusterdetailController struct {
	beego.Controller
}

var clientPool = util.Lock{}

// @Title Get
// @Description get cluster by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Cluster
// @Failure 403 :uid is empty
// @router /:clustername [get]
func (c *ClusterdetailController) Get() {
	name := c.GetString(":clustername")
	//beego.Info(name,"--------------------------------")
	//o := orm.NewOrm()
	//var cluster models.Cluster
  	//o.QueryTable("Cluster").Filter("name", name).One(&cluster)
	//config := restclient.Config{}
	//config.CAData=[]byte(cluster.Cacrt)
	//tlsCfg := restclient.TLSClientConfig{
	//	Insecure:false,
	//	CAData:[]byte(cluster.Cacrt),
	//	KeyData:[]byte(cluster.Privitekey),
	//	CertData:[]byte(cluster.Publickey),
	//}
	//config.TLSClientConfig = tlsCfg
	//var port string = cluster.Apiserver_port
	//var master string = cluster.Apiserver_ip
	//config.Host = "https://" + master + ":" + port
	//beego.Info(config,"8888888888--------------------------------")

    clientset := util.Getclient( name)
	//获取POD
	pods,err := clientset.CoreV1().Pods("").List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	//fmt.Println(len(pods.Items))
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(),metav1.ListOptions{})
	//fmt.Println("##################")
	services, err := clientset.CoreV1().Services("").List(context.TODO(),metav1.ListOptions{})
	configmaps, err := clientset.CoreV1().ConfigMaps("").List(context.TODO(),metav1.ListOptions{})
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(context.TODO(),metav1.ListOptions{})
	deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(),metav1.ListOptions{})
	source_count := make(map[string]int)
	source_count["pods"] = len(pods.Items)
	source_count["nodes"] = len(nodes.Items)
	source_count["services"] = len(services.Items)
	source_count["pvcs"] = len(pvcs.Items)
	source_count["configmaps"] = len(configmaps.Items)
	source_count["deployments"] = len(deployments.Items)
    beego.Info(source_count, "###############################################")
	c.Data["json"] = source_count
	c.ServeJSON()
}
