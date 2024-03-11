/*
Copyright 2022.

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

//
// Generated by:
//
// operator-sdk create webhook --group heat --version v1beta1 --kind Heat --programmatic-validation --defaulting
//

package v1beta1

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// HeatDefaults -
type HeatDefaults struct {
	APIContainerImageURL    string
	CfnAPIContainerImageURL string
	EngineContainerImageURL string
}

var heatDefaults HeatDefaults

// log is for logging in this package.
var heatlog = logf.Log.WithName("heat-resource")

// SetupHeatDefaults - initialize Heat spec defaults for use with either internal or external webhooks
func SetupHeatDefaults(defaults HeatDefaults) {
	heatDefaults = defaults
	heatlog.Info("Heat defaults initialized", "defaults", defaults)
}

// SetupWebhookWithManager sets up the webhook with the Manager
func (r *Heat) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-heat-openstack-org-v1beta1-heat,mutating=true,failurePolicy=fail,sideEffects=None,groups=heat.openstack.org,resources=heats,verbs=create;update,versions=v1beta1,name=mheat.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Heat{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Heat) Default() {
	heatlog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this Heat spec
func (spec *HeatSpec) Default() {
	if spec.HeatAPI.ContainerImage == "" {
		spec.HeatAPI.ContainerImage = heatDefaults.APIContainerImageURL
	}
	if spec.HeatCfnAPI.ContainerImage == "" {
		spec.HeatCfnAPI.ContainerImage = heatDefaults.CfnAPIContainerImageURL
	}
	if spec.HeatEngine.ContainerImage == "" {
		spec.HeatEngine.ContainerImage = heatDefaults.EngineContainerImageURL
	}
}

// Default - set defaults for this Heat spec core. This version is called
// by the OpenStackControlplane
func (spec *HeatSpecCore) Default() {
	// nothing here yet
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-heat-openstack-org-v1beta1-heat,mutating=false,failurePolicy=fail,sideEffects=None,groups=heat.openstack.org,resources=heats,verbs=create;update,versions=v1beta1,name=vheat.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Heat{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Heat) ValidateCreate() (admission.Warnings, error) {
	heatlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Heat) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	heatlog.Info("validate update", "name", r.Name)
	oldHeat, ok := old.(*Heat)
	if !ok {
		return nil, apierrors.NewInternalError(
			fmt.Errorf("Expected a Heatv1 object, but got %T", oldHeat))
	}

	var errors field.ErrorList
	// We currently have no logic in place to perform database migrations. Changing databases
	// would render all of the existing stacks unmanageable. We should block changes to the
	// databaseInstance to protect existing workloads.
	if r.Spec.DatabaseInstance != oldHeat.Spec.DatabaseInstance {
		errors = append(errors, field.Forbidden(
			field.NewPath("spec.databaseInstance"),
			"Changing the DatabaseInstance is not supported for existing deployments"))
	}

	if errors != nil {
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: "heat.openstack.org", Kind: "Heat"},
			r.Name,
			errors,
		)
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Heat) ValidateDelete() (admission.Warnings, error) {
	heatlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
