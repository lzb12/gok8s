package service

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct {
}

//定义列表的返回内容，Items是deployment元素列表，Total为deployment元素数量
type DeploymentResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

//定义DeployCreat结构体，用于创建deployment需要的参数属性的定义
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Imag          string            `json:"imag"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
}

//定义DeploysNP 用于返回namespace中deployment的数量
type DeploysNp struct {
	Namespace string `json:"namespace"`
	DeployNum int    `json:"deploy_num"`
}

//获取deployment列表，支持过滤，排序，分页
func (d *deployment) GetDeployment(filterName, namespace string, limit, page int) (deploymentResp *DeploymentResp, err error) {

	deploymentList, err := K8s.Clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取deployment列表失败" + err.Error())
		return nil, errors.New("获取deployment列表失败" + err.Error())
	}
	//将deployment中的deployment列表，放进dataselector对象中进行排序
	selectableData := &dataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
		DataSelect: &DataSelectQuery{
			Filter:   &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{Limit: limit, Page: page},
		},
	}
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()

	//将[]DataCell类型的deployment列表转为appsv1.deployment列表
	deployments := d.fromCells(data.GenericDataList)

	return &DeploymentResp{
		Items: deployments,
		Total: total,
	}, nil

}

//获取Deployment详情
func (d *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployments *appsv1.Deployment, err error) {
	deployments, err = K8s.Clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取deployment失败" + err.Error())
		return nil, errors.New("获取deployment失败" + err.Error())
	}
	return deployments, nil
}

//修改Deployment副本数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replica int32, err error) {
	//获取 autoscalingv1.Scale类型的对象，能点出当前的副本数
	scale, err := K8s.Clientset.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Deployment副本数信息失败" + err.Error())
		return 0, errors.New("获取Deployment副本数信息失败" + err.Error())
	}

	//修改副本数
	scale.Spec.Replicas = int32(scaleNum)

	//更新副本数传入scale对象
	newScale, err := K8s.Clientset.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Deployment副本数信息失败" + err.Error())
		return 0, errors.New("更新Deployment副本数信息失败" + err.Error())
	}
	return newScale.Spec.Replicas, nil

}

//创建Deployment

func (d *deployment) CreateDeployment(data *DeployCreate) (err error) {
	//初始化APPSv1.Deployment类型对象，并将入参的data数据放进去
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},

		//spec中定义副本数，选择器以及pod属性
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				//定义pod名和标签
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				//定义容器名，镜像和端口
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Imag,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: data.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		//status定义资源运行状态，这里由于是新建，传入空的appsv1.DeploymentStatus{}对象即可
		Status: appsv1.DeploymentStatus{},
	}
	//判断健康检查功能是否打开，若打开，则增加健康检查功能
	if data.HealthCheck {
		//设置第一个容器的ReadinessProbe，因为我们pod中只有一个容器，所以直接使用index 0 即可
		//若pod中有多个容器，则这里需要使用for循环去定义了
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
					//Type: 0则表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
					//Type: 1则表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			//初始化等待时间
			InitialDelaySeconds: 15,
			//超时时间
			TimeoutSeconds: 5,
			//执行间隔
			PeriodSeconds: 5,
		}

	}
	//定义容器的limit和request资源
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	//调用sdk更新deployment
	_, err = K8s.Clientset.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建deployment失败" + err.Error())
		return errors.New("创建deployment失败" + err.Error())
	}
	return nil
}

//删除Deployment

func (d *deployment) DeleteDeploy(deploymentName, namespace string) (err error) {

	err = K8s.Clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除deployment失败" + err.Error())
		return errors.New("删除deployment失败" + err.Error())
	}
	return nil
}

//重启Deployment
func (d *deployment) RestartDeployment(deploymentName, namespace string) (err error) {

	//使用patchData Map组装数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{"name": deploymentName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}
	//序列化为字节，因为patch方法只接收字节类型参数
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error("json序列化失败" + err.Error())
		return errors.New("json序列化失败" + err.Error())
	}
	//调用patch方法更新deployment
	_, err = K8s.Clientset.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error("重启deployment失败" + err.Error())
		return errors.New("重启deployment失败" + err.Error())
	}

	return nil

}

//更新deployment
func (d *deployment) UpdateDeployment(namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}

	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error("反序列化失败" + err.Error())
		return errors.New("反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新deployment失败" + err.Error())
		return errors.New("更新deployment失败" + err.Error())
	}

	return nil
}

//获取每个namespace的Deployment

func (d *deployment) GetDeployNumPerNp() (deploysNps []*DeploysNp, err error) {
	namespaceList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		deploymentList, err := K8s.Clientset.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		deploysNp := &DeploysNp{
			Namespace: namespace.Name,
			DeployNum: len(deploymentList.Items),
		}

		deploysNps = append(deploysNps, deploysNp)

	}
	return deploysNps, nil
}

//类型转换
func (d *deployment) toCells(deployments []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(deployments))
	for i := range deployments {
		cells[i] = deploymentCell(deployments[i])
	}
	return cells
}

//类型转换
func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	pods := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		pods[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return pods
}
