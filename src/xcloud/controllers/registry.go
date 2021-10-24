package controllers

import (
	"context"
	"encoding/json"
	"k8s.io/apimachinery/pkg/api/errors"
	//"fmt"
	//"reflect"
	"strconv"
	"k8s.io/apimachinery/pkg/util/intstr"
	//appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//restclient "k8s.io/client-go/rest"
	"xcloud/models"
	v1 "k8s.io/api/apps/v1"
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/beego/beego/v2/client/orm"
	"xcloud/util"
	//"k8s.io/apimachinery/pkg/runtime/schema"
)

// Operations about object
type RegistryController struct {
	beego.Controller
}



// @Title Create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (this *RegistryController) Post() {
	var registrydata map[string]interface{}
	json.Unmarshal([]byte(this.Ctx.Input.RequestBody), &registrydata)
	beego.Info(registrydata,"==============99999===============")
    //获取连接信息并连接集群
	o := orm.NewOrm()
	var cluster models.Cluster
	clustervalue := registrydata["clustervalue"].(float64)
	o.QueryTable("Cluster").Filter("Id", clustervalue).One(&cluster)
	clustername := cluster.Name
    registryname := registrydata["name"].(string)
	clientset := util.Getclient(clustername)
	//查看configmap是否已经存在
	beego.Info("start")
	_, err1 := clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), registryname + "-config", metav1.GetOptions{})
	if !errors.IsNotFound(err1) {
		this.CustomAbort(500 ,"ConfigMap已存在")
	}
	//创建secret组件
	if( registrydata["ssl"] == true && registrydata["cert_id"] != nil ){
		certid := registrydata["cert_id"].(float64)
		var cert models.Cert
		o.QueryTable("Cert").Filter("Id", certid).One(&cert)
		certCrt :=  cert.Crt
		certKey := cert.Key
		secret_yaml := &apiv1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind: "secret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: registryname + "-cert",
			},

			Data: map[string][]byte{
				"registry-crt": []byte(certCrt),
				"registry-key": []byte(certKey),
			},
		}
		_, err2 := clientset.CoreV1().Secrets("default").Create(context.TODO(), secret_yaml, metav1.CreateOptions{})
		//beego.Info(err2.Error(), "iiiiiiiiiiiiisceret")
		if err2 != nil {
			this.CustomAbort(500 ,err2.Error())
		}
		beego.Info(err2, "iiiiiiiiiiiiisceret")
	}
	//创建configmap对象
    if( registrydata["ssl"] == false ){
    	nossl_registry_yaml := `version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
`
		configmap_nossl_yaml := &apiv1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind: "configmap",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: registryname + "-config",
			},
			Data: map[string]string{
				"config.yml": nossl_registry_yaml,
			},
		}
		_, err3 := clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), configmap_nossl_yaml, metav1.CreateOptions{})
		beego.Info(err3, "iiiiiiiiiiiiiconfigmap")
		if err3 != nil {
			this.CustomAbort(500 ,err3.Error())
		}
	}else{
		ssl_registry_yaml := `version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
  tls:
    certificate: /opt/registry-crt
    key: /opt/registry-key
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
`
		configmap_ssl_yaml := &apiv1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind: "configmap",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: registryname + "-config",
			},
			Data: map[string]string{
				"config.yml": ssl_registry_yaml,
			},
		}
		_, err3_2 := clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), configmap_ssl_yaml, metav1.CreateOptions{})
		beego.Info(err3_2, "iiiiiiiiiiiiiconfigmap")
		if err3_2 != nil {
			this.CustomAbort(500 ,err3_2.Error())
		}
	}
	if (registrydata["ssl"] == true) {
		deployment := &v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: registryname,
			},
			Spec: v1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": registryname,
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":     registryname,
							"version": "V1",
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:            registryname,
								Image:           "registry:2",
								ImagePullPolicy: "IfNotPresent",
								Ports: []apiv1.ContainerPort{
									{
										Name:          "http",
										Protocol:      apiv1.ProtocolTCP,
										ContainerPort: 80,
									},
								},
								VolumeMounts: []apiv1.VolumeMount{
									{
										Name:      "secret",
										MountPath: "/opt/",
									},
									{
										Name:      "config",
										MountPath: "/etc/docker/registry",
									},
									{
										Name:      "image",
										MountPath: "/var/lib/registry",
									},
								},
							},
						},
						Volumes: []apiv1.Volume {
							{
								Name: "secret",
								VolumeSource: apiv1.VolumeSource{
									Secret: &apiv1.SecretVolumeSource{
										SecretName: registryname + "-cert",
									},
								},
							},
							{
								Name: "config",
								VolumeSource: apiv1.VolumeSource{
									ConfigMap: &apiv1.ConfigMapVolumeSource{
										LocalObjectReference : apiv1.LocalObjectReference{
											Name: registryname + "-config",
										},
									},
								},
							},
							{
								Name: "image",
								VolumeSource: apiv1.VolumeSource{
									HostPath: &apiv1.HostPathVolumeSource{
										Path: "/var/lib/registry",
									},
								},
							},
						},
					},
				},
			},
		}
		_, err4_1 := clientset.AppsV1().Deployments(apiv1.NamespaceDefault).Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err4_1 != nil {
			this.CustomAbort(500 ,err4_1.Error())
		}
	}else{
		deployment := &v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: registryname,
			},
			Spec: v1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": registryname,
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":     registryname,
							"version": "V1",
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:            registryname,
								Image:           "registry:2",
								ImagePullPolicy: "IfNotPresent",
								Ports: []apiv1.ContainerPort{
									{
										Name:          "http",
										Protocol:      apiv1.ProtocolTCP,
										ContainerPort: 80,
									},
								},
								VolumeMounts: []apiv1.VolumeMount{
									{
										Name:      "config",
										MountPath: "/etc/docker/registry",
									},
									{
										Name:      "image",
										MountPath: "/var/lib/registry",
									},
								},
							},
						},
						Volumes: []apiv1.Volume {
							{
								Name: "config",
								VolumeSource: apiv1.VolumeSource{
									ConfigMap: &apiv1.ConfigMapVolumeSource{
										LocalObjectReference : apiv1.LocalObjectReference{
											Name: registryname + "-config",
										},
									},
								},
							},
							{
								Name: "image",
								VolumeSource: apiv1.VolumeSource{
									HostPath: &apiv1.HostPathVolumeSource{
										Path: "/var/lib/registry",
									},
								},
							},
						},
					},
				},
			},
		}
		_, err4_2 := clientset.AppsV1().Deployments(apiv1.NamespaceDefault).Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err4_2 != nil {
			this.CustomAbort(500 ,err4_2.Error())
		}
	}
	service_yaml := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: registryname,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": registryname,
			},
			Type: "NodePort",
			Ports: []apiv1.ServicePort{
				{
					Name: "registry",
					Port:5000,
					NodePort:31000,
					TargetPort: intstr.IntOrString{
						Type: intstr.Int,
						IntVal: 5000,
					},
					Protocol: apiv1.ProtocolTCP,
				},
			},
		},
	}
	_, err5 := clientset.CoreV1().Services("default").Create(context.TODO(), service_yaml, metav1.CreateOptions{})
	beego.Info(err5, "iiiiiiiiiiiii")
	if err5 != nil {
		this.CustomAbort(500 ,err5.Error())
	}
	//数据库插入registry信息数据
	o1 := orm.NewOrm()
	registry := new(models.Registry)
	registry.Name = registrydata["name"].(string)
	registry.Domain = registrydata["domain"].(string)
	registry.Hostpath = registrydata["path"].(string)
	if(registrydata["ssl"] == true) {
		registry.IsSsl = true
	}
	if(registrydata["ssl"] == false) {
		registry.IsSsl = false
	}
	cluster_id := int(registrydata["clustervalue"].(float64))
	var clsr models.Cluster
	clsr.Id=cluster_id
	o1.Read(&clsr, "Id" )
	registry.Cluster = &clsr
	o2 := orm.NewOrm()
	env_id :=int(registrydata["envvalue"].(float64))
	var en models.Env
	en.Id = env_id
	o2.Read(&en, "Id")
	registry.Env = &en
	o3 := orm.NewOrm()
	if (registrydata["cert_id"] != nil) {
		cert_id := int(registrydata["cert_id"].(float64))
		var ce models.Cert
		ce.Id = cert_id
		o3.Read(&ce, "Id")
		registry.Cert = &ce
		o4 := orm.NewOrm()
		o4.Insert(registry)
	}
	o4 := orm.NewOrm()
	o4.Insert(registry)
	this.ServeJSON()
}

