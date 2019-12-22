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

package sensor

import (
	"github.com/argoproj/argo-events/common"
	apicommon "github.com/argoproj/argo-events/pkg/apis/common"
	"github.com/argoproj/argo-events/pkg/apis/sensor/v1alpha1"
	subscriptionv1alpha1 "github.com/argoproj/argo-events/pkg/apis/subscription/v1alpha1"
	sensorinformers "github.com/argoproj/argo-events/pkg/client/sensor/informers/externalversions"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/cache"
	"strconv"
)

func (controller *Controller) instanceIDReq() (*labels.Requirement, error) {
	var instanceIDReq *labels.Requirement
	var err error
	if controller.Config.InstanceID == "" {
		return nil, errors.New("controller instance id must be specified")
	}
	instanceIDReq, err = labels.NewRequirement(LabelControllerInstanceID, selection.Equals, []string{controller.Config.InstanceID})
	if err != nil {
		panic(err)
	}
	return instanceIDReq, nil
}

func getHTTPSubscriptionIndex(subscriptions []subscriptionv1alpha1.HTTPSubscription, name string) int {
	for index, subscription := range subscriptions {
		if subscription.Name == name {
			return index
		}
	}
	return -1
}

func getNATSSubscriptionIndex(subscriptions []subscriptionv1alpha1.NATSSubscription, name string) int {
	for index, subscription := range subscriptions {
		if subscription.Name == name {
			return index
		}
	}
	return -1
}

// updateSubscription updates the subscription resource
func (controller *Controller) updateSubscription(obj interface{}, eventType EventType) {
	sensor, ok := obj.(*v1alpha1.Sensor)
	if !ok {
		controller.logger.Errorln("failed to update the subscription for sensor. unable to parse the sensor object")
		return
	}
	if err := ValidateSensor(sensor); err != nil {
		controller.logger.WithError(err).Errorln("failed to validate the sensor object, won't process the subscription updates")
		return
	}

	if sensor.Spec.SubscriptionRef != nil {
		controller.logger.WithField("sensor-name", sensor.Name).Warnln("sensor doesn't have a subscription reference")
		return
	}

	logger := controller.logger.WithFields(logrus.Fields{
		"sensor-name":       sensor.Name,
		"subscription-name": sensor.Spec.SubscriptionRef.Name,
	})

	namespace := sensor.Namespace
	if sensor.Spec.SubscriptionRef.Namespace != "" {
		namespace = sensor.Spec.SubscriptionRef.Namespace
	}

	subscription, err := controller.subscriptionClient.ArgoprojV1alpha1().Subscriptions(namespace).Get(sensor.Name, metav1.GetOptions{})
	if err != nil {
		logger.WithError(err).Errorln("failed to retrieve the subscription resource, unable to unsubscribe")
		return
	}

	updated := false

	switch eventType {
	case UpdateEvent:
		switch sensor.Spec.EventProtocol.Type {
		case apicommon.HTTP:
			if sensor.Status.Resources != nil && sensor.Status.Resources.Service != nil {
				port := strconv.Itoa(common.SensorServerPort)
				if sensor.Spec.EventProtocol.Http.Port != "" {
					port = sensor.Spec.EventProtocol.Http.Port
				}

				url := common.FormatServiceURL("http", common.ServiceDNSName(sensor.Status.Resources.Service.Name, namespace), port, "/")

				index := getHTTPSubscriptionIndex(subscription.Spec.HTTP, sensor.Name)
				if index != -1 && subscription.Spec.HTTP[index].URL != url {
					subscription.Spec.HTTP[index].URL = url
					updated = true
				}
				if index == -1 {
					subscription.Spec.HTTP = append(subscription.Spec.HTTP, subscriptionv1alpha1.HTTPSubscription{
						Name: sensor.Name,
						URL:  url,
					})
					updated = true
				}
			}
		case apicommon.NATS:
			index := getNATSSubscriptionIndex(subscription.Spec.NATS, sensor.Name)
			if index != -1 {
				if subscription.Spec.NATS[index].Subject != sensor.Spec.EventProtocol.Nats.Subject {
					subscription.Spec.NATS[index].Subject = sensor.Spec.EventProtocol.Nats.Subject
					updated = true
				}
				if subscription.Spec.NATS[index].ServerURL != sensor.Spec.EventProtocol.Nats.ServerURL {
					subscription.Spec.NATS[index].ServerURL = sensor.Spec.EventProtocol.Nats.ServerURL
					updated = true
				}
			}
			if index == -1 {
				subscription.Spec.NATS = append(subscription.Spec.NATS, subscriptionv1alpha1.NATSSubscription{
					Name:      sensor.Name,
					ServerURL: sensor.Spec.EventProtocol.Nats.ServerURL,
					Subject:   sensor.Spec.EventProtocol.Nats.Subject,
				})
			}
		}
	case DeleteEvent:
		index := getHTTPSubscriptionIndex(subscription.Spec.HTTP, sensor.Name)
		if index != -1 {
			subscription.Spec.HTTP[index] = subscription.Spec.HTTP[len(subscription.Spec.HTTP)-1]
			subscription.Spec.HTTP = subscription.Spec.HTTP[:len(subscription.Spec.HTTP)-1]
			updated = true
		}
		index = getNATSSubscriptionIndex(subscription.Spec.NATS, sensor.Name)
		if index != -1 {
			subscription.Spec.NATS[index] = subscription.Spec.NATS[len(subscription.Spec.NATS)-1]
			subscription.Spec.NATS = subscription.Spec.NATS[:len(subscription.Spec.NATS)-1]
			updated = true
		}
	default:
		logger.WithField("event-type", string(eventType)).Errorln("unknown sensor resource change event type")
		return
	}

	if updated {
		logger.Infoln("updating subscription resource")
		if _, err := controller.subscriptionClient.ArgoprojV1alpha1().Subscriptions(namespace).Update(subscription); err != nil {
			logger.WithError(err).Errorf("failed to update the subscription resource")
			return
		}
		logger.Infoln("subscription resource updated")
		return
	}

	logger.Infoln("no updates necessary for the subscription resource")
}

// The sensor informer adds new sensors to the controller'sensor queue based on Add, Update, and Delete event handlers for the sensor resources
func (controller *Controller) newSensorInformer() (cache.SharedIndexInformer, error) {
	labelSelector, err := controller.instanceIDReq()
	if err != nil {
		return nil, err
	}

	sensorInformerFactory := sensorinformers.NewSharedInformerFactoryWithOptions(
		controller.sensorClient,
		sensorResyncPeriod,
		sensorinformers.WithNamespace(controller.Config.Namespace),
		sensorinformers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.Everything().String()
			options.LabelSelector = labelSelector.String()
		}),
	)
	informer := sensorInformerFactory.Argoproj().V1alpha1().Sensors().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err == nil {
					controller.queue.Add(key)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(new)
				if err == nil {
					controller.queue.Add(key)
					// Update the subscription resource for the sensor object
					controller.updateSubscription(new, UpdateEvent)
				}
			},
			DeleteFunc: func(obj interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err == nil {
					controller.queue.Add(key)
					// Update the subscription resource for the sensor object
					controller.updateSubscription(obj, DeleteEvent)
				}
			},
		},
	)
	return informer, nil
}
