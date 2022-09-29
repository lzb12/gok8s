package controller

import "github.com/gin-gonic/gin"

var Router router

type router struct {
}

func (r *router) InitApiRouter(router *gin.Engine) {

	//pod操作
	pod := router.Group("/api/k8s/pod")
	pod.GET("", Pod.GetPods)
	pod.GET("/detail", Pod.GetPodDetail)
	pod.DELETE("/del", Pod.DeletePod)
	pod.PUT("/update", Pod.UpdatePod)
	pod.GET("/container", Pod.GetPodContainer)
	pod.GET("/log", Pod.GetPodLog)
	pod.GET("/num", Pod.GetPodNumPerNp)

	//deployment操作
	dep := router.Group("/api/k8s/deployment")
	dep.GET("", Deployment.GetDeployments)
	dep.GET("/detail", Deployment.GetDeploymentDetail)
	dep.POST("/create", Deployment.CreateDeployment)
	dep.DELETE("/del", Deployment.DeleteDeploy)
	dep.PUT("/update", Deployment.UpdateDeployment)

}
