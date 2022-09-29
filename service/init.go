package service

import (
	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//用于初始化k8s clientset
var K8s k8s

type k8s struct {
	Clientset *kubernetes.Clientset
}

func (k *k8s) Init() {
	configPath := "config/config"
	//将kubeconfig文件转换成rest.config类型的对象
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		//fmt.Println("连接失败", err)
		panic("获取k8s client 配置失败" + err.Error())
	}

	// 根据rest.config类型的对象，newyigeclientset出来
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("创建k8s client失败" + err.Error())
	} else {
		logger.Info("k8s client 初始化成功！")
	}

	k.Clientset = clientset
}
