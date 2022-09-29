package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gok8s/service"
	"net/http"
)

var Pod pod

type pod struct {
}

//获取pod列表分页过滤排序
func (p *pod) GetPods(ctx *gin.Context) {
	//处理入参
	//匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
	})

}

//获取pod详情
func (p *pod) GetPodDetail(ctx *gin.Context) {
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod详情成功",
		"data": data,
	})

}

//删除pod
func (p *pod) DeletePod(ctx *gin.Context) {
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除pod成功",
		"data": nil,
	})

}

//更新pod
func (p *pod) UpdatePod(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBindJSON绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "ShouldBindJSON绑定失败" + err.Error(),
			"data": nil,
		})
		return
	}
	fmt.Println(params.Content)
	err := service.Pod.UpdatePod(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新pod成功",
		"data": nil,
	})

}

//获取pod中的容器名列表
func (p *pod) GetPodContainer(ctx *gin.Context) {
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodContainer(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取容器名列表失败",
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取容器名列表成功",
		"data": data,
	})

}

//获取容器日志
func (p *pod) GetPodLog(ctx *gin.Context) {

	params := new(struct {
		ContainerName string `form:"container_name"`
		PodName       string `form:"pod_name"`
		Namespace     string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取容器日志失败",
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取容器日志成功",
		"data": data,
	})

}

//获取每个namespace的pod数量

func (p *pod) GetPodNumPerNp(ctx *gin.Context) {
	data, err := service.Pod.GetPodNumPerNp()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod数量成功",
		"data": data,
	})

}
