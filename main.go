package main

import (
	"clientgo-learn/entity"
	"context"
	"encoding/json"
	"fmt"
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

}
