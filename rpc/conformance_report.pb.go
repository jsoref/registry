// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.3
// source: google/cloud/apigeeregistry/v1/style/conformance_report.proto

// (-- api-linter: core::0215::versioned-packages=disabled
//     aip.dev/not-precedent: Support protos for the apigeeregistry.v1 API. --)

package rpc

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ConformanceReport describes how well an API Spec or a series of
// API Specs conform to a specific API Style guide.
type ConformanceReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Artifact identifier. This will be auto-generated based on the id of the
	// StyleGuide used to calculate this.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Artifact kind. May be used in YAML representations to identify the type of
	// this artifact.
	Kind string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	// Full resource name of the StyleGuide artifact which was used
	// to generate this report.
	Styleguide string `protobuf:"bytes,3,opt,name=styleguide,proto3" json:"styleguide,omitempty"`
	// A list of guideline report groups.
	GuidelineReportGroups []*GuidelineReportGroup `protobuf:"bytes,4,rep,name=guideline_report_groups,json=guidelineReportGroups,proto3" json:"guideline_report_groups,omitempty"`
}

func (x *ConformanceReport) Reset() {
	*x = ConformanceReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConformanceReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConformanceReport) ProtoMessage() {}

func (x *ConformanceReport) ProtoReflect() protoreflect.Message {
	mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConformanceReport.ProtoReflect.Descriptor instead.
func (*ConformanceReport) Descriptor() ([]byte, []int) {
	return file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescGZIP(), []int{0}
}

func (x *ConformanceReport) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ConformanceReport) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *ConformanceReport) GetStyleguide() string {
	if x != nil {
		return x.Styleguide
	}
	return ""
}

func (x *ConformanceReport) GetGuidelineReportGroups() []*GuidelineReportGroup {
	if x != nil {
		return x.GuidelineReportGroups
	}
	return nil
}

// GuidelineReport describes how well an API Spec or a series of
// API Specs conform to a guideline within an API Style Guide.
type GuidelineReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID of the guideline that this report pertains to.
	GuidelineId string `protobuf:"bytes,1,opt,name=guideline_id,json=guidelineId,proto3" json:"guideline_id,omitempty"`
	// A list of rule report groups.
	RuleReportGroups []*RuleReportGroup `protobuf:"bytes,2,rep,name=rule_report_groups,json=ruleReportGroups,proto3" json:"rule_report_groups,omitempty"`
}

func (x *GuidelineReport) Reset() {
	*x = GuidelineReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuidelineReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuidelineReport) ProtoMessage() {}

