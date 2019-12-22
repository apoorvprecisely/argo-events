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

package main

import (
	subscriptionv1alpha1 "github.com/argoproj/argo-events/pkg/apis/subscription/v1alpha1"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventsnats "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
	"github.com/sirupsen/logrus"
)

func (gatewayContext *GatewayContext) updateSubscriptions(subscription subscriptionv1alpha1.Subscription) {
	gatewayContext.httpSubscriptions = subscription.Spec.HTTP
	gatewayContext.natsSubscriptions = subscription.Spec.NATS
	if gatewayContext.httpClients == nil {
		gatewayContext.httpClients = map[string]client.Client{}
	}
	if gatewayContext.natsClients == nil {
		gatewayContext.natsClients = map[string]client.Client{}
	}
	for _, subscription := range gatewayContext.httpSubscriptions {
		if _, ok := gatewayContext.httpClients[subscription.Name]; !ok {
			t, err := cloudevents.NewHTTPTransport(
				cloudevents.WithTarget(subscription.URL),
				cloudevents.WithEncoding(cloudevents.HTTPBinaryV02),
			)
			if err != nil {
				gatewayContext.logger.WithError(err).WithFields(logrus.Fields{
					"subscription-name": subscription.Name,
					"subscription-url":  subscription.URL,
				}).Warnln("failed to create a http transport")
				continue
			}

			c, err := cloudevents.NewClient(t)
			if err != nil {
				gatewayContext.logger.WithError(err).WithFields(logrus.Fields{
					"subscription-name": subscription.Name,
					"subscription-url":  subscription.URL,
				}).Warnln("failed to create a http client")
				continue
			}

			gatewayContext.logger.WithFields(logrus.Fields{
				"subscription-name": subscription.Name,
				"subscription-url":  subscription.URL,
			}).Infoln("created a http client")

			gatewayContext.httpClients[subscription.Name] = c
		}
	}
	for _, subscription := range gatewayContext.natsSubscriptions {
		if _, ok := gatewayContext.natsClients[subscription.Name]; !ok {
			t, err := cloudeventsnats.New(subscription.ServerURL, subscription.Subject)
			if err != nil {
				gatewayContext.logger.WithError(err).WithFields(logrus.Fields{
					"subscription-name":       subscription.Name,
					"subscription-server-url": subscription.ServerURL,
					"subscription-subject":    subscription.Subject,
				}).Warnln("failed to create a nats transport")
				continue
			}

			c, err := client.New(t)
			if err != nil {
				gatewayContext.logger.WithError(err).WithFields(logrus.Fields{
					"subscription-name":       subscription.Name,
					"subscription-server-url": subscription.ServerURL,
					"subscription-subject":    subscription.Subject,
				}).Warnln("failed to create a nats client")
				continue
			}

			gatewayContext.logger.WithFields(logrus.Fields{
				"subscription-name":       subscription.Name,
				"subscription-server-url": subscription.ServerURL,
				"subscription-subject":    subscription.Subject,
			}).Infoln("created a nats client")

			gatewayContext.natsClients[subscription.Name] = c
		}
	}
}
