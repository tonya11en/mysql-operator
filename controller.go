package main

import (
	"fmt"

	opkit "github.com/rook/operator-kit"
	mysql "github.com/tonya11en/mysql-operator/pkg/apis/myproject/v1alpha1"
	mysqlclient "github.com/tonya11en/mysql-operator/pkg/client/clientset/versioned/typed/myproject/v1alpha1"
	"k8s.io/client-go/tools/cache"
)

// MySqlController represents a controller object for mysql custom resources
type MySqlController struct {
	context        *opkit.Context
	mySqlClientset mysqlclient.MyprojectV1alpha1Interface
}

// newMySqlController creates a controller for watching mysql custom resources created
func newMySqlController(context *opkit.Context, mySqlClientset mysqlclient.MyprojectV1alpha1Interface) *MySqlController {
	return &MySqlController{
		context:        context,
		mySqlClientset: mySqlClientset,
	}
}

// Watch watches for instances of MySql custom resources and acts on them
func (c *MySqlController) StartWatch(namespace string, stopCh chan struct{}) error {

	resourceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	}
	restClient := c.mySqlClientset.RESTClient()
	watcher := opkit.NewWatcher(mysql.MySqlResource, namespace, resourceHandlers, restClient)
	go watcher.Watch(&mysql.MySql{}, stopCh)
	return nil
}

func (c *MySqlController) onAdd(obj interface{}) {
	s := obj.(*mysql.MySql).DeepCopy()

	fmt.Printf("Added MySql")
}

func (c *MySqlController) onUpdate(oldObj, newObj interface{}) {
	fmt.Printf("MySql update handler")
}

func (c *MySqlController) onDelete(obj interface{}) {
	fmt.Printf("MySql delete handler")
}
