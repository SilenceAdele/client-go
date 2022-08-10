package main

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
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

	//配置api路径  无组名资源组
	config.APIPath = "api" //api/v1/pods    Resource:开头字母小写，以s结尾；如pods Kind:以大写字母开头，结尾没有s； 如：Pod

	//配置分组版本
	config.GroupVersion = &corev1.SchemeGroupVersion //无组名资源组  Group=nil Version="v1"

	//配置数据的编码工具
	config.NegotiatedSerializer = scheme.Codecs

	//实例化RESTClient对象
	restclient, err2 := rest.RESTClientFor(config)
	if err2 != nil {
		log.Println(err2)
	}

	//定义接受返回值的变量
	result := &corev1.PodList{}

	//restClient与apiServer交互  通过k8s的api可以获取具体是哪种请求

	err3 := restclient.
		Get().                                                         //get请求方式
		Namespace("kube-system").                                      //指定命名空间
		Resource("pods").                                              //指定需要查询的资源，传递资源名称
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec). //参数及参数的序列化工具
		Do(context.TODO()).                                            //触发请求
		Into(result)                                                   //写入返回结果
	if err3 != nil {
		log.Println(err3)
	}
	/*
		get,定义请求方式，返回一个Request结构体对象。这个Request结构体对象，就是构建访问apiserver的请求的用的
		依次执行 Namespace(),Resource(),VersionedParams()等函数，构建与apiServer交互的参数。
		Do方法通过request发送请求，然后通过transformResponse解析请求返回，并绑定到对应资源对象的结构体上，
		这里的话，就表示corev1.podlist的对象。
		request先检查了有没有可用的Client。这里的开始调用net/http包的功能。
	*/

	for _, item := range result.Items {
		fmt.Printf("namesapce: %v,name: %v\n", item.Namespace, item.Name)
	}
}
