package util

import (
	"fmt"
	"log"
	//"time"
	"github.com/beego/beego/v2/client/orm"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	//"k8s.io/client-go/dynamic"
	//"k8s.io/apimachinery/pkg/runtime/schema"
	"xcloud/models"
)

// client信息缓存
var clientPool = Lock{}
func Getclient(clustername string) (kubernetes.Clientset) {
	key := clustername + "clientSet"
	cl, ok := clientPool.Get(key)
	if ok && cl != nil {
		return cl.(kubernetes.Clientset)
	} else {
		//var kubeconfig *string
		//if home := homedir.HomeDir(); home != "" {
		//	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		//} else {
		//	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		//}
		//flag.Parse()
		//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		//if err != nil {
		//	log.Println(err)
		//}
		//var config interface{}
		config := Gettnlsconfig(clustername)
		clientset, err := kubernetes.NewForConfig(&config)
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println("connect k8s success")
		}
		clientPool.Put(key, *clientset)
		return *clientset
	}


}

//func Getyamlclient(clustername string, groups string, version string, api string ) interface{} {
//	key := clustername + "yamlclient"
//	cl, ok := clientPool.Get(key)
//	if ok && cl != nil {
//		return cl.(kubernetes.Clientset)
//	} else {
//		config := Gettnlsconfig(clustername)
//		gv := &schema.GroupVersion{groups, version}
//		config.ContentConfig = restclient.ContentConfig{GroupVersion: gv}
//		config.Timeout = time.Second * 3
//		config.APIPath = api
//		//创建新的dynamic client
//		cl, err := dynamic.NewForConfig(&config)
//		if err != nil {
//			log.Fatalln(err)
//		} else {
//			fmt.Println("connect k8s success")
//		}
//		clientPool.Put(key, cl)
//		return cl
//	}
//}

func Gettnlsconfig(name string) restclient.Config {
	o := orm.NewOrm()
	var cluster models.Cluster
	o.QueryTable("Cluster").Filter("name", name).One(&cluster)
	config := restclient.Config{}
	config.CAData=[]byte(cluster.Cacrt)
	tlsCfg := restclient.TLSClientConfig{
		Insecure:false,
		CAData:[]byte(cluster.Cacrt),
		KeyData:[]byte(cluster.Privitekey),
		CertData:[]byte(cluster.Publickey),
	}
	config.TLSClientConfig = tlsCfg
	var port string = cluster.Apiserver_port
	var master string = cluster.Apiserver_ip
	config.Host = "https://" + master + ":" + port
	return config
}