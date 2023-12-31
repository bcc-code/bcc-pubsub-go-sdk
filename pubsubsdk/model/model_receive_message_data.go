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

// ReceiveMessageData struct for ReceiveMessageData
type ReceiveMessageData struct {
	Subscription NullableString `json:"subscription,omitempty"`
	Message *ReceiveMessage `json:"message,omitempty"`
}

// NewReceiveMessageData instantiates a new ReceiveMessageData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReceiveMessageData() *ReceiveMessageData {
	this := ReceiveMessageData{}
	return &this
}

// NewReceiveMessageDataWithDefaults instantiates a new ReceiveMessageData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReceiveMessageDataWithDefaults() *ReceiveMessageData {
	this := ReceiveMessageData{}
	return &this
}

// GetSubscription returns the Subscription field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ReceiveMessageData) GetSubscription() string {
	if o == nil || o.Subscription.Get() == nil {
		var ret string
		return ret
	}
	return *o.Subscription.Get()
}

// GetSubscriptionOk returns a tuple with the Subscription field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReceiveMessageData) GetSubscriptionOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return o.Subscription.Get(), o.Subscription.IsSet()
}

// HasSubscription returns a boolean if a field has been set.
func (o *ReceiveMessageData) HasSubscription() bool {
	if o != nil && o.Subscription.IsSet() {
		return true
	}

	return false
}

// SetSubscription gets a reference to the given NullableString and assigns it to the Subscription field.
func (o *ReceiveMessageData) SetSubscription(v string) {
	o.Subscription.Set(&v)
}
// SetSubscriptionNil sets the value for Subscription to be an explicit nil
func (o *ReceiveMessageData) SetSubscriptionNil() {
	o.Subscription.Set(nil)
}

// UnsetSubscription ensures that no value is present for Subscription, not even an explicit nil
func (o *ReceiveMessageData) UnsetSubscription() {
	o.Subscription.Unset()
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *ReceiveMessageData) GetMessage() ReceiveMessage {
	if o == nil || o.Message == nil {
		var ret ReceiveMessage
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReceiveMessageData) GetMessageOk() (*ReceiveMessage, bool) {
	if o == nil || o.Message == nil {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *ReceiveMessageData) HasMessage() bool {
	if o != nil && o.Message != nil {
		return true
	}

	return false
}

// SetMessage gets a reference to the given ReceiveMessage and assigns it to the Message field.
func (o *ReceiveMessageData) SetMessage(v ReceiveMessage) {
	o.Message = &v
}

func (o ReceiveMessageData) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Subscription.IsSet() {
		toSerialize["subscription"] = o.Subscription.Get()
	}
	if o.Message != nil {
		toSerialize["message"] = o.Message
	}
	return json.Marshal(toSerialize)
}

type NullableReceiveMessageData struct {
	value *ReceiveMessageData
	isSet bool
}

func (v NullableReceiveMessageData) Get() *ReceiveMessageData {
	return v.value
}

func (v *NullableReceiveMessageData) Set(val *ReceiveMessageData) {
	v.value = val
	v.isSet = true
}

func (v NullableReceiveMessageData) IsSet() bool {
	return v.isSet
}

func (v *NullableReceiveMessageData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReceiveMessageData(val *ReceiveMessageData) *NullableReceiveMessageData {
	return &NullableReceiveMessageData{value: val, isSet: true}
}

func (v NullableReceiveMessageData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReceiveMessageData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


