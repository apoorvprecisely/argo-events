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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Subscription is the definition of a subscription resource
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Subscription struct {
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,name=metadata"`
	metav1.TypeMeta   `json:",inline"`
	Spec              SubscriptionSpec   `json:"spec" protobuf:"bytes,2,name=spec"`
	Status            SubscriptionStatus `json:"status" protobuf:"bytes,3,name=status"`
}

// SubscriptionList is the list of subscription resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SubscriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,name=metadata"`
	// +listType=items
	Items []Subscription `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// SubscriptionSpec describes the specification of the subscription resource
type SubscriptionSpec struct {
	// HTTP refers to list of subscriptions over HTTP protocol
	// +mapType=subscriptions
	HTTP map[string]HTTPSubscription `json:"http,omitempty" protobuf:"bytes,1,rep,name=http"`
	// NATS refers to list of subscriptions over NATS protocol
	// +mapType=subscriptions
	NATS map[string]NATSSubscription `json:"nats,omitempty" protobuf:"bytes,2,rep,name=nats"`
	// AMQP refers to list of subscriptions over AMQP protocol
	// +mapType=subscriptions
	AMQP map[string]AMQPSubscription `json:"amqp,omitempty" protobuf:"bytes,3,rep,name=amqp"`
}

// HTTPSubscription describes the subscription details over HTTP
type HTTPSubscription struct {
	// URL of the event consumer server
	URL string `json:"url" protobuf:"bytes,1,name=url"`
	// Port is the port for the event consumer server to run on.
	// +optional
	Port string `json:"port,omitempty" protobuf:"bytes,1,opt,name=port"`
	// Endpoint is the endpoint to consume events on.
	// +optional
	Endpoint string `json:"endpoint,omitempty" protobuf:"bytes,2,opt,name=endpoint"`
}

// NATSSubscription describes the subscription details over NATS
type NATSSubscription struct {
	// ServerURL is NATS server URL
	ServerURL string `json:"serverURL" protobuf:"bytes,1,name=serverURL"`
	// Subject is the name of the NATS subject
	Subject string `json:"subject" protobuf:"bytes,2,name=subject"`
}

// NATSSubscription describes the subscription details over AMQP
type AMQPSubscription struct {
	// Server URL
	ServerURL string `json:"serverURL" protobuf:"bytes,1,name=serverURL"`
	// Queue name
	Queue string `json:"queue" protobuf:"bytes,2,name=queue"`
	// Credentials if any
	// +optional
	Credentials *Creds `json:"credentials,omitempty" protobuf:"bytes,3,opt,name=credentials"`
}

// Creds contain reference to username and password
type Creds struct {
	Username *corev1.SecretKeySelector `json:"username" protobuf:"bytes,1,name=username"`
	Password *corev1.SecretKeySelector `json:"password" protobuf:"bytes,2,name=password"`
}

// SubscriptionStatus describes the status of the subscription resource
type SubscriptionStatus struct {
	// CreatedAt refers to creation time
	CreatedAt *metav1.Time `json:"createdAt" protobuf:"bytes,1,name=createdAt"`
	// UpdatedAt refers to time at the resource was updated
	UpdatedAt *metav1.Time `json:"updatedAt" protobuf:"bytes,2,name=updatedAt"`
}
