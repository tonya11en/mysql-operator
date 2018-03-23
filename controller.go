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

// Creates a controller watching for mysql custom resources.
func newMySqlController(context *opkit.Context, mySqlClientset mysqlclient.MyprojectV1alpha1Interface) *MySqlController {
	return &MySqlController{
		context:        context,
		mySqlClientset: mySqlClientset,
	}
}

// Watch watches for instances of MySql custom resources and acts on them
func (c *MySqlController) StartWatch(namespace string, stopCh chan struct{}) error {
	fmt.Println("Starting watch on the mysql resource")

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

	fmt.Println("Handling MySql add")
}

func (c *MySqlController) onUpdate(oldObj, newObj interface{}) {
	sOld := oldObj.(*mysql.MySql).DeepCopy()
	sNew := newObj.(*mysql.MySql).DeepCopy()

	fmt.Println("Handling MySql update")
}

func (c *MySqlController) onDelete(obj interface{}) {
	s := obj.(*mysql.MySql).DeepCopy()

	fmt.Println("Handling MySql delete")
}
