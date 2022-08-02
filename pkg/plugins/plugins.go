package plugins

import (
	"fmt"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

type Args struct {
	FavoriteColor  string `json:"favorite_color,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
	ThanksTo       string `json:"thanks_to,omitempty"`
}

type QinchyPlugin struct {
	args   *Args
	handle framework.FrameworkHandle
}

const (
	// Name 插件的名字
	Name = "Qinchy-Plugin"
)

// Name 返回插件的名字
func (s *QinchyPlugin) Name() string {
	return Name
}

func (s *QinchyPlugin) PreFilter(_ *framework.PluginContext, pod *core.Pod) *framework.Status {
	klog.V(3).Infof("prefilter pod: %v", pod.Name)
	return framework.NewStatus(framework.Success, "")
}

func (s *QinchyPlugin) Filter(_ *framework.PluginContext, pod *core.Pod, nodeName string) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v", pod.Name, nodeName)
	// 1.先判断是否混部节点和混部应用，只有混部应用才能调度到混部节点上
	nodeInfo := s.handle.NodeInfoSnapshot().NodeInfoMap[nodeName]
	// 是混部应用但是不是混部节点的情况
	if nodeInfo.Node().Labels["mixed-schedule"] != "true" && pod.Labels["cmos/mixed-pod"] == "true" {
		klog.V(3).Infof("filter pod: %v is mixed pod,but node: %v is not mixed-schedule node，marked as Unschedulable!", pod.Name, nodeName)
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("filter pod: %s is mixed pod,but node: %s is not mixed node，marked as Unschedulable!", pod.Name, nodeName))
	}

	annotations := pod.Annotations
	klog.V(3).Infof("filter pod: %s, annotations: %v", pod.Name, annotations)
	cpuIntensity, memIntensity := annotations["cpuIntensity"], annotations["memIntensity"]
	klog.V(3).Infof("filter pod: %s, cpuIntensity: %s, memIntensity: %s", pod.Name, cpuIntensity, memIntensity)

	// cpuQuantity 这个pod中所有容器的cpu请求量
	// memQuantity 这个pod中所有容器的mem请求量
	cpuRequestQuantity, memoryRequestQuantity := &resource.Quantity{}, &resource.Quantity{}
	for _, container := range pod.Spec.Containers {
		cpuRequestQuantity.Add(container.Resources.Limits["cmos/mixed-cpu"])
		memoryRequestQuantity.Add(container.Resources.Limits["cmos/mixed-memory"])
	}
	klog.V(3).Infof("filter pod: %v, total cpuQuantity: %v, total memQuantity: %v", pod.Name, cpuRequestQuantity, memoryRequestQuantity)

	nodeMixedCpuAmount := nodeInfo.Node().Labels["cmos/mixed-cpu"]
	cpuAvaliableQuantity, err := resource.ParseQuantity(nodeMixedCpuAmount)
	if err != nil {
		klog.V(3).Infof("parseQuantity node %v mixed-cpu resource: %s error", nodeName, nodeMixedCpuAmount)
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("parseQuantity node %v mixed-cpu resource: %s error", nodeName, nodeMixedCpuAmount))
	}

	// 如果cpu请求大于可混部的cpu量则不能调度
	if cpuRequestQuantity.Cmp(cpuAvaliableQuantity) == 1 {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("Pod %s Request cpu greater than node %v supplyed mixed-cpu resource", pod.Name, nodeName))
	}

	nodeMixedMemoryAmount := nodeInfo.Node().Labels["cmos/mixed-memory"]
	memoryAvaliableQuantity, err := resource.ParseQuantity(nodeMixedMemoryAmount)
	if err != nil {
		klog.V(3).Infof("parseQuantity node %v mixed-memory resource: %s error", nodeName, nodeMixedMemoryAmount)
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("parseQuantity node %v mixed-cpu resource: %s error", nodeName, nodeMixedMemoryAmount))
	}

	// 如果cpu请求大于可混部的cpu量则不能调度
	if memoryRequestQuantity.Cmp(memoryAvaliableQuantity) == 1 {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("Pod %s Request memory greater than node %v supplyed mixed-memory resource", pod.Name, nodeName))
	}

	return framework.NewStatus(framework.Success, "")
}

func (s *QinchyPlugin) PreBind(_ *framework.PluginContext, pod *core.Pod, nodeName string) *framework.Status {
	if nodeInfo, ok := s.handle.NodeInfoSnapshot().NodeInfoMap[nodeName]; !ok {
		return framework.NewStatus(framework.Error, fmt.Sprintf("prebind get node info error: %+v", nodeName))
	} else {
		klog.V(3).Infof("prebind pod %v to node : %v", pod.Name, nodeInfo.Node().Name)
		return framework.NewStatus(framework.Success, "")
	}
}

// New 初始化新的插件
func New(configuration *runtime.Unknown, f framework.FrameworkHandle) (framework.Plugin, error) {
	args := &Args{}
	if err := framework.DecodeInto(configuration, args); err != nil {
		return nil, err
	}
	klog.V(3).Infof("get plugin config args: %+v", args)
	return &QinchyPlugin{
		args:   args,
		handle: f,
	}, nil
}
