package main

import (
	"clientgo-learn/entity"
	"context"
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"time"
)

func main() {
	//log := logr.Logger{}
	clientset, kubeConfig := GetKubeClient(context.Background(),
		"/Users/zhualong/.kube/config")
	kc := entity.NewKubeController(kubeConfig, clientset, time.Second*30)
	stopPadch := make(chan struct{})
	kc.Run(stopPadch)
	go func() {
		kc.Run(stopPadch)
		//<-stopPadch
	}()
	for {
		if kc.Status == 1 {
			break
		}
		time.Sleep(time.Second * 1)
		fmt.Println("sleep 1S")
	}
	//client.CoreV1().Pods(namespace).List(context.TODO(), options)
	pod, err := kc.PodLister.Pods("default").
		Get("tomcat-d9c7df887-lml26")
	if err != nil {
		//log.Error(err, "get pods err")
		fmt.Errorf("get pods err")
	}
	fmt.Println("the pod hostname is: %s", pod.Spec.NodeName)
	b, err := json.Marshal(pod)
	fmt.Println("the pod ：%s", string(b))
	//log.Info("the pod hostname is: %s", pod.Spec.NodeName)
	container, err := json.Marshal(pod.Spec.Containers[0])
	fmt.Println("the pod hostname is: %s", string(container))

	//测试删除
	//pkc := &kubeclient.PodKubeController{kc}
	//err = pkc.Delete("default", "tomcat-d9c7df887-vl6sp")
	//if err != nil {
	//	fmt.Println("delete pod failed ：%s", err.Error())
	//}

	informer := kc.DeploymentInformer.Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deploy := obj.(*v1.Deployment)
			fmt.Println("add a deployment:", deploy.Name)
		},
		DeleteFunc: func(obj interface{}) {
			deploy := obj.(*v1.Deployment)
			fmt.Println("delete a deployment:", deploy.Name)
		},
		UpdateFunc: func(old, new interface{}) {
			oldDeploy := old.(*v1.Deployment)
			newDeploy := new.(*v1.Deployment)
			fmt.Println("update deployment:", oldDeploy.Name, newDeploy.Name)
		},
	})

	deploymentNamespaceLister := kc.DeploymentLister.Deployments("default")
	list, err := deploymentNamespaceLister.List(labels.Everything())
	if err != nil {
		return
	}

	for index, deployment := range list {
		bytes, err := json.Marshal(deployment)
		if err != nil {
			return
		}
		fmt.Println(fmt.Printf("顺序:%s---deployment:%s\n", index, string(bytes)))
	}

	pods, err := kc.PodLister.List(labels.Everything())
	if err != nil {
		return
	}

	for i, v := range pods {
		bytes, err := json.Marshal(v)
		if err != nil {
			return
		}
		fmt.Println(fmt.Printf("pod顺序:%i---deployment:%s\n", i, string(bytes)))
	}

	labeles := make(map[string]string, 1)
	labeles["control-plane"] = "controller-manager"
	selector := labels.SelectorFromSet(labeles)
	services, err := kc.ServiceInformer.Lister().Services("nodedemo-system").List(selector)
	if err != nil {
		return
	}
	for k, service := range services {
		bytes, err := json.Marshal(service)
		if err != nil {
			return
		}
		fmt.Println(fmt.Printf("service key:%i---service:%s\n", k, string(bytes)))
	}
}
