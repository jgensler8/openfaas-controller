package controller

import (
	"context"
	"fmt"
	"log"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/jgensler8/openfaas-controller/pkg/apis/cr/v1"
	"github.com/jgensler8/openfaas-client-go"
	"encoding/json"
)

// Watcher is an example of watching on resource create/update/delete events
type FunctionController struct {
	KubernetesFunctionInterface rest.Interface
	OpenFaaSFunctionAPIClient *swagger.DefaultApi
}

// Run starts an Example resource controller
func (c *FunctionController) Run(ctx context.Context) error {
	fmt.Print("Watch Example objects\n")

	// Watch Example objects
	_, err := c.watchFunctions(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for Example resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *FunctionController) watchFunctions(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.KubernetesFunctionInterface,
		"functions",
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&v1.Function{},

		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		})

	go controller.Run(ctx.Done())
	return controller, nil
}

func (c *FunctionController) onAdd(obj interface{}) {
	function := obj.(*v1.Function)
	log.Printf("[CONTROLLER] OnAdd %v\n", function)

	bytes, err := json.Marshal(function.Spec)
	if err != nil {
		log.Printf("Failed to Marshal Function.Spec for function (%s) in namespace (%s)", function.Name, function.Namespace)
		log.Printf("%v", err)
		return
	}

	req := swagger.CreateFunctionRequest{}
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		log.Printf("Failed to Unmarshal Function.Spec for function (%s) in namespace (%s)", function.Name, function.Namespace)
		log.Printf("%v", err)
		return
	}
	res, err := c.OpenFaaSFunctionAPIClient.SystemFunctionsPost(req)
	if err != nil || res.StatusCode != 202 {
		log.Printf("API call (CREATE) to OpenFaaS server failed for function (%s) in namespace (%s)", function.Name, function.Namespace)
		log.Printf("%v", err)
		return
	}
	log.Printf("Success Create: %v", res)
}

func (c *FunctionController) onUpdate(oldObj, newObj interface{}) {
	oldFunction := oldObj.(*v1.Function)
	newFunction := newObj.(*v1.Function)
	log.Printf("[CONTROLLER] OnUpdate oldObj: %v\n", oldFunction)
	log.Printf("[CONTROLLER] OnUpdate newObj: %v\n", newFunction)

	if oldFunction.Spec.Service != newFunction.Spec.Service {
		c.onDelete(oldObj)
		c.onAdd(newObj)
	} else {
		bytes, err := json.Marshal(newFunction.Spec)
		if err != nil {
			log.Printf("Failed to Marshal Function.Spec for function (%s) in namespace (%s)", newFunction.Name, newFunction.Namespace)
			log.Printf("%v", err)
			return
		}

		req := swagger.CreateFunctionRequest{}
		err = json.Unmarshal(bytes, &req)
		if err != nil {
			log.Printf("Failed to Unmarshal Function.Spec for function (%s) in namespace (%s)", newFunction.Name, newFunction.Namespace)
			log.Printf("%v", err)
			return
		}
		res, err := c.OpenFaaSFunctionAPIClient.SystemFunctionsPut(req)
		if err != nil || res.StatusCode != 200 {
			log.Printf("API call (PUT) to OpenFaaS server failed for function (%s) in namespace (%s)", newFunction.Name, newFunction.Namespace)
			log.Printf("%v", err)
			log.Printf("Payload: %s", res.Payload)
			return
		}
	}
	log.Printf("Success Update: %s", newFunction.Spec.Service)
}

func (c *FunctionController) onDelete(obj interface{}) {
	function := obj.(*v1.Function)
	log.Printf("[CONTROLLER] OnDelete %v\n", function)

	req := swagger.DeleteFunctionRequest{
		FunctionName: function.Spec.Service,
	}
	res, err := c.OpenFaaSFunctionAPIClient.SystemFunctionsDelete(req)
	if err != nil || res.StatusCode != 200 {
		log.Printf("API call (DELETE) to OpenFaaS server failed for function (%s) in namespace (%s)", function.Name, function.Namespace)
		log.Printf("%v", err)
		return
	}
	log.Printf("Success Delete: %v", res)
}