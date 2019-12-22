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
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/argoproj/argo-events/common"
	"github.com/argoproj/argo-events/gateways"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/google/uuid"
)

// dispatchEvent dispatches event to gateway transformer for further processing
func (gatewayContext *GatewayContext) dispatchEvent(gatewayEvent *gateways.Event) error {
	logger := gatewayContext.logger.WithField(common.LabelEventSource, gatewayEvent.Name)
	logger.Infoln("dispatching event to subscribers")

	cloudEvent, err := gatewayContext.transformEvent(gatewayEvent)
	if err != nil {
		return err
	}

	completeSuccess := true

	for _, subscription := range gatewayContext.httpSubscriptions {
		client, ok := gatewayContext.httpClients[subscription.Name]
		if !ok {
			logger.WithField("subscription-name", subscription.Name).Errorln("failed to send event. no http client available")
			continue
		}
		if _, _, err := client.Send(context.Background(), *cloudEvent); err != nil {
			logger.WithError(err).WithFields(logrus.Fields{
				"subscription-name": subscription.Name,
				"subscription-url":  subscription.URL,
			}).Warnln("failed to send the event")
			completeSuccess = false
		}
	}

	for _, subscription := range gatewayContext.natsSubscriptions {
		client, ok := gatewayContext.natsClients[subscription.Name]
		if !ok {
			logger.WithField("subscription-name", subscription.Name).Errorln("failed to send event. no nats client available")
			continue
		}
		if _, _, err := client.Send(context.Background(), *cloudEvent); err != nil {
			logger.WithError(err).WithFields(logrus.Fields{
				"subscription-name":       subscription.Name,
				"subscription-server-url": subscription.ServerURL,
				"subscription-subject":    subscription.Subject,
			}).Warnln("failed to send the event")
			completeSuccess = false
		}
	}

	response := "dispatched event to all subscribers"
	if !completeSuccess {
		response = fmt.Sprintf("%s.%s", response, " although some of the dispatch operations failed, check logs for more info")
	}

	logger.Infoln(response)
	return nil
}

// transformEvent transforms an event from gateway server into a CloudEvent
// See https://github.com/cloudevents/spec for more info.
func (gatewayContext *GatewayContext) transformEvent(gatewayEvent *gateways.Event) (*cloudevents.Event, error) {
	event := cloudevents.NewEvent(cloudevents.VersionV03)
	event.SetID(fmt.Sprintf("%x", uuid.New()))
	event.SetType(string(gatewayContext.gateway.Spec.Type))
	event.SetSource(gatewayContext.gateway.Name)
	event.SetDataContentType("application/json")
	event.SetSubject(gatewayEvent.Name)
	event.SetTime(time.Now())
	if err := event.SetData(gatewayEvent.Payload); err != nil {
		return nil, err
	}
	return &event, nil
}
