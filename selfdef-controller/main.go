package main

import (
	"client-go/selfdef-controller/pkg"
	"log"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. config

	//get config from cluster out
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		//if get config from out cluster failed,then from cluster in
		inClusterConfig, err2 := rest.InClusterConfig()
		if err2 != nil {
			log.Fatalln("cant't get config")
		}
		config = inClusterConfig
	}

	// 2. client
	clientset, err2 := kubernetes.NewForConfig(config)
	if err2 != nil {
		log.Fatalln("can't create client")
	}
	// 3. informer
	factory := informers.NewSharedInformerFactory(clientset, 0)
	serverInformer := factory.Core().V1().Services()
	ingressInfromer := factory.Networking().V1().Ingresses()

	controller := pkg.NewController(clientset, serverInformer, ingressInfromer)

	stopCh := make(chan struct{})
	factory.Start(stopCh)

	factory.WaitForCacheSync(stopCh)

	controller.Run(stopCh)
	// 4. add event handler
	// 5. informer.Start
}
