package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//加载配置，生成config对象
	config, err := clientcmd.BuildConfigFromFlags("", "/home/hachi/.kube/config")
	if err != nil {
		log.Println(config)
	}

	//获取clientSet对象
	clientSet, err2 := kubernetes.NewForConfig(config)
	if err2 != nil {
		log.Println(err2)
	}

	// 3.调用监听的方法
	watch, err3 := clientSet.AppsV1().Deployments("default").Watch(context.TODO(), metav1.ListOptions{})
	if err3 != nil {
		log.Println(err3)
	}

	fmt.Println("start...")
	for {
		select {
		case e, _ := <-watch.ResultChan():
			fmt.Println(e.Type, e.Object)
			//e.Type表示时间变化的类型，如add、delete、update等
			// e.Object:表示变化后的数据
		}
	}
}
