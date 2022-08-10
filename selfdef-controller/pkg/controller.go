package pkg

import (
	"context"
	"log"
	"reflect"
	"time"

	v14 "k8s.io/api/core/v1"
	v12 "k8s.io/api/networking/v1"
	errors_k8s "k8s.io/apimachinery/pkg/api/errors"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	informer "k8s.io/client-go/informers/core/v1"
	networking "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	corelister "k8s.io/client-go/listers/core/v1"
	v1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	workNum  = 5
	maxRetry = 10
)

type controller struct {
	client        kubernetes.Interface
	ingressLister v1.IngressLister
	serviceLister corelister.ServiceLister
	queue         workqueue.RateLimitingInterface
}

func (c *controller) addService(obj interface{}) {
	c.enqueue(obj)
}
func (c *controller) updateService(oldObj interface{}, newObj interface{}) {
	//compare annatation
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	c.enqueue(newObj)
}
func (c *controller) deleteIngress(obj interface{}) {
	ingress := obj.(*v12.Ingress)
	ownerReference := v13.GetControllerOf(ingress)
	if ownerReference == nil {
		return
	}
	if ownerReference.Kind != "" {
		return
	}
	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceIndexFunc(obj)
	if err != nil {
		return
	}
	c.queue.Add(key)
}
func (c *controller) syncService(key string) error {
	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	//delete
	service, err2 := c.serviceLister.Services(namespaceKey).Get(name)

	if errors_k8s.IsNotFound(err2) {
		return nil
	}

	if err2 != nil {
		return err2
	}

	//add and delete
	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err3 := c.ingressLister.Ingresses(namespaceKey).Get(name)

	if err3 != nil && !errors_k8s.IsNotFound(err3) {
		return err3
	}
	if ok && errors_k8s.IsNotFound(err3) {
		//create ingress
		ig := c.constructIngress(service)
		_, err4 := c.client.NetworkingV1().Ingresses(namespaceKey).Create(context.TODO(), ig, v13.CreateOptions{})
		if err4 != nil {
			return err4
		}

	} else if !ok && ingress != nil {
		//delete ingress
		err5 := c.client.NetworkingV1().IngressClasses().Delete(context.TODO(), name, v13.DeleteOptions{})
		if err5 != nil {
			return err5
		}
	}
	return nil
}
func (c *controller) constructIngress(service *v14.Service) *v12.Ingress {
	ingress := v12.Ingress{}
	ingress.Name = service.Name
	ingress.Namespace = service.Namespace
	pathType := v12.PathTypePrefix
	ingress.Spec = v12.IngressSpec{
		Rules: []v12.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: v12.IngressRuleValue{
					HTTP: &v12.HTTPIngressRuleValue{
						Paths: []v12.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: v12.IngressBackend{
									Service: &v12.IngressServiceBackend{
										Name: service.Name,
										Port: v12.ServiceBackendPort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return &ingress
}
func (c *controller) HandleError(key string, err error) {
	if c.queue.NumRequeues(key) <= maxRetry {
		c.queue.AddRateLimited(key)
		return
	}
	if err != nil {
		log.Fatalln(err)
	}
	c.queue.Forget(key)
}
func (c *controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(item)

	key := item.(string)
	err := c.syncService(key)
	if err != nil {
		c.HandleError(key, err)
	}
	return true
}

//loop workqueue get key,then deal
func (c *controller) worker() {
	for c.processNextItem() {

	}
}
func (c *controller) Run(stopCh chan struct{}) {
	for i := 0; i < workNum; i++ {
		go wait.Until(c.worker, time.Minute, stopCh)
	}
	<-stopCh
}
func NewController(client kubernetes.Interface, serviceInformer informer.ServiceInformer, ingressInfromer networking.IngressInformer) controller {
	c := controller{
		client:        client,
		ingressLister: ingressInfromer.Lister(),
		serviceLister: serviceInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingressInfromer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})

	return c
}
