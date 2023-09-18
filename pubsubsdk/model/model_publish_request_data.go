/*
BCC PubSub API

Event messaging and webhook service

API version: v1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package model

import (
	"encoding/json"
)

// PublishRequestData struct for PublishRequestData
type PublishRequestData struct {
	Type string `json:"type"`
	Message string `json:"message"`
	Time NullableString `json:"time,omitempty"`
	Subject NullableString `json:"subject,omitempty"`
	TargetOrgs []string `json:"targetOrgs,omitempty"`
	RequiredScopes []string `json:"requiredScopes,omitempty"`
}

// NewPublishRequestData instantiates a new PublishRequestData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPublishRequestData(type_ string, message string) *PublishRequestData {
	this := PublishRequestData{}
	this.Type = type_
	this.Message = message
	return &this
}

// NewPublishRequestDataWithDefaults instantiates a new PublishRequestData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPublishRequestDataWithDefaults() *PublishRequestData {
	this := PublishRequestData{}
	return &this
}

// GetType returns the Type field value
func (o *PublishRequestData) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *PublishRequestData) GetTypeOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *PublishRequestData) SetType(v string) {
	o.Type = v
}

// GetMessage returns the Message field value
func (o *PublishRequestData) GetMessage() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Message
}

// GetMessageOk returns a tuple with the Message field value
// and a boolean to check if the value has been set.
func (o *PublishRequestData) GetMessageOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Message, true
}

// SetMessage sets field value
func (o *PublishRequestData) SetMessage(v string) {
	o.Message = v
}

// GetTime returns the Time field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PublishRequestData) GetTime() string {
	if o == nil || o.Time.Get() == nil {
		var ret string
		return ret
	}
	return *o.Time.Get()
}

// GetTimeOk returns a tuple with the Time field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PublishRequestData) GetTimeOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return o.Time.Get(), o.Time.IsSet()
}

// HasTime returns a boolean if a field has been set.
func (o *PublishRequestData) HasTime() bool {
	if o != nil && o.Time.IsSet() {
		return true
	}

	return false
}

// SetTime gets a reference to the given NullableString and assigns it to the Time field.
func (o *PublishRequestData) SetTime(v string) {
	o.Time.Set(&v)
}
// SetTimeNil sets the value for Time to be an explicit nil
func (o *PublishRequestData) SetTimeNil() {
	o.Time.Set(nil)
}

// UnsetTime ensures that no value is present for Time, not even an explicit nil
func (o *PublishRequestData) UnsetTime() {
	o.Time.Unset()
}

// GetSubject returns the Subject field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PublishRequestData) GetSubject() string {
	if o == nil || o.Subject.Get() == nil {
		var ret string
		return ret
	}
	return *o.Subject.Get()
}

// GetSubjectOk returns a tuple with the Subject field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PublishRequestData) GetSubjectOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return o.Subject.Get(), o.Subject.IsSet()
}

// HasSubject returns a boolean if a field has been set.
func (o *PublishRequestData) HasSubject() bool {
	if o != nil && o.Subject.IsSet() {
		return true
	}

	return false
}

// SetSubject gets a reference to the given NullableString and assigns it to the Subject field.
func (o *PublishRequestData) SetSubject(v string) {
	o.Subject.Set(&v)
}
// SetSubjectNil sets the value for Subject to be an explicit nil
func (o *PublishRequestData) SetSubjectNil() {
	o.Subject.Set(nil)
}

// UnsetSubject ensures that no value is present for Subject, not even an explicit nil
func (o *PublishRequestData) UnsetSubject() {
	o.Subject.Unset()
}

// GetTargetOrgs returns the TargetOrgs field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PublishRequestData) GetTargetOrgs() []string {
	if o == nil  {
		var ret []string
		return ret
	}
	return o.TargetOrgs
}

// GetTargetOrgsOk returns a tuple with the TargetOrgs field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PublishRequestData) GetTargetOrgsOk() (*[]string, bool) {
	if o == nil || o.TargetOrgs == nil {
		return nil, false
	}
	return &o.TargetOrgs, true
}

// HasTargetOrgs returns a boolean if a field has been set.
func (o *PublishRequestData) HasTargetOrgs() bool {
	if o != nil && o.TargetOrgs != nil {
		return true
	}

	return false
}

// SetTargetOrgs gets a reference to the given []string and assigns it to the TargetOrgs field.
func (o *PublishRequestData) SetTargetOrgs(v []string) {
	o.TargetOrgs = v
}

// GetRequiredScopes returns the RequiredScopes field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PublishRequestData) GetRequiredScopes() []string {
	if o == nil  {
		var ret []string
		return ret
	}
	return o.RequiredScopes
}

// GetRequiredScopesOk returns a tuple with the RequiredScopes field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PublishRequestData) GetRequiredScopesOk() (*[]string, bool) {
	if o == nil || o.RequiredScopes == nil {
		return nil, false
	}
	return &o.RequiredScopes, true
}

// HasRequiredScopes returns a boolean if a field has been set.
func (o *PublishRequestData) HasRequiredScopes() bool {
	if o != nil && o.RequiredScopes != nil {
		return true
	}

	return false
}

// SetRequiredScopes gets a reference to the given []string and assigns it to the RequiredScopes field.
func (o *PublishRequestData) SetRequiredScopes(v []string) {
	o.RequiredScopes = v
}

func (o PublishRequestData) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["type"] = o.Type
	}
	if true {
		toSerialize["message"] = o.Message
	}
	if o.Time.IsSet() {
		toSerialize["time"] = o.Time.Get()
	}
	if o.Subject.IsSet() {
		toSerialize["subject"] = o.Subject.Get()
	}
	if o.TargetOrgs != nil {
		toSerialize["targetOrgs"] = o.TargetOrgs
	}
	if o.RequiredScopes != nil {
		toSerialize["requiredScopes"] = o.RequiredScopes
	}
	return json.Marshal(toSerialize)
}

type NullablePublishRequestData struct {
	value *PublishRequestData
	isSet bool
}

func (v NullablePublishRequestData) Get() *PublishRequestData {
	return v.value
}

func (v *NullablePublishRequestData) Set(val *PublishRequestData) {
	v.value = val
	v.isSet = true
}

func (v NullablePublishRequestData) IsSet() bool {
	return v.isSet
}

func (v *NullablePublishRequestData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePublishRequestData(val *PublishRequestData) *NullablePublishRequestData {
	return &NullablePublishRequestData{value: val, isSet: true}
}

func (v NullablePublishRequestData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePublishRequestData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

