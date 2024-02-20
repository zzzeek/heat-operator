/*

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

package heat

import (
	"github.com/openstack-k8s-operators/lib-common/modules/common/env"

	corev1 "k8s.io/api/core/v1"
)

// APIDetails ..
type APIDetails struct {
	ContainerImage            string
	DatabaseHost              string
	DatabaseName              string
	TransportURL              string
	OSPSecret                 string
	UserPasswordSelector      string
	AuthEncryptionKeySelector string
	VolumeMounts              []corev1.VolumeMount
	Privileged                bool
}

// InitContainerCommand is
const (
	InitContainerCommand = "/usr/local/bin/container-scripts/init.sh"
)

// InitContainer ..
func InitContainer(init APIDetails) []corev1.Container {
	runAsUser := int64(0)
	trueVar := true

	securityContext := &corev1.SecurityContext{
		RunAsUser: &runAsUser,
	}

	if init.Privileged {
		securityContext.Privileged = &trueVar
	}

	envVars := map[string]env.Setter{}
	envVars["DatabaseHost"] = env.SetValue(init.DatabaseHost)
	envVars["DatabaseName"] = env.SetValue(init.DatabaseName)

	envs := []corev1.EnvVar{
		{
			Name: "HeatPassword",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: init.OSPSecret,
					},
					Key: init.UserPasswordSelector,
				},
			},
		},
		{
			Name: "AuthEncryptionKey",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: init.OSPSecret,
					},
					Key: init.AuthEncryptionKeySelector,
				},
			},
		},
	}
	envs = env.MergeEnvs(envs, envVars)

	if init.TransportURL != "" {
		envTransport := corev1.EnvVar{
			Name: "TransportURL",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: init.TransportURL,
					},
					Key: "transport_url",
				},
			},
		}
		envs = append(envs, envTransport)
	}

	return []corev1.Container{
		{
			Name:            "init",
			Image:           init.ContainerImage,
			SecurityContext: securityContext,
			Command: []string{
				"/bin/bash",
			},
			Args: []string{
				InitContainerCommand,
			},
			Env:          envs,
			VolumeMounts: init.VolumeMounts,
		},
	}
}
