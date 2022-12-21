/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.2.10
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package terraform-provider-gravity

import (
	"encoding/json"
)

// TsdbAPIRoleConfigOutput struct for TsdbAPIRoleConfigOutput
type TsdbAPIRoleConfigOutput struct {
	Config TsdbRoleConfig `json:"config"`
}

// NewTsdbAPIRoleConfigOutput instantiates a new TsdbAPIRoleConfigOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTsdbAPIRoleConfigOutput(config TsdbRoleConfig) *TsdbAPIRoleConfigOutput {
	this := TsdbAPIRoleConfigOutput{}
	this.Config = config
	return &this
}

// NewTsdbAPIRoleConfigOutputWithDefaults instantiates a new TsdbAPIRoleConfigOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTsdbAPIRoleConfigOutputWithDefaults() *TsdbAPIRoleConfigOutput {
	this := TsdbAPIRoleConfigOutput{}
	return &this
}

// GetConfig returns the Config field value
func (o *TsdbAPIRoleConfigOutput) GetConfig() TsdbRoleConfig {
	if o == nil {
		var ret TsdbRoleConfig
		return ret
	}

	return o.Config
}

// GetConfigOk returns a tuple with the Config field value
// and a boolean to check if the value has been set.
func (o *TsdbAPIRoleConfigOutput) GetConfigOk() (*TsdbRoleConfig, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Config, true
}

// SetConfig sets field value
func (o *TsdbAPIRoleConfigOutput) SetConfig(v TsdbRoleConfig) {
	o.Config = v
}

func (o TsdbAPIRoleConfigOutput) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["config"] = o.Config
	}
	return json.Marshal(toSerialize)
}

type NullableTsdbAPIRoleConfigOutput struct {
	value *TsdbAPIRoleConfigOutput
	isSet bool
}

func (v NullableTsdbAPIRoleConfigOutput) Get() *TsdbAPIRoleConfigOutput {
	return v.value
}

func (v *NullableTsdbAPIRoleConfigOutput) Set(val *TsdbAPIRoleConfigOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableTsdbAPIRoleConfigOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableTsdbAPIRoleConfigOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTsdbAPIRoleConfigOutput(val *TsdbAPIRoleConfigOutput) *NullableTsdbAPIRoleConfigOutput {
	return &NullableTsdbAPIRoleConfigOutput{value: val, isSet: true}
}

func (v NullableTsdbAPIRoleConfigOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTsdbAPIRoleConfigOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


