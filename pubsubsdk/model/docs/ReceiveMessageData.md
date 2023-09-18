# ReceiveMessageData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Subscription** | Pointer to **NullableString** |  | [optional] 
**Message** | Pointer to [**ReceiveMessage**](ReceiveMessage.md) |  | [optional] 

## Methods

### NewReceiveMessageData

`func NewReceiveMessageData() *ReceiveMessageData`

NewReceiveMessageData instantiates a new ReceiveMessageData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReceiveMessageDataWithDefaults

`func NewReceiveMessageDataWithDefaults() *ReceiveMessageData`

NewReceiveMessageDataWithDefaults instantiates a new ReceiveMessageData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSubscription

`func (o *ReceiveMessageData) GetSubscription() string`

GetSubscription returns the Subscription field if non-nil, zero value otherwise.

### GetSubscriptionOk

`func (o *ReceiveMessageData) GetSubscriptionOk() (*string, bool)`

GetSubscriptionOk returns a tuple with the Subscription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubscription

`func (o *ReceiveMessageData) SetSubscription(v string)`

SetSubscription sets Subscription field to given value.

### HasSubscription

`func (o *ReceiveMessageData) HasSubscription() bool`

HasSubscription returns a boolean if a field has been set.

### SetSubscriptionNil

`func (o *ReceiveMessageData) SetSubscriptionNil(b bool)`

 SetSubscriptionNil sets the value for Subscription to be an explicit nil

### UnsetSubscription
`func (o *ReceiveMessageData) UnsetSubscription()`

UnsetSubscription ensures that no value is present for Subscription, not even an explicit nil
### GetMessage

`func (o *ReceiveMessageData) GetMessage() ReceiveMessage`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ReceiveMessageData) GetMessageOk() (*ReceiveMessage, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ReceiveMessageData) SetMessage(v ReceiveMessage)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ReceiveMessageData) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


