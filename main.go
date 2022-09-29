package main

import (

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
	"github.com/gin-gonic/gin"
	"gok8s/config"
	"gok8s/controller"
	"gok8s/service"
)

func main() {

	//初始化k8s client

	service.K8s.Init() //可以使用service.K8s.Clientset调用

	r := gin.Default()
	//跨包调用
	controller.Router.InitApiRouter(r)
	r.Run(config.ListenAddr)

}
