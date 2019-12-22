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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/argoproj/argo-events/pkg/apis/eventsource/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// EventSourceLister helps list EventSources.
type EventSourceLister interface {
	// List lists all EventSources in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.EventSource, err error)
	// EventSources returns an object that can list and get EventSources.
	EventSources(namespace string) EventSourceNamespaceLister
	EventSourceListerExpansion
}

// eventSourceLister implements the EventSourceLister interface.
type eventSourceLister struct {
	indexer cache.Indexer
}

// NewEventSourceLister returns a new EventSourceLister.
func NewEventSourceLister(indexer cache.Indexer) EventSourceLister {
	return &eventSourceLister{indexer: indexer}
}

// List lists all EventSources in the indexer.
func (s *eventSourceLister) List(selector labels.Selector) (ret []*v1alpha1.EventSource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EventSource))
	})
	return ret, err
}

// EventSources returns an object that can list and get EventSources.
func (s *eventSourceLister) EventSources(namespace string) EventSourceNamespaceLister {
	return eventSourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EventSourceNamespaceLister helps list and get EventSources.
type EventSourceNamespaceLister interface {
	// List lists all EventSources in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.EventSource, err error)
	// Get retrieves the EventSource from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.EventSource, error)
	EventSourceNamespaceListerExpansion
}

// eventSourceNamespaceLister implements the EventSourceNamespaceLister
// interface.
type eventSourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all EventSources in the indexer for a given namespace.
func (s eventSourceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.EventSource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EventSource))
	})
	return ret, err
}

// Get retrieves the EventSource from the indexer for a given namespace and name.
func (s eventSourceNamespaceLister) Get(name string) (*v1alpha1.EventSource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("eventsource"), name)
	}
	return obj.(*v1alpha1.EventSource), nil
}
