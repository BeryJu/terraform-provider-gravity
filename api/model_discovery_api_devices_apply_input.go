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

// DiscoveryAPIDevicesApplyInput struct for DiscoveryAPIDevicesApplyInput
type DiscoveryAPIDevicesApplyInput struct {
	DhcpScope string `json:"dhcpScope"`
	DnsZone string `json:"dnsZone"`
	To string `json:"to"`
}

// NewDiscoveryAPIDevicesApplyInput instantiates a new DiscoveryAPIDevicesApplyInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiscoveryAPIDevicesApplyInput(dhcpScope string, dnsZone string, to string) *DiscoveryAPIDevicesApplyInput {
	this := DiscoveryAPIDevicesApplyInput{}
	this.DhcpScope = dhcpScope
	this.DnsZone = dnsZone
	this.To = to
	return &this
}

// NewDiscoveryAPIDevicesApplyInputWithDefaults instantiates a new DiscoveryAPIDevicesApplyInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiscoveryAPIDevicesApplyInputWithDefaults() *DiscoveryAPIDevicesApplyInput {
	this := DiscoveryAPIDevicesApplyInput{}
	return &this
}

// GetDhcpScope returns the DhcpScope field value
func (o *DiscoveryAPIDevicesApplyInput) GetDhcpScope() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DhcpScope
}

// GetDhcpScopeOk returns a tuple with the DhcpScope field value
// and a boolean to check if the value has been set.
func (o *DiscoveryAPIDevicesApplyInput) GetDhcpScopeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DhcpScope, true
}

// SetDhcpScope sets field value
func (o *DiscoveryAPIDevicesApplyInput) SetDhcpScope(v string) {
	o.DhcpScope = v
}

// GetDnsZone returns the DnsZone field value
func (o *DiscoveryAPIDevicesApplyInput) GetDnsZone() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DnsZone
}

// GetDnsZoneOk returns a tuple with the DnsZone field value
// and a boolean to check if the value has been set.
func (o *DiscoveryAPIDevicesApplyInput) GetDnsZoneOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DnsZone, true
}

// SetDnsZone sets field value
func (o *DiscoveryAPIDevicesApplyInput) SetDnsZone(v string) {
	o.DnsZone = v
}

// GetTo returns the To field value
func (o *DiscoveryAPIDevicesApplyInput) GetTo() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.To
}

// GetToOk returns a tuple with the To field value
// and a boolean to check if the value has been set.
func (o *DiscoveryAPIDevicesApplyInput) GetToOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.To, true
}

// SetTo sets field value
func (o *DiscoveryAPIDevicesApplyInput) SetTo(v string) {
	o.To = v
}

func (o DiscoveryAPIDevicesApplyInput) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["dhcpScope"] = o.DhcpScope
	}
	if true {
		toSerialize["dnsZone"] = o.DnsZone
	}
	if true {
		toSerialize["to"] = o.To
	}
	return json.Marshal(toSerialize)
}

type NullableDiscoveryAPIDevicesApplyInput struct {
	value *DiscoveryAPIDevicesApplyInput
	isSet bool
}

func (v NullableDiscoveryAPIDevicesApplyInput) Get() *DiscoveryAPIDevicesApplyInput {
	return v.value
}

func (v *NullableDiscoveryAPIDevicesApplyInput) Set(val *DiscoveryAPIDevicesApplyInput) {
	v.value = val
	v.isSet = true
}

func (v NullableDiscoveryAPIDevicesApplyInput) IsSet() bool {
	return v.isSet
}

func (v *NullableDiscoveryAPIDevicesApplyInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDiscoveryAPIDevicesApplyInput(val *DiscoveryAPIDevicesApplyInput) *NullableDiscoveryAPIDevicesApplyInput {
	return &NullableDiscoveryAPIDevicesApplyInput{value: val, isSet: true}
}

func (v NullableDiscoveryAPIDevicesApplyInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDiscoveryAPIDevicesApplyInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


