package kubeclient

import (
	"clientgo-learn/entity"
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
)
import corev1 "k8s.io/api/core/v1"
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type PodKubeController struct {
	*entity.KubeController
}

func (pkc *PodKubeController) Get(namespace string, name string) (interface{}, error) {
	return pkc.PodLister.Pods(namespace).Get(namespace)
}

func (pkc *PodKubeController) Create(namespace string, name string, kubePod interface{}) (interface{}, error) {
	pod := kubePod.(*corev1.Pod)
	return pkc.Clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
}

func (pkc *PodKubeController) Update(namespace string, name string, kubePod interface{}) (interface{}, error) {
	pod := kubePod.(*corev1.Pod)
	return pkc.Clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
}

func (pkc *PodKubeController) GetAllInfo() (interface{}, error) {
	return pkc.DeploymentLister.List(labels.Everything())
}

func (pkc *PodKubeController) GetFromLabelApp(namespace string, appName string) (interface{}, error) {
	return pkc.PodLister.Pods(namespace).List(labels.SelectorFromSet(map[string]string{"app": appName}))
}

func (pkc *PodKubeController) Delete(namespace string, podName string) error {
	return pkc.Clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
}

func (pkc *PodKubeController) OffOnline(namespace string, name string, kubePod interface{}) (interface{}, error) {

	pod, ok := kubePod.(*corev1.Pod)
	if ok {
		//对pod进行深拷贝操作
		podCopy := pod.DeepCopy()
		if podCopy.ObjectMeta.Labels["online"] == "true" {
			podCopy.ObjectMeta.Labels["online"] = "false"
		}
		//更新操作
		podUpdated, err := pkc.Clientset.CoreV1().Pods(namespace).
			Update(context.TODO(), podCopy, metav1.UpdateOptions{})
		if err != nil {
			fmt.Printf("update pod error:%s", err)
			panic(err)
		}
		return podUpdated, err
	}

	return nil, nil

}
