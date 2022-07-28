package plugins

import (
	"fmt"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

type Args struct {
	FavoriteColor  string `json:"favorite_color,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
	ThanksTo       string `json:"thanks_to,omitempty"`
}

type QinchyScheduler struct {
	args   *Args
	handle framework.FrameworkHandle
}

const (
	// Name 插件的名字
	Name = "Qinchy-Plugin"
)

// Name 返回插件的名字
func (s *QinchyScheduler) Name() string {
	return Name
}

func (s *QinchyScheduler) PreFilter(pc *framework.PluginContext, pod *core.Pod) *framework.Status {
	klog.V(3).Infof("prefilter pod: %v", pod.Name)
	return framework.NewStatus(framework.Success, "")
}

func (s *QinchyScheduler) Filter(pc *framework.PluginContext, pod *core.Pod, nodeName string) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v", pod.Name, nodeName)
	return framework.NewStatus(framework.Success, "")
}

func (s *QinchyScheduler) PreBind(pc *framework.PluginContext, pod *core.Pod, nodeName string) *framework.Status {
	if nodeInfo, ok := s.handle.NodeInfoSnapshot().NodeInfoMap[nodeName]; !ok {
		return framework.NewStatus(framework.Error, fmt.Sprintf("prebind get node info error: %+v", nodeName))
	} else {
		klog.V(3).Infof("prebind node info: %+v", nodeInfo.Node())
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
	return &QinchyScheduler{
		args:   args,
		handle: f,
	}, nil
}
