package main

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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
	dynamicClient, err2 := dynamic.NewForConfig(config)
	if err2 != nil {
		log.Println(err2)
	}

	// 3.配置需要调用的GVR
	gvr := schema.GroupVersionResource{
		Group:    "", //不需要写，因为是无名资源组，也就是core资源组
		Version:  "v1",
		Resource: "pods",
	}

	// 4.发送请求，且得到结果
	unStruceData, err3 := dynamicClient.Resource(gvr).Namespace("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err3 != nil {
		log.Println(err3)
	}
	podList := &corev1.PodList{}
	/*
		Rescorce,基于gvr生成一个针对于资源的客户端，也可以称之为动态客户端，dynamciRescorceClient
		NameSpace,指定一个可操作的命名空间。同时它是dynamciRescorceClient的方法
		List，首先是通过RESTClient调用K8s APIServer的接口返回了Pod的数据。返回的数据格式是二进制Json格式，
		然后通过一些了的解析方法，转换成UnstructuredList
	*/
	err4 := runtime.DefaultUnstructuredConverter.FromUnstructured(
		unStruceData.UnstructuredContent(),
		podList,
	)

	if err4 != nil {
		log.Println(err4)
	}

	for _, item := range podList.Items {
		fmt.Printf("namespace: %v,name: %v \n", item.Namespace, item.Name)
	}
}
