/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	v1alpha1 "github.com/tonya11en/mysql-operator/pkg/apis/myproject/v1alpha1"
	scheme "github.com/tonya11en/mysql-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MySqlsGetter has a method to return a MySqlInterface.
// A group's client should implement this interface.
type MySqlsGetter interface {
	MySqls(namespace string) MySqlInterface
}

// MySqlInterface has methods to work with MySql resources.
type MySqlInterface interface {
	Create(*v1alpha1.MySql) (*v1alpha1.MySql, error)
	Update(*v1alpha1.MySql) (*v1alpha1.MySql, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.MySql, error)
	List(opts v1.ListOptions) (*v1alpha1.MySqlList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MySql, err error)
	MySqlExpansion
}

// mySqls implements MySqlInterface
type mySqls struct {
	client rest.Interface
	ns     string
}

// newMySqls returns a MySqls
func newMySqls(c *MyprojectV1alpha1Client, namespace string) *mySqls {
	return &mySqls{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mySql, and returns the corresponding mySql object, and an error if there is any.
func (c *mySqls) Get(name string, options v1.GetOptions) (result *v1alpha1.MySql, err error) {
	result = &v1alpha1.MySql{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mysqls").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MySqls that match those selectors.
func (c *mySqls) List(opts v1.ListOptions) (result *v1alpha1.MySqlList, err error) {
	result = &v1alpha1.MySqlList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mysqls").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mySqls.
func (c *mySqls) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mysqls").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a mySql and creates it.  Returns the server's representation of the mySql, and an error, if there is any.
func (c *mySqls) Create(mySql *v1alpha1.MySql) (result *v1alpha1.MySql, err error) {
	result = &v1alpha1.MySql{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mysqls").
		Body(mySql).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mySql and updates it. Returns the server's representation of the mySql, and an error, if there is any.
func (c *mySqls) Update(mySql *v1alpha1.MySql) (result *v1alpha1.MySql, err error) {
	result = &v1alpha1.MySql{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mysqls").
		Name(mySql.Name).
		Body(mySql).
		Do().
		Into(result)
	return
}

// Delete takes name of the mySql and deletes it. Returns an error if one occurs.
func (c *mySqls) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mysqls").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mySqls) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mysqls").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mySql.
func (c *mySqls) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MySql, err error) {
	result = &v1alpha1.MySql{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mysqls").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
