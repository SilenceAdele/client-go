package main

import (
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
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

	// 2.实例化DynamicClient对象
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)

	//发送请求，获取gvr数据
	_, apiSource, err1 := discoveryClient.ServerGroupsAndResources()
	if err1 != nil {
		log.Println(err1)
	}

	for _, list := range apiSource {
		gv, err2 := schema.ParseGroupVersion(list.GroupVersion)
		if err2 != nil {
			log.Println(err2)
		}
		for _, resource := range list.APIResources {
			fmt.Printf("name: %v,group: %v,version: %v \n", resource.Name, gv.Group, gv.Version)
		}
	}
}
