# SubscribeRequestData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**EndPoint** | **string** |  | 
**SubscriptionId** | Pointer to **NullableString** |  | [optional] 
**Scopes** | Pointer to **[]string** |  | [optional] 

## Methods

### NewSubscribeRequestData

`func NewSubscribeRequestData(type_ string, endPoint string, ) *SubscribeRequestData`

NewSubscribeRequestData instantiates a new SubscribeRequestData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSubscribeRequestDataWithDefaults

`func NewSubscribeRequestDataWithDefaults() *SubscribeRequestData`

NewSubscribeRequestDataWithDefaults instantiates a new SubscribeRequestData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *SubscribeRequestData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *SubscribeRequestData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *SubscribeRequestData) SetType(v string)`

SetType sets Type field to given value.


### GetEndPoint

`func (o *SubscribeRequestData) GetEndPoint() string`

GetEndPoint returns the EndPoint field if non-nil, zero value otherwise.

### GetEndPointOk

`func (o *SubscribeRequestData) GetEndPointOk() (*string, bool)`

GetEndPointOk returns a tuple with the EndPoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndPoint

`func (o *SubscribeRequestData) SetEndPoint(v string)`

SetEndPoint sets EndPoint field to given value.


### GetSubscriptionId

`func (o *SubscribeRequestData) GetSubscriptionId() string`

GetSubscriptionId returns the SubscriptionId field if non-nil, zero value otherwise.

### GetSubscriptionIdOk

`func (o *SubscribeRequestData) GetSubscriptionIdOk() (*string, bool)`

GetSubscriptionIdOk returns a tuple with the SubscriptionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubscriptionId

`func (o *SubscribeRequestData) SetSubscriptionId(v string)`

SetSubscriptionId sets SubscriptionId field to given value.

### HasSubscriptionId

`func (o *SubscribeRequestData) HasSubscriptionId() bool`

HasSubscriptionId returns a boolean if a field has been set.

### SetSubscriptionIdNil

`func (o *SubscribeRequestData) SetSubscriptionIdNil(b bool)`

 SetSubscriptionIdNil sets the value for SubscriptionId to be an explicit nil

### UnsetSubscriptionId
`func (o *SubscribeRequestData) UnsetSubscriptionId()`

UnsetSubscriptionId ensures that no value is present for SubscriptionId, not even an explicit nil
### GetScopes

`func (o *SubscribeRequestData) GetScopes() []string`

GetScopes returns the Scopes field if non-nil, zero value otherwise.

### GetScopesOk

`func (o *SubscribeRequestData) GetScopesOk() (*[]string, bool)`

GetScopesOk returns a tuple with the Scopes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopes

`func (o *SubscribeRequestData) SetScopes(v []string)`

SetScopes sets Scopes field to given value.

### HasScopes

`func (o *SubscribeRequestData) HasScopes() bool`

HasScopes returns a boolean if a field has been set.

### SetScopesNil

`func (o *SubscribeRequestData) SetScopesNil(b bool)`

 SetScopesNil sets the value for Scopes to be an explicit nil

### UnsetScopes
`func (o *SubscribeRequestData) UnsetScopes()`

UnsetScopes ensures that no value is present for Scopes, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


