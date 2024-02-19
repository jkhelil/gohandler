package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//ctx, _ := context.WithCancel(context.Background())
	var kubeclient kubernetes.Interface
	if kubeclient, err := loadKubernetesClientSet(); err != nil {
		fmt.Printf("Failed to load kubernetes API", err.Error())
		os.Exit(1)
	} else {
		kubeclient = kubeclient
	}
	shareInformerFactory := informers.NewSharedInformerFactoryWithOptions(kubeclient, time.Minute*5)
	podInformer := shareInformerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("Pod Added: %s\n", pod.Name)

		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newPod := newObj.(*corev1.Pod)
			fmt.Printf("Pod Updated: %s\n", newPod.Name)
		},

		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("Pod Deleted: %s\n", pod.Name)
		},
	})
	shareInformerFactory.Start(wait.NeverStop)
	for gvr, ok := range shareInformerFactory.WaitForCacheSync(wait.NeverStop) {
		if !ok {
			fmt.Errorf(fmt.Sprintf("Failed to sync cache for resource %v", gvr))
		}
	}
	fmt.Println("start kubelet informer...")
}

func GetConfig() (*rest.Config, error) {
	if len(os.Getenv("KUBECONFIG")) > 0 {
		return clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	}
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}
	if usr, err := user.Current(); err == nil {
		if c, err := clientcmd.BuildConfigFromFlags("", filepath.Join(usr.HomeDir, ".kube", "config")); err == nil {
			return c, nil
		}
	}
	return nil, fmt.Errorf("could not locate kubeconfig")
}

func loadKubernetesClientSet() (kubernetes.Interface, error) {
	kubeRestConfig, err := GetConfig()
	//fmt.Println(kubeRestConfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(kubeRestConfig)
}
