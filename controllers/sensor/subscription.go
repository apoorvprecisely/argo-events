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
	"strconv"

	"github.com/argoproj/argo-events/common"
	"github.com/argoproj/argo-events/pkg/apis/sensor/v1alpha1"
	subscriptionv1alpha1 "github.com/argoproj/argo-events/pkg/apis/subscription/v1alpha1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func (controller *Controller) getSubscriptionResource(namespace string, ref *v1alpha1.Subscription) (*subscriptionv1alpha1.Subscription, error) {
	if ref.Namespace != "" {
		namespace = ref.Namespace
	}
	return controller.subscriptionClient.ArgoprojV1alpha1().Subscriptions(namespace).Get(ref.Name, metav1.GetOptions{})
}

func (controller *Controller) getUpdatedHTTPSubscription(sensor *v1alpha1.Sensor) (*subscriptionv1alpha1.Subscription, error) {
	if sensor.Spec.EventProtocol.HTTP != nil && sensor.Status.Resources != nil && sensor.Status.Resources.Service != nil {
		updated := false
		protocol := sensor.Spec.EventProtocol.HTTP
		port := strconv.Itoa(common.SensorServerPort)
		if protocol.Port != "" {
			port = protocol.Port
		}
		endpoint := common.SensorServiceEndpoint
		if protocol.Endpoint != "" {
			endpoint = protocol.Endpoint
		}

		url := common.FormatServiceURL("http", common.ServiceDNSName(sensor.Status.Resources.Service.Name, sensor.Status.Resources.Service.Namespace), port, endpoint)

		subscription, err := controller.getSubscriptionResource(sensor.Namespace, protocol.SubscriptionRef)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve the subscription %s resource, unable to continue sensor subscription update", protocol.SubscriptionRef.Name)
		}

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

		if updated {
			return subscription, nil
		}
	}
	return nil, nil
}

func (controller *Controller) getUpdatedNATSSubscriptions(sensor *v1alpha1.Sensor) (*subscriptionv1alpha1.Subscription, error) {
	if sensor.Spec.EventProtocol.NATS != nil {
		updated := false
		protocol := sensor.Spec.EventProtocol.NATS
		subscription, err := controller.getSubscriptionResource(sensor.Namespace, protocol.SubscriptionRef)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve the subscription %s resource, unable to continue sensor subscription update", protocol.SubscriptionRef.Name)
		}

		index := getNATSSubscriptionIndex(subscription.Spec.NATS, sensor.Name)
		if index != -1 {
			if subscription.Spec.NATS[index].Subject != sensor.Spec.EventProtocol.NATS.Subject {
				subscription.Spec.NATS[index].Subject = sensor.Spec.EventProtocol.NATS.Subject
				updated = true
			}
			if subscription.Spec.NATS[index].ServerURL != sensor.Spec.EventProtocol.NATS.ServerURL {
				subscription.Spec.NATS[index].ServerURL = sensor.Spec.EventProtocol.NATS.ServerURL
				updated = true
			}
		}
		if index == -1 {
			subscription.Spec.NATS = append(subscription.Spec.NATS, subscriptionv1alpha1.NATSSubscription{
				Name:      sensor.Name,
				ServerURL: sensor.Spec.EventProtocol.NATS.ServerURL,
				Subject:   sensor.Spec.EventProtocol.NATS.Subject,
			})
			updated = true
		}
		if updated {
			return subscription, nil
		}
	}
	return nil, nil
}

func (controller *Controller) deleteHTTPSubscription(sensor *v1alpha1.Sensor) (*subscriptionv1alpha1.Subscription, error) {
	protocol := sensor.Spec.EventProtocol.HTTP
	subscription, err := controller.getSubscriptionResource(sensor.Namespace, protocol.SubscriptionRef)
	if err != nil {
		return nil, err
	}
	if index := getHTTPSubscriptionIndex(subscription.Spec.HTTP, sensor.Name); index != -1 {
		subscription.Spec.HTTP[index] = subscription.Spec.HTTP[len(subscription.Spec.HTTP)-1]
		subscription.Spec.HTTP = subscription.Spec.HTTP[:len(subscription.Spec.HTTP)-1]
		return subscription, nil
	}
	return nil, nil
}

func (controller *Controller) deleteNATSSubscription(sensor *v1alpha1.Sensor) (*subscriptionv1alpha1.Subscription, error) {
	protocol := sensor.Spec.EventProtocol.NATS
	subscription, err := controller.getSubscriptionResource(sensor.Namespace, protocol.SubscriptionRef)
	if err != nil {
		return nil, err
	}
	if index := getNATSSubscriptionIndex(subscription.Spec.NATS, sensor.Name); index != -1 {
		subscription.Spec.NATS[index] = subscription.Spec.NATS[len(subscription.Spec.NATS)-1]
		subscription.Spec.NATS = subscription.Spec.NATS[:len(subscription.Spec.NATS)-1]
		return subscription, nil
	}
	return nil, nil
}

func (controller *Controller) updateSubscriptionResource(sensorName string, subscription *subscriptionv1alpha1.Subscription) {
	if subscription != nil {
		if _, err := controller.subscriptionClient.ArgoprojV1alpha1().Subscriptions(subscription.Namespace).Update(subscription); err != nil {
			controller.logger.WithError(err).WithField("sensor-name", sensorName).Errorln("failed to update the subscription resource")
			return
		}
		controller.logger.WithField("sensor-name", sensorName).Infoln("successfully processed http update for the subscription resource")
	}
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

	switch eventType {
	case UpdateEvent:
		controller.logger.WithField("sensor-name", sensor.Name).Infoln("processing subscription update...")
		subscription, err := controller.getUpdatedHTTPSubscription(sensor)
		if err != nil {
			controller.logger.WithError(err).WithField("sensor-name", sensor.Name).Errorln("failed to update http subscriptions")
			return
		}
		controller.updateSubscriptionResource(sensor.Name, subscription)

		subscription, err = controller.getUpdatedNATSSubscriptions(sensor)
		if err != nil {
			controller.logger.WithError(err).WithField("sensor-name", sensor.Name).Errorln("failed to update nats subscriptions")
			return
		}
		controller.updateSubscriptionResource(sensor.Name, subscription)

		controller.logger.WithField("sensor-name", sensor.Name).Infoln("completed subscription update")

	case DeleteEvent:
		controller.logger.WithField("sensor-name", sensor.Name).Infoln("processing subscription deletion...")
		subscription, err := controller.deleteHTTPSubscription(sensor)
		if err != nil {
			controller.logger.WithError(err).WithField("sensor-name", sensor.Name).Errorln("failed to update http subscriptions")
			return
		}
		controller.updateSubscriptionResource(sensor.Name, subscription)

		subscription, err = controller.deleteNATSSubscription(sensor)
		if err != nil {
			controller.logger.WithError(err).WithField("sensor-name", sensor.Name).Errorln("failed to update nats subscriptions")
			return
		}
		controller.updateSubscriptionResource(sensor.Name, subscription)

		controller.logger.WithField("sensor-name", sensor.Name).Infoln("completed subscription deletion")

	default:
		controller.logger.WithField("event-type", string(eventType)).Errorln("unknown sensor resource change event type")
		return
	}
	controller.logger.Infoln("no updates necessary for the subscription resource")
}
