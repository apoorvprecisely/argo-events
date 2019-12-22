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

package common

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EventProtocolType is type of the event dispatch protocol. Used for dispatching events
type EventProtocolType string

// possible types of event dispatch protocol
const (
	HTTP EventProtocolType = "HTTP"
	NATS EventProtocolType = "NATS"
)

// Type of nats connection.
type NatsType string

// possible values of nats connection type
const (
	Standard  NatsType = "Standard"
	Streaming NatsType = "Streaming"
)

// Event is a data and its context.
// Adheres to the CloudEvents v0.3 specification
type Event struct {
	Context EventContext `json:"context" protobuf:"bytes,1,opt,name=context"`
	Data    []byte       `json:"data" protobuf:"bytes,2,opt,name=data"`
}

// EventContext contains metadata that provides circumstantial information about the occurrence.
type EventContext struct {
	// The type of occurrence which has happened. Often this attribute is used for
	// routing, observability, policy enforcement, etc.
	// should be prefixed with a reverse-DNS name. The prefixed domain dictates
	// the organization which defines the semantics of this event type. ex: com.github.pull.create
	Type string `json:"type" protobuf:"bytes,1,name=type"`
	// The version of the CloudEvents specification which the event uses.
	// Enables the interpretation of the context.
	SpecVersion string `json:"specVersion" protobuf:"bytes,2,name=specVersion"`
	// This describes the event producer.
	Source string `json:"source" protobuf:"bytes,3,name=source"`
	// ID of the event. The semantics are explicitly undefined to ease the implementation of producers.
	ID string `json:"eventID" protobuf:"bytes,4,name=eventID"`
	// Time when the event happened. Must adhere to format specified in RFC 3339.
	Time metav1.MicroTime `json:"time" protobuf:"bytes,5,name=time"`
	// Content type of the data attribute value. Enables the data attribute to carry any type of content,
	// whereby format and encoding might differ from that of the chosen event format.
	// For example, the data attribute may carry an XML or JSON payload and the consumer is informed
	// by this attribute being set to "application/xml" or "application/json" respectively.
	DataContentType string `json:"dataContentType" protobuf:"bytes,6,name=dataContentType"`
	// Subject of the event
	Subject string `json:"subject" protobuf:"bytes,7,name=subject"`
}

// Dispatch protocol contains configuration necessary to dispatch an event to sensor over different communication protocols
type EventProtocol struct {
	Type EventProtocolType `json:"type" protobuf:"bytes,1,name=type"`
	Http Http              `json:"http,omitempty" protobuf:"bytes,2,opt,name=http"`
	Nats Nats              `json:"nats,omitempty" protobuf:"bytes,3,opt,name=nats"`
}

// Http contains the information required to setup a http server and listen to incoming events
type Http struct {
	// Port on which server will run
	Port string `json:"port" protobuf:"bytes,1,name=port"`
}

// Nats contains the information required to connect to nats server and get subscriptions
type Nats struct {
	// URL is nats server/service URL
	ServerURL string `json:"url" protobuf:"bytes,1,name=url"`
	// NATS subject name
	Subject string `json:"subject" protobuf:"bytes,2,name=subject"`
}
