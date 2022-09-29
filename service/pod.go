package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/wonderivan/logger"
	"gok8s/config"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct {
}

//定义列表返回内容， items是pod的元素列表 total是元素数量
type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}

type PodsNp struct {
	Namespace string
	PodNum    int
}

//获取pod列表，支持过滤，排序，分页
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	//context.TODO()用于声明一个空的contex上下文，用于list方法内这个请求超时(源码)
	//metav1.ListOptions{} 用于过滤list数据，如使用label，field等
	podList, err := K8s.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取pod列表失败" + err.Error())
		return nil, errors.New("获取pod列表失败" + err.Error())
	}

	//实例化dataSelector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	//先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)

	//排序分页
	data := filtered.Sort().Paginate()

	//将dataCell类型转成pod
	pods := p.fromCells(data.GenericDataList)

	return &PodsResp{
		Total: total,
		Items: pods,
	}, nil
}

//获取pod详情
func (p *pod) GetDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取pod详情失败" + err.Error())
		return nil, errors.New("获取pod详情失败" + err.Error())
	}
	return pod, nil
}

//删除pod
func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.Clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除pod失败" + err.Error())
		return errors.New("删除pod失败" + err.Error())
	}
	return nil
}

//更新pod
func (p *pod) UpdatePod(namespace, content string) (err error) {
	var pod = &corev1.Pod{}
	//讲json反序列化为pod类型

	err = json.Unmarshal([]byte(content), pod)
	if err != nil {
		logger.Error("反序列化失败" + err.Error())
		return errors.New("反序列化失败" + err.Error())
	}
	//更新pod
	_, err = K8s.Clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新pod失败" + err.Error())
		return errors.New("更新pod失败" + err.Error())
	}
	return nil
}

//获取pod中的容器名列表
func (p *pod) GetPodContainer(podName, namespace string) (containers []string, err error) {
	//获取pod详情
	pod, err := p.GetDetail(podName, namespace)
	if err != nil {
		return nil, err
	}
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

//获取容器日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	//设置日志的配置，容器名，获取的内容的配置
	lineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}

	//获取一个request实例
	req := K8s.Clientset.CoreV1().Pods(namespace).GetLogs(podName, option)

	//发起stream连接，得到Response.body
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error("获取podLog失败" + err.Error())
		return "", errors.New("获取podLog失败" + err.Error())
	}
	defer podLogs.Close()
	//将response body写入到缓冲区，目的是为了转换成string类型
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error("复制podLog失败" + err.Error())
		return "", errors.New("复制podLog失败" + err.Error())
	}

	return buf.String(), nil
}

//类型转换的方法 covev1.pod -> DataCell,DataCell -> corev1.pod
func (p *pod) toCells(pods []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods {
		cells[i] = podCell(pods[i])
	}

	return cells
}

//获取每个namespace的pod数量
func (p *pod) GetPodNumPerNp() (podsNps []*PodsNp, err error) {
	//获取namespce
	namespaceList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		//获取pod列表
		podList, err := K8s.Clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		//组装数据
		podsNp := &PodsNp{
			Namespace: namespace.Name,
			PodNum:    len(podList.Items),
		}

		//添加到podsNps数组中
		podsNps = append(podsNps, podsNp)

	}
	return podsNps, nil
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range pods {
		//cells[i].(podCell) 是将DataCell类型转成podCell
		pods[i] = corev1.Pod(cells[i].(podCell))
	}

	return pods
}
