/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.2.10
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package terraform-provider-gravity

import (
	"encoding/json"
	"time"
)

// BackupBackupStatus struct for BackupBackupStatus
type BackupBackupStatus struct {
	Duration *int32 `json:"duration,omitempty"`
	Error *string `json:"error,omitempty"`
	Filename *string `json:"filename,omitempty"`
	Size *int32 `json:"size,omitempty"`
	Status *string `json:"status,omitempty"`
	Time *time.Time `json:"time,omitempty"`
}

// NewBackupBackupStatus instantiates a new BackupBackupStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupBackupStatus() *BackupBackupStatus {
	this := BackupBackupStatus{}
	return &this
}

// NewBackupBackupStatusWithDefaults instantiates a new BackupBackupStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupBackupStatusWithDefaults() *BackupBackupStatus {
	this := BackupBackupStatus{}
	return &this
}

// GetDuration returns the Duration field value if set, zero value otherwise.
func (o *BackupBackupStatus) GetDuration() int32 {
	if o == nil || o.Duration == nil {
		var ret int32
		return ret
	}
	return *o.Duration
}

// GetDurationOk returns a tuple with the Duration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupBackupStatus) GetDurationOk() (*int32, bool) {
	if o == nil || o.Duration == nil {
		return nil, false
	}
	return o.Duration, true
}

// HasDuration returns a boolean if a field has been set.
func (o *BackupBackupStatus) HasDuration() bool {
	if o != nil && o.Duration != nil {
		return true
	}

	return false
}

// SetDuration gets a reference to the given int32 and assigns it to the Duration field.
func (o *BackupBackupStatus) SetDuration(v int32) {
	o.Duration = &v
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *BackupBackupStatus) GetError() string {
	if o == nil || o.Error == nil {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupBackupStatus) GetErrorOk() (*string, bool) {
	if o == nil || o.Error == nil {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *BackupBackupStatus) HasError() bool {
	if o != nil && o.Error != nil {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *BackupBackupStatus) SetError(v string) {
	o.Error = &v
}

// GetFilename returns the Filename field value if set, zero value otherwise.
func (o *BackupBackupStatus) GetFilename() string {
	if o == nil || o.Filename == nil {
		var ret string
		return ret
	}
	return *o.Filename
}

// GetFilenameOk returns a tuple with the Filename field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupBackupStatus) GetFilenameOk() (*string, bool) {
	if o == nil || o.Filename == nil {
		return nil, false
	}
	return o.Filename, true
}

// HasFilename returns a boolean if a field has been set.
func (o *BackupBackupStatus) HasFilename() bool {
	if o != nil && o.Filename != nil {
		return true
	}

	return false
}

// SetFilename gets a reference to the given string and assigns it to the Filename field.
func (o *BackupBackupStatus) SetFilename(v string) {
	o.Filename = &v
}

// GetSize returns the Size field value if set, zero value otherwise.
func (o *BackupBackupStatus) GetSize() int32 {
	if o == nil || o.Size == nil {
		var ret int32
		return ret
	}
	return *o.Size
}

// GetSizeOk returns a tuple with the Size field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupBackupStatus) GetSizeOk() (*int32, bool) {
	if o == nil || o.Size == nil {
		return nil, false
	}
	return o.Size, true
}

// HasSize returns a boolean if a field has been set.
func (o *BackupBackupStatus) HasSize() bool {
	if o != nil && o.Size != nil {
		return true
	}

	return false
}

// SetSize gets a reference to the given int32 and assigns it to the Size field.
func (o *BackupBackupStatus) SetSize(v int32) {
	o.Size = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *BackupBackupStatus) GetStatus() string {
	if o == nil || o.Status == nil {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupBackupStatus) GetStatusOk() (*string, bool) {
	if o == nil || o.Status == nil {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *BackupBackupStatus) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *BackupBackupStatus) SetStatus(v string) {
	o.Status = &v
}

// GetTime returns the Time field value if set, zero value otherwise.
func (o *BackupBackupStatus) GetTime() time.Time {
	if o == nil || o.Time == nil {
		var ret time.Time
		return ret
	}
	return *o.Time
}

// GetTimeOk returns a tuple with the Time field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupBackupStatus) GetTimeOk() (*time.Time, bool) {
	if o == nil || o.Time == nil {
		return nil, false
	}
	return o.Time, true
}

// HasTime returns a boolean if a field has been set.
func (o *BackupBackupStatus) HasTime() bool {
	if o != nil && o.Time != nil {
		return true
	}

	return false
}

// SetTime gets a reference to the given time.Time and assigns it to the Time field.
func (o *BackupBackupStatus) SetTime(v time.Time) {
	o.Time = &v
}

func (o BackupBackupStatus) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Duration != nil {
		toSerialize["duration"] = o.Duration
	}
	if o.Error != nil {
		toSerialize["error"] = o.Error
	}
	if o.Filename != nil {
		toSerialize["filename"] = o.Filename
	}
	if o.Size != nil {
		toSerialize["size"] = o.Size
	}
	if o.Status != nil {
		toSerialize["status"] = o.Status
	}
	if o.Time != nil {
		toSerialize["time"] = o.Time
	}
	return json.Marshal(toSerialize)
}

type NullableBackupBackupStatus struct {
	value *BackupBackupStatus
	isSet bool
}

func (v NullableBackupBackupStatus) Get() *BackupBackupStatus {
	return v.value
}

func (v *NullableBackupBackupStatus) Set(val *BackupBackupStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableBackupBackupStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableBackupBackupStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBackupBackupStatus(val *BackupBackupStatus) *NullableBackupBackupStatus {
	return &NullableBackupBackupStatus{value: val, isSet: true}
}

func (v NullableBackupBackupStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBackupBackupStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


