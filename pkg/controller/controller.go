package controller

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/jgensler8/openfaas-controller/pkg/apis/cr/v1"
)

// Watcher is an example of watching on resource create/update/delete events
type FunctionController struct {
	FunctionInterface rest.Interface
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
		c.FunctionInterface,
		"functions",
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		//&crv1.Example{},
		&v1.FunctionList{},

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
	fmt.Printf("[CONTROLLER] OnAdd %v\n", function)
}

func (c *FunctionController) onUpdate(oldObj, newObj interface{}) {
	oldFunction := oldObj.(*v1.Function)
	newFunction := newObj.(*v1.Function)
	fmt.Printf("[CONTROLLER] OnUpdate oldObj: %v\n", oldFunction)
	fmt.Printf("[CONTROLLER] OnUpdate newObj: %v\n", newFunction)
}

func (c *FunctionController) onDelete(obj interface{}) {
	function := obj.(*v1.Function)
	fmt.Printf("[CONTROLLER] OnDelete %v\n", function)
}