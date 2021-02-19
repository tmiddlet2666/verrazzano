// +build !ignore_autogenerated

// Copyright (c) 2021, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	corev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentTemplate) DeepCopyInto(out *DeploymentTemplate) {
	*out = *in
	in.Metadata.DeepCopyInto(&out.Metadata)
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	in.PodSpec.DeepCopyInto(&out.PodSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentTemplate.
func (in *DeploymentTemplate) DeepCopy() *DeploymentTemplate {
	if in == nil {
		return nil
	}
	out := new(DeploymentTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressPath) DeepCopyInto(out *IngressPath) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressPath.
func (in *IngressPath) DeepCopy() *IngressPath {
	if in == nil {
		return nil
	}
	out := new(IngressPath)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressRule) DeepCopyInto(out *IngressRule) {
	*out = *in
	if in.Hosts != nil {
		in, out := &in.Hosts, &out.Hosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]IngressPath, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressRule.
func (in *IngressRule) DeepCopy() *IngressRule {
	if in == nil {
		return nil
	}
	out := new(IngressRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressSecurity) DeepCopyInto(out *IngressSecurity) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSecurity.
func (in *IngressSecurity) DeepCopy() *IngressSecurity {
	if in == nil {
		return nil
	}
	out := new(IngressSecurity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressTrait) DeepCopyInto(out *IngressTrait) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressTrait.
func (in *IngressTrait) DeepCopy() *IngressTrait {
	if in == nil {
		return nil
	}
	out := new(IngressTrait)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IngressTrait) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressTraitList) DeepCopyInto(out *IngressTraitList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IngressTrait, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressTraitList.
func (in *IngressTraitList) DeepCopy() *IngressTraitList {
	if in == nil {
		return nil
	}
	out := new(IngressTraitList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IngressTraitList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressTraitSpec) DeepCopyInto(out *IngressTraitSpec) {
	*out = *in
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]IngressRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.TLS = in.TLS
	out.WorkloadReference = in.WorkloadReference
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressTraitSpec.
func (in *IngressTraitSpec) DeepCopy() *IngressTraitSpec {
	if in == nil {
		return nil
	}
	out := new(IngressTraitSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressTraitStatus) DeepCopyInto(out *IngressTraitStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]corev1alpha1.TypedReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressTraitStatus.
func (in *IngressTraitStatus) DeepCopy() *IngressTraitStatus {
	if in == nil {
		return nil
	}
	out := new(IngressTraitStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingScope) DeepCopyInto(out *LoggingScope) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingScope.
func (in *LoggingScope) DeepCopy() *LoggingScope {
	if in == nil {
		return nil
	}
	out := new(LoggingScope)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LoggingScope) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingScopeList) DeepCopyInto(out *LoggingScopeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LoggingScope, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingScopeList.
func (in *LoggingScopeList) DeepCopy() *LoggingScopeList {
	if in == nil {
		return nil
	}
	out := new(LoggingScopeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LoggingScopeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingScopeSpec) DeepCopyInto(out *LoggingScopeSpec) {
	*out = *in
	if in.WorkloadReferences != nil {
		in, out := &in.WorkloadReferences, &out.WorkloadReferences
		*out = make([]corev1alpha1.TypedReference, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingScopeSpec.
func (in *LoggingScopeSpec) DeepCopy() *LoggingScopeSpec {
	if in == nil {
		return nil
	}
	out := new(LoggingScopeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoggingScopeStatus) DeepCopyInto(out *LoggingScopeStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]QualifiedResourceRelation, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoggingScopeStatus.
func (in *LoggingScopeStatus) DeepCopy() *LoggingScopeStatus {
	if in == nil {
		return nil
	}
	out := new(LoggingScopeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsTrait) DeepCopyInto(out *MetricsTrait) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsTrait.
func (in *MetricsTrait) DeepCopy() *MetricsTrait {
	if in == nil {
		return nil
	}
	out := new(MetricsTrait)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MetricsTrait) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsTraitList) DeepCopyInto(out *MetricsTraitList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MetricsTrait, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsTraitList.
func (in *MetricsTraitList) DeepCopy() *MetricsTraitList {
	if in == nil {
		return nil
	}
	out := new(MetricsTraitList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MetricsTraitList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsTraitSpec) DeepCopyInto(out *MetricsTraitSpec) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int)
		**out = **in
	}
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(string)
		**out = **in
	}
	if in.Secret != nil {
		in, out := &in.Secret, &out.Secret
		*out = new(string)
		**out = **in
	}
	if in.Scraper != nil {
		in, out := &in.Scraper, &out.Scraper
		*out = new(string)
		**out = **in
	}
	out.WorkloadReference = in.WorkloadReference
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsTraitSpec.
func (in *MetricsTraitSpec) DeepCopy() *MetricsTraitSpec {
	if in == nil {
		return nil
	}
	out := new(MetricsTraitSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsTraitStatus) DeepCopyInto(out *MetricsTraitStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]QualifiedResourceRelation, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsTraitStatus.
func (in *MetricsTraitStatus) DeepCopy() *MetricsTraitStatus {
	if in == nil {
		return nil
	}
	out := new(MetricsTraitStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QualifiedResourceRelation) DeepCopyInto(out *QualifiedResourceRelation) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QualifiedResourceRelation.
func (in *QualifiedResourceRelation) DeepCopy() *QualifiedResourceRelation {
	if in == nil {
		return nil
	}
	out := new(QualifiedResourceRelation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoCoherenceWorkload) DeepCopyInto(out *VerrazzanoCoherenceWorkload) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoCoherenceWorkload.
func (in *VerrazzanoCoherenceWorkload) DeepCopy() *VerrazzanoCoherenceWorkload {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoCoherenceWorkload)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VerrazzanoCoherenceWorkload) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoCoherenceWorkloadList) DeepCopyInto(out *VerrazzanoCoherenceWorkloadList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VerrazzanoCoherenceWorkload, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoCoherenceWorkloadList.
func (in *VerrazzanoCoherenceWorkloadList) DeepCopy() *VerrazzanoCoherenceWorkloadList {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoCoherenceWorkloadList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VerrazzanoCoherenceWorkloadList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoCoherenceWorkloadSpec) DeepCopyInto(out *VerrazzanoCoherenceWorkloadSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoCoherenceWorkloadSpec.
func (in *VerrazzanoCoherenceWorkloadSpec) DeepCopy() *VerrazzanoCoherenceWorkloadSpec {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoCoherenceWorkloadSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoCoherenceWorkloadStatus) DeepCopyInto(out *VerrazzanoCoherenceWorkloadStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoCoherenceWorkloadStatus.
func (in *VerrazzanoCoherenceWorkloadStatus) DeepCopy() *VerrazzanoCoherenceWorkloadStatus {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoCoherenceWorkloadStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoHelidonWorkload) DeepCopyInto(out *VerrazzanoHelidonWorkload) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoHelidonWorkload.
func (in *VerrazzanoHelidonWorkload) DeepCopy() *VerrazzanoHelidonWorkload {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoHelidonWorkload)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VerrazzanoHelidonWorkload) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoHelidonWorkloadList) DeepCopyInto(out *VerrazzanoHelidonWorkloadList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VerrazzanoHelidonWorkload, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoHelidonWorkloadList.
func (in *VerrazzanoHelidonWorkloadList) DeepCopy() *VerrazzanoHelidonWorkloadList {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoHelidonWorkloadList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VerrazzanoHelidonWorkloadList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoHelidonWorkloadSpec) DeepCopyInto(out *VerrazzanoHelidonWorkloadSpec) {
	*out = *in
	in.DeploymentTemplate.DeepCopyInto(&out.DeploymentTemplate)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoHelidonWorkloadSpec.
func (in *VerrazzanoHelidonWorkloadSpec) DeepCopy() *VerrazzanoHelidonWorkloadSpec {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoHelidonWorkloadSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoHelidonWorkloadStatus) DeepCopyInto(out *VerrazzanoHelidonWorkloadStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]QualifiedResourceRelation, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoHelidonWorkloadStatus.
func (in *VerrazzanoHelidonWorkloadStatus) DeepCopy() *VerrazzanoHelidonWorkloadStatus {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoHelidonWorkloadStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoWebLogicWorkload) DeepCopyInto(out *VerrazzanoWebLogicWorkload) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoWebLogicWorkload.
func (in *VerrazzanoWebLogicWorkload) DeepCopy() *VerrazzanoWebLogicWorkload {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoWebLogicWorkload)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VerrazzanoWebLogicWorkload) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoWebLogicWorkloadList) DeepCopyInto(out *VerrazzanoWebLogicWorkloadList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VerrazzanoWebLogicWorkload, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoWebLogicWorkloadList.
func (in *VerrazzanoWebLogicWorkloadList) DeepCopy() *VerrazzanoWebLogicWorkloadList {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoWebLogicWorkloadList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VerrazzanoWebLogicWorkloadList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoWebLogicWorkloadSpec) DeepCopyInto(out *VerrazzanoWebLogicWorkloadSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoWebLogicWorkloadSpec.
func (in *VerrazzanoWebLogicWorkloadSpec) DeepCopy() *VerrazzanoWebLogicWorkloadSpec {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoWebLogicWorkloadSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VerrazzanoWebLogicWorkloadStatus) DeepCopyInto(out *VerrazzanoWebLogicWorkloadStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VerrazzanoWebLogicWorkloadStatus.
func (in *VerrazzanoWebLogicWorkloadStatus) DeepCopy() *VerrazzanoWebLogicWorkloadStatus {
	if in == nil {
		return nil
	}
	out := new(VerrazzanoWebLogicWorkloadStatus)
	in.DeepCopyInto(out)
	return out
}
