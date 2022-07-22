package plugins

import (
	"context"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

type Args struct {
	FavoriteColor  string `json:"favorite_color,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
	ThanksTo       string `json:"thanks_to,omitempty"`
}

type QinchyScheduler struct {
	handle framework.FrameworkHandle
}

var _ framework.PreFilterPlugin = &QinchyScheduler{}
var _ framework.FilterPlugin = &QinchyScheduler{}
var _ framework.ScorePlugin = &QinchyScheduler{}
var _ framework.PreBindPlugin = &QinchyScheduler{}

const (
	// Name 插件的名字
	Name = "Qinchy-Plugin"
)

// Name 返回插件的名字
func (s *QinchyScheduler) Name() string {
	return Name
}

//PreFilter 预过滤
func (s *QinchyScheduler) PreFilter(ctx context.Context, state *framework.CycleState, pod *core.Pod) *framework.Status {
	klog.V(3).Infof("PreFilter pod:%#v\n", pod)
	return framework.NewStatus(framework.Success, "预过滤通过")
}

func (s *QinchyScheduler) AddPod(ctx context.Context, state *framework.CycleState, podToSchedule *core.Pod, podToAdd *core.Pod, nodeInfo *nodeinfo.NodeInfo) *framework.Status {
	klog.V(3).Infof("PreFilterExtensions -> AddPod pod:%#v, SchedulePod pod:%#v, nodeInfo:%#v\n", podToAdd, podToSchedule, nodeInfo)
	return framework.NewStatus(framework.Success, "添加pod成功")
}

func (s *QinchyScheduler) RemovePod(ctx context.Context, state *framework.CycleState, podToSchedule *core.Pod, podToRemove *core.Pod, nodeInfo *nodeinfo.NodeInfo) *framework.Status {
	klog.V(3).Infof("PreFilterExtensions -> RemovePod pod:%#v, SchedulePod pod:%#v, nodeInfo:%#v\n", podToRemove, podToSchedule, nodeInfo)
	return framework.NewStatus(framework.Success, "移除pod成功")
}

//PreFilterExtensions 预过滤扩展器
func (s *QinchyScheduler) PreFilterExtensions() framework.PreFilterExtensions {
	return s
}

//Filter 过滤节点
func (s *QinchyScheduler) Filter(ctx context.Context, state *framework.CycleState, pod *core.Pod, nodeInfo *nodeinfo.NodeInfo) *framework.Status {
	klog.V(3).Infof("Filter pod:%#v\n", pod)
	return framework.NewStatus(framework.Success, "预过滤通过")
}

//Score 打分
func (s *QinchyScheduler) Score(ctx context.Context, state *framework.CycleState, pod *core.Pod, nodeName string) (int64, *framework.Status) {
	klog.V(3).Infof("Score pod:%s nodeName:%s\n", pod.Name, nodeName)
	return 10, nil
}

//NormalizeScore 归一化
func (s *QinchyScheduler) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *core.Pod, scores framework.NodeScoreList) *framework.Status {
	klog.V(3).Infof("NormalizeScore pod:%v\n", pod.Name)
	return framework.NewStatus(framework.Success, "预过滤通过")
}

//ScoreExtensions 打分扩展
func (s *QinchyScheduler) ScoreExtensions() framework.ScoreExtensions {
	return s
}

//PreBind 预绑定
func (s *QinchyScheduler) PreBind(ctx context.Context, state *framework.CycleState, pod *core.Pod, nodeName string) *framework.Status {
	if nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(nodeName); err != nil {
		klog.V(3).Infof("prebind get node info error: %v\n", nodeName)
		return framework.NewStatus(framework.Error, fmt.Sprintf("prebind get node info error: %v", nodeName))
	} else {
		klog.V(3).Infof("prebind node info: %v\n", nodeInfo.Node().Name)
		return framework.NewStatus(framework.Success, "")
	}
}

// New 初始化新的插件
func New(configuration *runtime.Unknown, f framework.FrameworkHandle) (framework.Plugin, error) {
	klog.V(3).Infof("使用配置：%#v初始化调度器：%s\n",configuration, Name)
	return &QinchyScheduler{
		handle: f,
	}, nil
}