func (x *GuidelineReport) ProtoReflect() protoreflect.Message {
	mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuidelineReport.ProtoReflect.Descriptor instead.
func (*GuidelineReport) Descriptor() ([]byte, []int) {
	return file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescGZIP(), []int{1}
}

func (x *GuidelineReport) GetGuidelineId() string {
	if x != nil {
		return x.GuidelineId
	}
	return ""
}

func (x *GuidelineReport) GetRuleReportGroups() []*RuleReportGroup {
	if x != nil {
		return x.RuleReportGroups
	}
	return nil
}

// RuleReport provides information and feedback on a rule that
// a spec breaches within a guideline on an API Style Guide.
type RuleReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID of the rule that the spec breaches.
	RuleId string `protobuf:"bytes,1,opt,name=rule_id,json=ruleId,proto3" json:"rule_id,omitempty"`
	// Resource name of the spec that the rule was breached on.
	Spec string `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec,omitempty"`
	// File name of the spec that the rule was breached on.
	File string `protobuf:"bytes,3,opt,name=file,proto3" json:"file,omitempty"`
	// A suggestion for resolving the problem.
	Suggestion string `protobuf:"bytes,4,opt,name=suggestion,proto3" json:"suggestion,omitempty"`
	// The location of the problem in the spec file.
	Location *LintLocation `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *RuleReport) Reset() {
	*x = RuleReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuleReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuleReport) ProtoMessage() {}

func (x *RuleReport) ProtoReflect() protoreflect.Message {
	mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuleReport.ProtoReflect.Descriptor instead.
func (*RuleReport) Descriptor() ([]byte, []int) {
	return file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescGZIP(), []int{2}
}

func (x *RuleReport) GetRuleId() string {
	if x != nil {
		return x.RuleId
	}
	return ""
}

func (x *RuleReport) GetSpec() string {
	if x != nil {
		return x.Spec
	}
	return ""
}

func (x *RuleReport) GetFile() string {
	if x != nil {
		return x.File
	}
	return ""
}

func (x *RuleReport) GetSuggestion() string {
	if x != nil {
		return x.Suggestion
	}
	return ""
}

func (x *RuleReport) GetLocation() *LintLocation {
	if x != nil {
		return x.Location
	}
	return nil
}

// GuidelineReportGroup is an abstraction that maps state
// (PROPOSED, ACTIVE, DEPRECATED, DISABLED) to a list of
// guideline reports for guidelines of that state.
type GuidelineReportGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// State of the guidelines in this report group.
	State Guideline_State `protobuf:"varint,1,opt,name=state,proto3,enum=google.cloud.apigeeregistry.v1.style.Guideline_State" json:"state,omitempty"`
	// A list of guideline reports.
	GuidelineReports []*GuidelineReport `protobuf:"bytes,2,rep,name=guideline_reports,json=guidelineReports,proto3" json:"guideline_reports,omitempty"`
}

func (x *GuidelineReportGroup) Reset() {
	*x = GuidelineReportGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GuidelineReportGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GuidelineReportGroup) ProtoMessage() {}

func (x *GuidelineReportGroup) ProtoReflect() protoreflect.Message {
	mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GuidelineReportGroup.ProtoReflect.Descriptor instead.
func (*GuidelineReportGroup) Descriptor() ([]byte, []int) {
	return file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescGZIP(), []int{3}
}

func (x *GuidelineReportGroup) GetState() Guideline_State {
	if x != nil {
		return x.State
	}
	return Guideline_STATE_UNSPECIFIED
}

func (x *GuidelineReportGroup) GetGuidelineReports() []*GuidelineReport {
	if x != nil {
		return x.GuidelineReports
	}
	return nil
}

// RuleReportGroup is an abstraction that maps severity
// (ERROR WARNING, INFO, HINT) to a list of rule reports for
// rules of that severity.
type RuleReportGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Severity of the rules in this report group.
	Severity Rule_Severity `protobuf:"varint,1,opt,name=severity,proto3,enum=google.cloud.apigeeregistry.v1.style.Rule_Severity" json:"severity,omitempty"`
	// A list of rule reports.
	RuleReports []*RuleReport `protobuf:"bytes,2,rep,name=rule_reports,json=ruleReports,proto3" json:"rule_reports,omitempty"`
}

func (x *RuleReportGroup) Reset() {
	*x = RuleReportGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuleReportGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuleReportGroup) ProtoMessage() {}

func (x *RuleReportGroup) ProtoReflect() protoreflect.Message {
	mi := &file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuleReportGroup.ProtoReflect.Descriptor instead.
func (*RuleReportGroup) Descriptor() ([]byte, []int) {
	return file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescGZIP(), []int{4}
}

func (x *RuleReportGroup) GetSeverity() Rule_Severity {
	if x != nil {
		return x.Severity
	}
	return Rule_SEVERITY_UNSPECIFIED
}

func (x *RuleReportGroup) GetRuleReports() []*RuleReport {
	if x != nil {
		return x.RuleReports
	}
	return nil
}

var File_google_cloud_apigeeregistry_v1_style_conformance_report_proto protoreflect.FileDescriptor

var file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDesc = []byte{
	0x0a, 0x3d, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x61,
	0x70, 0x69, 0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31,
	0x2f, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e,
	0x63, 0x65, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x24, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x61, 0x70,
	0x69, 0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e,
	0x73, 0x74, 0x79, 0x6c, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x2f, 0x6c, 0x69, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x36, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x2f, 0x73, 0x74,
	0x79, 0x6c, 0x65, 0x5f, 0x67, 0x75, 0x69, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xd5, 0x01, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x23,
	0x0a, 0x0a, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x67, 0x75, 0x69, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x0a, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x67, 0x75,
	0x69, 0x64, 0x65, 0x12, 0x72, 0x0a, 0x17, 0x67, 0x75, 0x69, 0x64, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x3a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x2e, 0x47, 0x75, 0x69, 0x64,
	0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x52, 0x15, 0x67, 0x75, 0x69, 0x64, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x22, 0x9e, 0x01, 0x0a, 0x0f, 0x47, 0x75, 0x69, 0x64,
	0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x26, 0x0a, 0x0c, 0x67,
	0x75, 0x69, 0x64, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x0b, 0x67, 0x75, 0x69, 0x64, 0x65, 0x6c, 0x69, 0x6e,
	0x65, 0x49, 0x64, 0x12, 0x63, 0x0a, 0x12, 0x72, 0x75, 0x6c, 0x65, 0x5f, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x35, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x61,
	0x70, 0x69, 0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31,
	0x2e, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x10, 0x72, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x22, 0xcc, 0x01, 0x0a, 0x0a, 0x52, 0x75, 0x6c,
	0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1c, 0x0a, 0x07, 0x72, 0x75, 0x6c, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x06, 0x72,
	0x75, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x12, 0x17,
	0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41,
	0x02, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x75, 0x67, 0x67, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x75, 0x67,
	0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4e, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x67, 0x65, 0x65, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x74, 0x79, 0x6c, 0x65,
	0x2e, 0x4c, 0x69, 0x6e, 0x74, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xd1, 0x01, 0x0a, 0x14, 0x47, 0x75, 0x69, 0x64,
	0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x12, 0x50, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x35, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x61,
	0x70, 0x69, 0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31,
	0x2e, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x2e, 0x47, 0x75, 0x69, 0x64, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x05, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x67, 0x0a, 0x11, 0x67, 0x75, 0x69, 0x64, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x5f,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x35, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x61, 0x70, 0x69,
	0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x73,
	0x74, 0x79, 0x6c, 0x65, 0x2e, 0x47, 0x75, 0x69, 0x64, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x10, 0x67, 0x75, 0x69, 0x64, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x22, 0xc1, 0x01, 0x0a, 0x0f,
	0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x54, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x33, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x2e, 0x61, 0x70, 0x69, 0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e,
	0x76, 0x31, 0x2e, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x2e, 0x53, 0x65,
	0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x08, 0x73, 0x65, 0x76,
	0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x58, 0x0a, 0x0c, 0x72, 0x75, 0x6c, 0x65, 0x5f, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x67, 0x65,
	0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x74, 0x79,
	0x6c, 0x65, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x03, 0xe0,
	0x41, 0x02, 0x52, 0x0b, 0x72, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x42,
	0x68, 0x0a, 0x28, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x67, 0x65, 0x65, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x42, 0x16, 0x43, 0x6f, 0x6e,
	0x66, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x70, 0x69, 0x67, 0x65, 0x65, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x2f, 0x72, 0x70, 0x63, 0x3b, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescOnce sync.Once
	file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescData = file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDesc
)

func file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescGZIP() []byte {
	file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescOnce.Do(func() {
		file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescData = protoimpl.X.CompressGZIP(file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescData)
	})
	return file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDescData
}

var file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_goTypes = []interface{}{
	(*ConformanceReport)(nil),    // 0: google.cloud.apigeeregistry.v1.style.ConformanceReport
	(*GuidelineReport)(nil),      // 1: google.cloud.apigeeregistry.v1.style.GuidelineReport
	(*RuleReport)(nil),           // 2: google.cloud.apigeeregistry.v1.style.RuleReport
	(*GuidelineReportGroup)(nil), // 3: google.cloud.apigeeregistry.v1.style.GuidelineReportGroup
	(*RuleReportGroup)(nil),      // 4: google.cloud.apigeeregistry.v1.style.RuleReportGroup
	(*LintLocation)(nil),         // 5: google.cloud.apigeeregistry.v1.style.LintLocation
	(Guideline_State)(0),         // 6: google.cloud.apigeeregistry.v1.style.Guideline.State
	(Rule_Severity)(0),           // 7: google.cloud.apigeeregistry.v1.style.Rule.Severity
}
var file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_depIdxs = []int32{
	3, // 0: google.cloud.apigeeregistry.v1.style.ConformanceReport.guideline_report_groups:type_name -> google.cloud.apigeeregistry.v1.style.GuidelineReportGroup
	4, // 1: google.cloud.apigeeregistry.v1.style.GuidelineReport.rule_report_groups:type_name -> google.cloud.apigeeregistry.v1.style.RuleReportGroup
	5, // 2: google.cloud.apigeeregistry.v1.style.RuleReport.location:type_name -> google.cloud.apigeeregistry.v1.style.LintLocation
	6, // 3: google.cloud.apigeeregistry.v1.style.GuidelineReportGroup.state:type_name -> google.cloud.apigeeregistry.v1.style.Guideline.State
	1, // 4: google.cloud.apigeeregistry.v1.style.GuidelineReportGroup.guideline_reports:type_name -> google.cloud.apigeeregistry.v1.style.GuidelineReport
	7, // 5: google.cloud.apigeeregistry.v1.style.RuleReportGroup.severity:type_name -> google.cloud.apigeeregistry.v1.style.Rule.Severity
	2, // 6: google.cloud.apigeeregistry.v1.style.RuleReportGroup.rule_reports:type_name -> google.cloud.apigeeregistry.v1.style.RuleReport
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_init() }
func file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_init() {
	if File_google_cloud_apigeeregistry_v1_style_conformance_report_proto != nil {
		return
	}
	file_google_cloud_apigeeregistry_v1_style_lint_proto_init()
	file_google_cloud_apigeeregistry_v1_style_style_guide_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConformanceReport); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuidelineReport); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuleReport); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GuidelineReportGroup); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuleReportGroup); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_goTypes,
		DependencyIndexes: file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_depIdxs,
		MessageInfos:      file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_msgTypes,
	}.Build()
	File_google_cloud_apigeeregistry_v1_style_conformance_report_proto = out.File
	file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_rawDesc = nil
	file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_goTypes = nil
	file_google_cloud_apigeeregistry_v1_style_conformance_report_proto_depIdxs = nil
}
