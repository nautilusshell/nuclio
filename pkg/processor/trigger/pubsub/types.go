/*
Copyright 2017 The Nuclio Authors.

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

package pubsub

import (
	"github.com/nuclio/nuclio/pkg/functionconfig"
	"github.com/nuclio/nuclio/pkg/processor/runtime"
	"github.com/nuclio/nuclio/pkg/processor/trigger"

	"github.com/mitchellh/mapstructure"
	"github.com/nuclio/errors"
)

type Subscription struct {
	Topic         string
	MaxNumWorkers int
	Shared        bool
	AckDeadline   string
}

type Configuration struct {
	trigger.Configuration
	Subscriptions []Subscription
	ProjectID     string
	AckDeadline   string
	Credentials   trigger.Secret
}

func NewConfiguration(ID string,
	triggerConfiguration *functionconfig.Trigger,
	runtimeConfiguration *runtime.Configuration) (*Configuration, error) {
	newConfiguration := Configuration{}

	// create base
	newConfiguration.Configuration = *trigger.NewConfiguration(ID, triggerConfiguration, runtimeConfiguration)

	// parse attributes
	if err := mapstructure.Decode(newConfiguration.Configuration.Attributes, &newConfiguration); err != nil {
		return nil, errors.Wrap(err, "Failed to decode attributes")
	}

	for subscriptionIdx, subscriptions := range newConfiguration.Subscriptions {

		if subscriptions.Topic == "" {
			return nil, errors.New("Subscription topic must be set")
		}

		if subscriptions.MaxNumWorkers == 0 {
			newConfiguration.Subscriptions[subscriptionIdx].MaxNumWorkers = 1
		}
	}

	return &newConfiguration, nil
}
