# PublishRequestData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**Message** | **string** |  | 
**Time** | Pointer to **NullableString** |  | [optional] 
**TargetOrgs** | Pointer to **[]string** |  | [optional] 
**RequiredScopes** | Pointer to **[]string** |  | [optional] 

## Methods

### NewPublishRequestData

`func NewPublishRequestData(type_ string, message string, ) *PublishRequestData`

NewPublishRequestData instantiates a new PublishRequestData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPublishRequestDataWithDefaults

`func NewPublishRequestDataWithDefaults() *PublishRequestData`

NewPublishRequestDataWithDefaults instantiates a new PublishRequestData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *PublishRequestData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *PublishRequestData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *PublishRequestData) SetType(v string)`

SetType sets Type field to given value.


### GetMessage

`func (o *PublishRequestData) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *PublishRequestData) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *PublishRequestData) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetTime

`func (o *PublishRequestData) GetTime() string`

GetTime returns the Time field if non-nil, zero value otherwise.

### GetTimeOk

`func (o *PublishRequestData) GetTimeOk() (*string, bool)`

GetTimeOk returns a tuple with the Time field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTime

`func (o *PublishRequestData) SetTime(v string)`

SetTime sets Time field to given value.

### HasTime

`func (o *PublishRequestData) HasTime() bool`

HasTime returns a boolean if a field has been set.

### SetTimeNil

`func (o *PublishRequestData) SetTimeNil(b bool)`

 SetTimeNil sets the value for Time to be an explicit nil

### UnsetTime
`func (o *PublishRequestData) UnsetTime()`

UnsetTime ensures that no value is present for Time, not even an explicit nil
### GetTargetOrgs

`func (o *PublishRequestData) GetTargetOrgs() []string`

GetTargetOrgs returns the TargetOrgs field if non-nil, zero value otherwise.

### GetTargetOrgsOk

`func (o *PublishRequestData) GetTargetOrgsOk() (*[]string, bool)`

GetTargetOrgsOk returns a tuple with the TargetOrgs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetOrgs

`func (o *PublishRequestData) SetTargetOrgs(v []string)`

SetTargetOrgs sets TargetOrgs field to given value.

### HasTargetOrgs

`func (o *PublishRequestData) HasTargetOrgs() bool`

HasTargetOrgs returns a boolean if a field has been set.

### SetTargetOrgsNil

`func (o *PublishRequestData) SetTargetOrgsNil(b bool)`

 SetTargetOrgsNil sets the value for TargetOrgs to be an explicit nil

### UnsetTargetOrgs
`func (o *PublishRequestData) UnsetTargetOrgs()`

UnsetTargetOrgs ensures that no value is present for TargetOrgs, not even an explicit nil
### GetRequiredScopes

`func (o *PublishRequestData) GetRequiredScopes() []string`

GetRequiredScopes returns the RequiredScopes field if non-nil, zero value otherwise.

### GetRequiredScopesOk

`func (o *PublishRequestData) GetRequiredScopesOk() (*[]string, bool)`

GetRequiredScopesOk returns a tuple with the RequiredScopes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequiredScopes

`func (o *PublishRequestData) SetRequiredScopes(v []string)`

SetRequiredScopes sets RequiredScopes field to given value.

### HasRequiredScopes

`func (o *PublishRequestData) HasRequiredScopes() bool`

HasRequiredScopes returns a boolean if a field has been set.

### SetRequiredScopesNil

`func (o *PublishRequestData) SetRequiredScopesNil(b bool)`

 SetRequiredScopesNil sets the value for RequiredScopes to be an explicit nil

### UnsetRequiredScopes
`func (o *PublishRequestData) UnsetRequiredScopes()`

UnsetRequiredScopes ensures that no value is present for RequiredScopes, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


