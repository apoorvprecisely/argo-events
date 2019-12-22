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

import "github.com/argoproj/argo-events/pkg/apis/sensor"

// EventType is the type of sensor resource change event
type EventType string

// Possible values for EventType
const (
	UpdateEvent EventType = "UPDATE"
	DeleteEvent EventType = "DELETE"
)

// Labels
const (
	//LabelControllerInstanceID is the label which allows to separate application among multiple running controllers.
	LabelControllerInstanceID = sensor.FullName + "/sensor-controller-instanceid"
	// LabelPhase is a label applied to sensors to indicate the current phase of the sensor (for filtering purposes)
	LabelPhase = sensor.FullName + "/phase"
	// LabelComplete is the label to mark sensors as complete
	LabelComplete = sensor.FullName + "/complete"
)
