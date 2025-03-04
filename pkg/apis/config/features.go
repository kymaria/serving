/*
Copyright 2020 The Knative Authors

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

package config

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	cm "knative.dev/pkg/configmap"
)

// Flag is a string value which can be either Enabled, Disabled, or Allowed.
type Flag string

const (
	// FeaturesConfigName is the name of the ConfigMap for the features.
	FeaturesConfigName = "config-features"

	// Enabled turns on an optional behavior.
	Enabled Flag = "Enabled"
	// Disabled turns off an optional behavior.
	Disabled Flag = "Disabled"
	// Allowed neither explicitly disables or enables a behavior.
	// eg. allow a client to control behavior with an annotation or allow a new value through validation.
	Allowed Flag = "Allowed"
)

func defaultFeaturesConfig() *Features {
	return &Features{
		MultiContainer:               Enabled,
		PodSpecAffinity:              Disabled,
		PodSpecDryRun:                Allowed,
		PodSpecHostAliases:           Disabled,
		PodSpecFieldRef:              Disabled,
		PodSpecNodeSelector:          Disabled,
		PodSpecRuntimeClassName:      Disabled,
		PodSpecSecurityContext:       Disabled,
		ContainerSpecAddCapabilities: Disabled,
		PodSpecTolerations:           Disabled,
		TagHeaderBasedRouting:        Disabled,
		AutoDetectHTTP2:              Disabled,
		InjectPodIP:				  Enabled,
	}
}

// NewFeaturesConfigFromMap creates a Features from the supplied Map
func NewFeaturesConfigFromMap(data map[string]string) (*Features, error) {
	nc := defaultFeaturesConfig()

	if err := cm.Parse(data,
		asFlag("multi-container", &nc.MultiContainer),
		asFlag("kubernetes.podspec-affinity", &nc.PodSpecAffinity),
		asFlag("kubernetes.podspec-dryrun", &nc.PodSpecDryRun),
		asFlag("kubernetes.podspec-hostaliases", &nc.PodSpecHostAliases),
		asFlag("kubernetes.podspec-fieldref", &nc.PodSpecFieldRef),
		asFlag("kubernetes.podspec-nodeselector", &nc.PodSpecNodeSelector),
		asFlag("kubernetes.podspec-runtimeclassname", &nc.PodSpecRuntimeClassName),
		asFlag("kubernetes.podspec-securitycontext", &nc.PodSpecSecurityContext),
		asFlag("kubernetes.containerspec-addcapabilities", &nc.ContainerSpecAddCapabilities),
		asFlag("kubernetes.podspec-tolerations", &nc.PodSpecTolerations),
		asFlag("tag-header-based-routing", &nc.TagHeaderBasedRouting),
		asFlag("autodetect-http2", &nc.AutoDetectHTTP2),
		asFlag("inject-pod-ip", &nc.InjectPodIP)); err != nil {
		return nil, err
	}
	return nc, nil
}

// NewFeaturesConfigFromConfigMap creates a Features from the supplied ConfigMap
func NewFeaturesConfigFromConfigMap(config *corev1.ConfigMap) (*Features, error) {
	return NewFeaturesConfigFromMap(config.Data)
}

// Features specifies which features are allowed by the webhook.
type Features struct {
	MultiContainer               Flag
	PodSpecAffinity              Flag
	PodSpecDryRun                Flag
	PodSpecFieldRef              Flag
	PodSpecHostAliases           Flag
	PodSpecNodeSelector          Flag
	PodSpecRuntimeClassName      Flag
	PodSpecSecurityContext       Flag	
	ContainerSpecAddCapabilities Flag
	PodSpecTolerations           Flag
	TagHeaderBasedRouting        Flag
	AutoDetectHTTP2              Flag
	InjectPodIP					 Flag
}

// asFlag parses the value at key as a Flag into the target, if it exists.
func asFlag(key string, target *Flag) cm.ParseFunc {
	return func(data map[string]string) error {
		if raw, ok := data[key]; ok {
			for _, flag := range []Flag{Enabled, Allowed, Disabled} {
				if strings.EqualFold(raw, string(flag)) {
					*target = flag
					return nil
				}
			}
		}
		return nil
	}
}
