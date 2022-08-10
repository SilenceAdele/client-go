package main

import (
	"fmt"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//create config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	//create client
	clientset, err2 := kubernetes.NewForConfig(config)
	if err2 != nil {
		panic(err2)
	}
	//get infromer

	//1)get all namespace pod
	//factory := informers.NewSharedInformerFactory(clientset, 0)
	//informer := factory.Core().V1().Pods().Informer()

	//2)get special namespace pod
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace("zzy"))
	informer := factory.Core().V1().Pods().Informer()

	//add queue
	//rateLimiterQueue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "controller")

	//add event handler
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("...AddFunc")
			// key, err3 := cache.MetaNamespaceIndexFunc(obj)
			// if err3 != nil {
			// 	fmt.Println("can't get key")
			// }
			// rateLimiterQueue.AddRateLimited(key)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("...DeleteFunc")
			// key, err3 := cache.MetaNamespaceIndexFunc(obj)
			// if err3 != nil {
			// 	fmt.Println("can't get key")
			// }
			// rateLimiterQueue.AddRateLimited(key)
		},
		UpdateFunc: func(oldobj, obj interface{}) {
			fmt.Println("...UpdateFunc")
			// key, err3 := cache.MetaNamespaceIndexFunc(obj)
			// if err3 != nil {
			// 	fmt.Println("can't get key")
			// }
			// rateLimiterQueue.AddRateLimited(key)
		},
	})
	//start informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
	<-stopCh
}