func int32Ptr(i int32) *int32 { return &i }

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (this *RegistryController) GetAll() {
	o := orm.NewOrm()
	var reg []*models.Registry
	o.QueryTable("Registry").RelatedSel().All(&reg)
	for _,v := range reg {
		o.LoadRelated(v, "Cluster")
		o.LoadRelated(v, "Env")
		o.LoadRelated(v, "Cert")
	}

	this.Data["json"] = reg
	this.ServeJSON()
}

// @Title Delete
// @Description delete registry
// @Success 200 {object} models.Registry
// @router /:id [delete]
func (this *RegistryController) Delete() {
	id,err := strconv.Atoi(this.GetString(":id"))
	if (err != nil) {
		beego.Info(err.Error())
	}
	registry_obj := models.Registry{}
	var o = orm.NewOrm()
	o.QueryTable("Registry").Filter("Id", id).RelatedSel().One(&registry_obj)
	registryname := registry_obj.Name
	clustername := registry_obj.Cluster.Name
	clientset := util.Getclient(clustername)
	clientset.CoreV1().ConfigMaps("default").Delete(context.TODO(), registryname + "-config", metav1.DeleteOptions{})
	clientset.AppsV1().Deployments(apiv1.NamespaceDefault).Delete(context.TODO(), registryname, metav1.DeleteOptions{})
	clientset.CoreV1().Services("default").Delete(context.TODO(), registryname, metav1.DeleteOptions{})
	clientset.CoreV1().Secrets("default").Delete(context.TODO(), registryname + "-cert", metav1.DeleteOptions{})
	if _, err := o.Delete(&models.Registry{Id: id}); err != nil {
		beego.Info(err.Error())
	}
	this.ServeJSON()
}