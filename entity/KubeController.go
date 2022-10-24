package entity

import (
	"fmt"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"time"
)

type KubeController struct {
	KubeConfig *rest.Config
	Status     int32
	ClusterId  []string
	Env        []string
	Clientset  *kubernetes.Clientset
	Factory    informers.SharedInformerFactory

	//定义Deployment资源对象的Informer、Lister以及HasSynce
	deploymentInformer appsinformers.DeploymentInformer
	deploymentLister   appslisters.DeploymentLister
	deploymentSynced   cache.InformerSynced

	//定义Pod资源对象的Informer、Lister以及HashSynce
	PodInformer coreinformers.PodInformer
	PodLister   corelisters.PodLister
	PodSynced   cache.InformerSynced

	//定义Service资源对象的Informer、Lister以及HashSynce
	ServiceInformer coreinformers.ServiceInformer
	ServiceLister   corelisters.ServiceLister
	ServiceSynced   cache.InformerSynced
}

func NewKubeController(kubeConfig *rest.Config,
	clientset *kubernetes.Clientset, defaultResync time.Duration) *KubeController {

	kc := &KubeController{KubeConfig: kubeConfig, Clientset: clientset}

	//通过clientset生成SharedInformerFactory
	//defaultResync参数可以控制reflector调用List的周期，如果设置为0，启动后获取
	//(List)一次全量的资源对象并存入缓存，后续不会再同步
	kc.Factory = informers.NewSharedInformerFactory(clientset, defaultResync)
	//生成Deployment、Pod、Service等资源对象的Informer、Lister以及HasSysnced

	kc.deploymentInformer = kc.Factory.Apps().V1().Deployments()
	kc.deploymentLister = kc.deploymentInformer.Lister()
	kc.deploymentSynced = kc.deploymentInformer.Informer().HasSynced

	kc.PodInformer = kc.Factory.Core().V1().Pods()
	kc.PodLister = kc.PodInformer.Lister()
	kc.PodSynced = kc.PodInformer.Informer().HasSynced

	kc.ServiceInformer = kc.Factory.Core().V1().Services()
	kc.ServiceLister = kc.ServiceInformer.Lister()
	kc.ServiceSynced = kc.ServiceInformer.Informer().HasSynced

	return kc

}
func (kc *KubeController) Run(stopPodch chan struct{}) {
	//log := logr.Logger{}
	//defer close(stopPodCh)'
	defer utilruntime.HandleCrash()
	//defer log.Error(nil, "KubeController shutdown")
	defer fmt.Errorf("KubeController shutdown")

	//传入停止的stopCh
	kc.Factory.Start(stopPodch)

	//等待资源查询(List)完成后同步到缓存
	if !cache.WaitForCacheSync(stopPodch, kc.deploymentSynced, kc.PodSynced, kc.ServiceSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for kuberesource caches to sync"))
		return
	}
	//同步成功，设置标志位 为1
	kc.Status = 1
	//log.Info("KubeController start")
	fmt.Println("KubeController start")
	//<-stopPodch
}
