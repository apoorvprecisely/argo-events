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

package fake

import (
	v1alpha1 "github.com/argoproj/argo-events/pkg/apis/eventsource/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEventSources implements EventSourceInterface
type FakeEventSources struct {
	Fake *FakeArgoprojV1alpha1
	ns   string
}

var eventsourcesResource = schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "eventsources"}

var eventsourcesKind = schema.GroupVersionKind{Group: "argoproj.io", Version: "v1alpha1", Kind: "EventSource"}

// Get takes name of the eventSource, and returns the corresponding eventSource object, and an error if there is any.
func (c *FakeEventSources) Get(name string, options v1.GetOptions) (result *v1alpha1.EventSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(eventsourcesResource, c.ns, name), &v1alpha1.EventSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventSource), err
}

// List takes label and field selectors, and returns the list of EventSources that match those selectors.
func (c *FakeEventSources) List(opts v1.ListOptions) (result *v1alpha1.EventSourceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(eventsourcesResource, eventsourcesKind, c.ns, opts), &v1alpha1.EventSourceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.EventSourceList{ListMeta: obj.(*v1alpha1.EventSourceList).ListMeta}
	for _, item := range obj.(*v1alpha1.EventSourceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested eventSources.
func (c *FakeEventSources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(eventsourcesResource, c.ns, opts))

}

// Create takes the representation of a eventSource and creates it.  Returns the server's representation of the eventSource, and an error, if there is any.
func (c *FakeEventSources) Create(eventSource *v1alpha1.EventSource) (result *v1alpha1.EventSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(eventsourcesResource, c.ns, eventSource), &v1alpha1.EventSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventSource), err
}

// Update takes the representation of a eventSource and updates it. Returns the server's representation of the eventSource, and an error, if there is any.
func (c *FakeEventSources) Update(eventSource *v1alpha1.EventSource) (result *v1alpha1.EventSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(eventsourcesResource, c.ns, eventSource), &v1alpha1.EventSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventSource), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeEventSources) UpdateStatus(eventSource *v1alpha1.EventSource) (*v1alpha1.EventSource, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(eventsourcesResource, "status", c.ns, eventSource), &v1alpha1.EventSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventSource), err
}

// Delete takes name of the eventSource and deletes it. Returns an error if one occurs.
func (c *FakeEventSources) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(eventsourcesResource, c.ns, name), &v1alpha1.EventSource{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEventSources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(eventsourcesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.EventSourceList{})
	return err
}

// Patch applies the patch and returns the patched eventSource.
func (c *FakeEventSources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EventSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(eventsourcesResource, c.ns, name, pt, data, subresources...), &v1alpha1.EventSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventSource), err
}
