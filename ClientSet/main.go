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
	/*
		1.k8s的配置文件
		2.通过配置文件，使得程序能够链接到k8s集群
	*/

	// 1.加载config文件，生成config对象
	config, err := clientcmd.BuildConfigFromFlags("", "/home/hachi/.kube/config")
	if err != nil {
		log.Println(err)
	}

	//实例化ClientSet对象
	client, err2 := kubernetes.NewForConfig(config)
	if err2 != nil {
		log.Println(err2)
	}

	pods, err3 := client.
		CoreV1().                                  // 返回CoreV1Client
		Pods("kube-system").                       //指定查询的资源以及指定资源的namespace，如果为空，则表示查询所有的namespace
		List(context.TODO(), metav1.ListOptions{}) //这里表示查询的pods列表

		/*
			Corev1 返回 CoreV1Client实例对象
			Pods调用newPods函数，该函数返回的是PodInterface对象 PodInterface对象实现了Pods资源相关的全部方法，同时在newPods
			里面还将RESTClient实例对象赋值给对应的Client属性。
			List内使用RestClient与k8s APIServer 进行交互。
		*/
	if err3 != nil {
		log.Println(err3)
	}
	for _, item := range pods.Items {
		fmt.Printf("namesapce: %v,name: %v\n", item.Namespace, item.Name)
	}
}
