package main

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeClient(ctx context.Context, cfgpath string) (*kubernetes.Clientset, *rest.Config) {

	log := logr.FromContextOrDiscard(ctx)
	//配置文件的存放路径
	configfile := cfgpath
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", configfile)
	if err != nil {
		//log.Error(err, "BuildConfigFromFlags kube clientset faild")
		fmt.Errorf("BuildConfigFromFlags kube clientset faild")
		panic(err)
	}
	//生成clientset
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		log.Error(err, "clientset get faild")
		panic(err)
	}
	log.Info("GetKubeClient success")
	return clientset, kubeconfig

}
