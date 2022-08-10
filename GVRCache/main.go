package main

import (
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1.获取配置文件
	config, err := clientcmd.BuildConfigFromFlags("", "/home/hachi/.kube/config")
	if err != nil {
		log.Println(err)
	}

	// 2.实例化客户端，本客户端将负责将GVR数据缓存到客户端
	cacheDiscoveryClient, err2 := disk.NewCachedDiscoveryClientForConfig(config, "./cache/discovery", "./cache/http", time.Minute*60)
	if err2 != nil {
		log.Println(err2)
	}

	/*
		1.先从缓存文件中找到GVR数据，有则直接返回，否则调用apiserver
		2.调用apiserver，获取gvr数据
		3.将获取的gvr数据缓存到本地，然后返回到客户端
	*/
	_, apiSource, err3 := cacheDiscoveryClient.ServerGroupsAndResources()

	if err3 != nil {
		log.Println(err3)
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
