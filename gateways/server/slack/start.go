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

package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/argoproj/argo-events/common"
	"github.com/argoproj/argo-events/gateways"
	"github.com/argoproj/argo-events/gateways/server"
	"github.com/argoproj/argo-events/gateways/server/common/webhook"
	"github.com/argoproj/argo-events/pkg/apis/eventsource/v1alpha1"
	"github.com/argoproj/argo-events/store"
	"github.com/ghodss/yaml"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/pkg/errors"
)

// controller controls the webhook operations
var (
	controller = webhook.NewController()
)

// set up the activation and inactivation channels to control the state of routes.
func init() {
	go webhook.ProcessRouteStatus(controller)
}

// Implement Router
// 1. GetRoute
// 2. HandleRoute
// 3. PostActivate
// 4. PostDeactivate

// GetRoute returns the route
func (rc *Router) GetRoute() *webhook.Route {
	return rc.route
}

// HandleRoute handles incoming requests on the route
func (rc *Router) HandleRoute(writer http.ResponseWriter, request *http.Request) {
	route := rc.route

	logger := route.Logger.WithFields(
		map[string]interface{}{
			common.LabelEventSource: route.EventSource.Name,
			common.LabelEndpoint:    route.Context.Endpoint,
			common.LabelHTTPMethod:  route.Context.Method,
		})

	logger.Info("request a received, processing it...")

	if !route.Active {
		logger.Warn("endpoint is not active, won't process it")
		common.SendErrorResponse(writer, "endpoint is inactive")
		return
	}

	logger.Infoln("verifying the request...")
	err := rc.verifyRequest(request)
	if err != nil {
		logger.WithError(err).Error("failed to validate the request")
		common.SendInternalErrorResponse(writer, err.Error())
		return
	}

	var data []byte
	// Interactive element actions are always
	// sent as application/x-www-form-urlencoded
	// If request was generated by an interactive element, it will be a POST form
	if len(request.Header["Content-Type"]) > 0 && request.Header["Content-Type"][0] == "application/x-www-form-urlencoded" {
		logger.Infoln("handling slack interaction...")
		data, err = rc.handleInteraction(request)
		if err != nil {
			logger.WithError(err).Error("failed to process the interaction")
			common.SendInternalErrorResponse(writer, err.Error())
			return
		}
	} else {
		// If there's no payload in the post body, this is likely an
		// Event API request. Parse and process if valid.
		logger.Infoln("handling slack event...")
		var response []byte
		data, response, err = rc.handleEvent(request)
		if err != nil {
			logger.WithError(err).Error("failed  to handle the event")
			common.SendInternalErrorResponse(writer, err.Error())
			return
		}
		if response != nil {
			writer.Header().Set("Content-Type", "text")
			if _, err := writer.Write(response); err != nil {
				logger.WithError(err).Error("failed to write the response for url verification")
				// don't return, we want to keep this running to give user chance to retry
			}
		}
	}

	if data != nil {
		logger.Infoln("dispatching event on route's data channel...")
		route.DataCh <- data
	}

	logger.Info("request successfully processed")
	common.SendSuccessResponse(writer, "success")
}

// PostActivate performs operations once the route is activated and ready to consume requests
func (rc *Router) PostActivate() error {
	return nil
}

// PostInactivate performs operations after the route is inactivated
func (rc *Router) PostInactivate() error {
	return nil
}

// handleEvent parse the slack notification and validates the event type
func (rc *Router) handleEvent(request *http.Request) ([]byte, []byte, error) {
	var err error
	var response []byte
	var data []byte
	body, err := getRequestBody(request)
	if err != nil {
		return data, response, errors.Wrap(err, "failed to fetch request body")
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: rc.token}))
	if err != nil {
		return data, response, errors.Wrap(err, "failed to extract event")
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err = json.Unmarshal([]byte(body), &r)
		if err != nil {
			return data, response, errors.Wrap(err, "failed to verify the challenge")
		}
		response = []byte(r.Challenge)
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		data, err = json.Marshal(eventsAPIEvent.InnerEvent.Data)
		if err != nil {
			return data, response, errors.Wrap(err, "failed to marshal event data")
		}
	}

	return data, response, nil
}

func (rc *Router) handleInteraction(request *http.Request) ([]byte, error) {
	var err error
	err = request.ParseForm()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse post body")
	}

	payload := request.PostForm.Get("payload")
	ie := &slack.InteractionCallback{}
	err = json.Unmarshal([]byte(payload), ie)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse interaction event")
	}

	data, err := json.Marshal(ie)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal action data")
	}

	return data, nil
}

func getRequestBody(request *http.Request) ([]byte, error) {
	// Read request payload
	body, err := ioutil.ReadAll(request.Body)
	// Reset request.Body ReadCloser to prevent side-effect if re-read
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request body")
	}
	return body, nil
}

// If a signing secret is provided, validate the request against the
// X-Slack-Signature header value.
// The signature is a hash generated as per Slack documentation at:
// https://api.slack.com/docs/verifying-requests-from-slack
func (rc *Router) verifyRequest(request *http.Request) error {
	signingSecret := rc.signingSecret
	if len(signingSecret) > 0 {
		sv, err := slack.NewSecretsVerifier(request.Header, signingSecret)
		if err != nil {
			return errors.Wrap(err, "cannot create secrets verifier")
		}

		// Read the request body
		body, err := getRequestBody(request)
		if err != nil {
			return err
		}

		_, err = sv.Write([]byte(string(body)))
		if err != nil {
			return errors.Wrap(err, "error writing body: cannot verify signature")
		}

		err = sv.Ensure()
		if err != nil {
			return errors.Wrap(err, "signature validation failed")
		}
	}
	return nil
}

// StartEventSource starts a event source
func (listener *EventListener) StartEventSource(eventSource *gateways.EventSource, eventStream gateways.Eventing_StartEventSourceServer) error {
	defer server.Recover(eventSource.Name)

	logger := listener.Logger.WithField(common.LabelEventSource, eventSource.Name)

	logger.Infoln("started processing the event source...")

	logger.Infoln("parsing slack event source...")

	var slackEventSource *v1alpha1.SlackEventSource
	if err := yaml.Unmarshal(eventSource.Value, &slackEventSource); err != nil {
		logger.WithError(err).Errorln("failed to parse the event source")
		return err
	}

	logger.Infoln("retrieving the slack token...")
	token, err := store.GetSecrets(listener.K8sClient, slackEventSource.Namespace, slackEventSource.Token.Name, slackEventSource.Token.Key)
	if err != nil {
		logger.WithError(err).Error("failed to retrieve the token")
		return err
	}

	logger.Infoln("retrieving the signing secret...")
	signingSecret, err := store.GetSecrets(listener.K8sClient, slackEventSource.Namespace, slackEventSource.SigningSecret.Name, slackEventSource.SigningSecret.Key)
	if err != nil {
		logger.WithError(err).Warn("failed to retrieve the signing secret")
		return err
	}

	route := webhook.NewRoute(slackEventSource.Webhook, listener.Logger, eventSource)

	return webhook.ManageRoute(&Router{
		route:            route,
		token:            token,
		signingSecret:    signingSecret,
		k8sClient:        listener.K8sClient,
		slackEventSource: slackEventSource,
	}, controller, eventStream)
}
