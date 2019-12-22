/*
Copyright 2018 BlackRock, Inc.

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
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/argoproj/argo-events/pkg/apis/eventsource/v1alpha1"
	scheme "github.com/argoproj/argo-events/pkg/client/eventsources/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// EventSourcesGetter has a method to return a EventSourceInterface.
// A group's client should implement this interface.
type EventSourcesGetter interface {
	EventSources(namespace string) EventSourceInterface
}

// EventSourceInterface has methods to work with EventSource resources.
type EventSourceInterface interface {
	Create(*v1alpha1.EventSource) (*v1alpha1.EventSource, error)
	Update(*v1alpha1.EventSource) (*v1alpha1.EventSource, error)
	UpdateStatus(*v1alpha1.EventSource) (*v1alpha1.EventSource, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.EventSource, error)
	List(opts v1.ListOptions) (*v1alpha1.EventSourceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EventSource, err error)
	EventSourceExpansion
}

// eventSources implements EventSourceInterface
type eventSources struct {
	client rest.Interface
	ns     string
}

// newEventSources returns a EventSources
func newEventSources(c *ArgoprojV1alpha1Client, namespace string) *eventSources {
	return &eventSources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the eventSource, and returns the corresponding eventSource object, and an error if there is any.
func (c *eventSources) Get(name string, options v1.GetOptions) (result *v1alpha1.EventSource, err error) {
	result = &v1alpha1.EventSource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("eventsources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of EventSources that match those selectors.
func (c *eventSources) List(opts v1.ListOptions) (result *v1alpha1.EventSourceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.EventSourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("eventsources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested eventSources.
func (c *eventSources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("eventsources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a eventSource and creates it.  Returns the server's representation of the eventSource, and an error, if there is any.
func (c *eventSources) Create(eventSource *v1alpha1.EventSource) (result *v1alpha1.EventSource, err error) {
	result = &v1alpha1.EventSource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("eventsources").
		Body(eventSource).
		Do().
		Into(result)
	return
}

// Update takes the representation of a eventSource and updates it. Returns the server's representation of the eventSource, and an error, if there is any.
func (c *eventSources) Update(eventSource *v1alpha1.EventSource) (result *v1alpha1.EventSource, err error) {
	result = &v1alpha1.EventSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("eventsources").
		Name(eventSource.Name).
		Body(eventSource).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *eventSources) UpdateStatus(eventSource *v1alpha1.EventSource) (result *v1alpha1.EventSource, err error) {
	result = &v1alpha1.EventSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("eventsources").
		Name(eventSource.Name).
		SubResource("status").
		Body(eventSource).
		Do().
		Into(result)
	return
}

// Delete takes name of the eventSource and deletes it. Returns an error if one occurs.
func (c *eventSources) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("eventsources").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *eventSources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("eventsources").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched eventSource.
func (c *eventSources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EventSource, err error) {
	result = &v1alpha1.EventSource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("eventsources").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
