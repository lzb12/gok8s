package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gok8s/service"
	"net/http"
)

var Deployment deployment

type deployment struct {
}

//获取deployment列表，支持过滤，排序，分页
func (d *deployment) GetDeployments(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("bind请求参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	data, err := service.Deployment.GetDeployment(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logger.Error("获取deployment列表失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取deployment列表成功",
		"data": data,
	})
}

//获取deployment详情
func (d *deployment) GetDeploymentDetail(ctx *gin.Context) {
	params := new(struct {
		DeploymentName string `form:"deployment_name"`
		Namespace      string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("bind请求参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Deployment.GetDeploymentDetail(params.DeploymentName, params.Namespace)
	if err != nil {
		logger.Error("获取deployment详情失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取deployment详情成功",
		"data": data,
	})
}

//创建deployment
func (d *deployment) CreateDeployment(ctx *gin.Context) {
	params := new(service.DeployCreate)
	if err := ctx.Bind(params); err != nil {
		logger.Error("bind请求参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Deployment.CreateDeployment(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建deployment成功",
		"data": nil,
	})

}

//删除Deployment

func (d *deployment) DeleteDeploy(ctx *gin.Context) {
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
	})

	err := ctx.ShouldBindJSON(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "ShouldBindJson绑定失败",
			"data": nil,
		})
	}
	err = service.Deployment.DeleteDeploy(params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "删除deployment失败",
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除deployment成功",
		"data": nil,
	})
}

//重启Deployment
func (d *deployment) RestartDeployment(ctx *gin.Context) {
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
	})

	err := ctx.ShouldBindJSON(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "ShouldBindJson绑定失败",
			"data": nil,
		})
	}
	err = service.Deployment.RestartDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "重启deployment失败",
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "重启deployment成功",
		"data": nil,
	})

}

//更新deployment
func (d *deployment) UpdateDeployment(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})

	err := ctx.ShouldBindJSON(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "ShouldBindJson绑定失败",
			"data": nil,
		})
	}
	err = service.Deployment.UpdateDeployment(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "更新deployment失败",
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新deployment成功",
		"data": nil,
	})
}
