# ReceiveMessage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to **NullableString** |  | [optional] 
**MessageId** | Pointer to **NullableString** |  | [optional] 
**PublishTime** | Pointer to **NullableString** |  | [optional] 
**Attributes** | Pointer to **map[string]string** |  | [optional] 

## Methods

### NewReceiveMessage

`func NewReceiveMessage() *ReceiveMessage`

NewReceiveMessage instantiates a new ReceiveMessage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReceiveMessageWithDefaults

`func NewReceiveMessageWithDefaults() *ReceiveMessage`

NewReceiveMessageWithDefaults instantiates a new ReceiveMessage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *ReceiveMessage) GetData() string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ReceiveMessage) GetDataOk() (*string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ReceiveMessage) SetData(v string)`

SetData sets Data field to given value.

### HasData

`func (o *ReceiveMessage) HasData() bool`

HasData returns a boolean if a field has been set.

### SetDataNil

`func (o *ReceiveMessage) SetDataNil(b bool)`

 SetDataNil sets the value for Data to be an explicit nil

### UnsetData
`func (o *ReceiveMessage) UnsetData()`

UnsetData ensures that no value is present for Data, not even an explicit nil
### GetMessageId

`func (o *ReceiveMessage) GetMessageId() string`

GetMessageId returns the MessageId field if non-nil, zero value otherwise.

### GetMessageIdOk

`func (o *ReceiveMessage) GetMessageIdOk() (*string, bool)`

GetMessageIdOk returns a tuple with the MessageId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessageId

`func (o *ReceiveMessage) SetMessageId(v string)`

SetMessageId sets MessageId field to given value.

### HasMessageId

`func (o *ReceiveMessage) HasMessageId() bool`

HasMessageId returns a boolean if a field has been set.

### SetMessageIdNil

`func (o *ReceiveMessage) SetMessageIdNil(b bool)`

 SetMessageIdNil sets the value for MessageId to be an explicit nil

### UnsetMessageId
`func (o *ReceiveMessage) UnsetMessageId()`

UnsetMessageId ensures that no value is present for MessageId, not even an explicit nil
### GetPublishTime

`func (o *ReceiveMessage) GetPublishTime() string`

GetPublishTime returns the PublishTime field if non-nil, zero value otherwise.

### GetPublishTimeOk

`func (o *ReceiveMessage) GetPublishTimeOk() (*string, bool)`

GetPublishTimeOk returns a tuple with the PublishTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublishTime

`func (o *ReceiveMessage) SetPublishTime(v string)`

SetPublishTime sets PublishTime field to given value.

### HasPublishTime

`func (o *ReceiveMessage) HasPublishTime() bool`

HasPublishTime returns a boolean if a field has been set.

### SetPublishTimeNil

`func (o *ReceiveMessage) SetPublishTimeNil(b bool)`

 SetPublishTimeNil sets the value for PublishTime to be an explicit nil

### UnsetPublishTime
`func (o *ReceiveMessage) UnsetPublishTime()`

UnsetPublishTime ensures that no value is present for PublishTime, not even an explicit nil
### GetAttributes

`func (o *ReceiveMessage) GetAttributes() map[string]string`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *ReceiveMessage) GetAttributesOk() (*map[string]string, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *ReceiveMessage) SetAttributes(v map[string]string)`

SetAttributes sets Attributes field to given value.

### HasAttributes

`func (o *ReceiveMessage) HasAttributes() bool`

HasAttributes returns a boolean if a field has been set.

### SetAttributesNil

`func (o *ReceiveMessage) SetAttributesNil(b bool)`

 SetAttributesNil sets the value for Attributes to be an explicit nil

### UnsetAttributes
`func (o *ReceiveMessage) UnsetAttributes()`

UnsetAttributes ensures that no value is present for Attributes, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


