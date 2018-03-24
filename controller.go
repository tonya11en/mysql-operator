/*
Copyright 2018 Tony Allen. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Some of the code below came from:
https://github.com/rook/rook
https://github.com/rook/operator-kit

which have the same license.
*/

package main

import (
	"fmt"

	opkit "github.com/rook/operator-kit"
	mysql "github.com/tonya11en/mysql-operator/pkg/apis/myproject/v1alpha1"
	mysqlclient "github.com/tonya11en/mysql-operator/pkg/client/clientset/versioned/typed/myproject/v1alpha1"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Create a pod.
func (c *MySqlController) makePod(objName string, ctrName string, ctrImage string, image string, port int32, podGroup string) (*v1.Pod, error) {
	coreV1Client := c.context.Clientset.CoreV1()
	pod, err := coreV1Client.Pods(v1.NamespaceDefault).Create(&v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: objName,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  ctrName,
					Image: ctrImage,
					Ports: []v1.ContainerPort{
						{
							ContainerPort: port,
						},
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Errorf("failed to create pod. %+v", err)
		return pod, err
	}

	pod.SetLabels(map[string]string{"pod-group": podGroup})
	pod, err = coreV1Client.Pods(v1.NamespaceDefault).Update(pod)
	if err != nil {
		fmt.Errorf("failed to set label on pod. %+v", err)
	}

	return pod, err
}

func (c *MySqlController) onAdd(obj interface{}) {
	fmt.Println("Handling MySql add")

	// Make mysql service
	// Make persistent volume claim
	// Make deployment
}

func (c *MySqlController) onUpdate(oldObj, newObj interface{}) {
	fmt.Println("Handling MySql update")
}

func (c *MySqlController) onDelete(obj interface{}) {
	fmt.Println("Handling MySql delete")
}
