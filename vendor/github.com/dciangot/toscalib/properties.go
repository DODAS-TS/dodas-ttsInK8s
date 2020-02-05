/*
Copyright 2015 - Olivier Wulveryck

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package toscalib

import "fmt"

// PropertyDefinition as described in Appendix 5.7:
// A property definition defines a named, typed value and related data
// that can be associated with an entity defined in this specification
// (e.g., Node Types, Relation ship Types, Capability Types, etc.).
// Properties are used by template authors to provide input values to
// TOSCA entities which indicate their “desired state” when they are instantiated.
// The value of a property can be retrieved using the
// get_property function within TOSCA Service Templates
type PropertyDefinition struct {
	// Value is not part of PropertyDefinition but an extension to represent both
	// PropertyDefinition and ParameterDefinition within a single type.
	Value       PropertyAssignment `yaml:"value,omitempty"`
	Type        string             `yaml:"type" json:"type"`                                   // The required data type for the property
	Description string             `yaml:"description,omitempty" json:"description,omitempty"` // The optional description for the property.
	Required    bool               `yaml:"required,omitempty" json:"required,omitempty"`       // An optional key that declares a property as required ( true) or not ( false) Default: true
	Default     string             `yaml:"default,omitempty" json:"default,omitempty"`
	Status      Status             `yaml:"status,omitempty" json:"status,omitempty"`
	Constraints Constraints        `yaml:"constraints,omitempty,flow" json:"constraints,omitempty"`
	EntrySchema interface{}        `yaml:"entry_schema,omitempty" json:"entry_schema,omitempty"`
}

// UnmarshalYAML converts YAML text to a type
func (p *PropertyDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err == nil {
		v := newPAValue(s)
		p.Value = *v
		return nil
	}
	var test2 struct {
		Value       PropertyAssignment     `yaml:"value,omitempty"`
		Type        string                 `yaml:"type" json:"type"`                                   // The required data type for the property
		Description string                 `yaml:"description,omitempty" json:"description,omitempty"` // The optional description for the property.
		Required    bool                   `yaml:"required,omitempty" json:"required,omitempty"`       // An optional key that declares a property as required ( true) or not ( false) Default: true
		Default     string                 `yaml:"default,omitempty" json:"default,omitempty"`
		Status      Status                 `yaml:"status,omitempty" json:"status,omitempty"`
		Constraints Constraints            `yaml:"constraints,omitempty,flow" json:"constraints,omitempty"`
		EntrySchema map[string]interface{} `yaml:"entry_schema,omitempty" json:"entry_schema,omitempty"`
	}
	err = unmarshal(&test2)
	//fmt.Println(test2.Required)
	if err == nil {
		p.Value = test2.Value
		p.Type = test2.Type
		p.Description = test2.Description
		p.Required = test2.Required
		p.Default = test2.Default
		p.Status = test2.Status
		p.Constraints = test2.Constraints
		p.EntrySchema = test2.EntrySchema
		return nil
	}
	var test3 struct {
		Value       PropertyAssignment     `yaml:"value,omitempty"`
		Type        string                 `yaml:"type" json:"type"`                                   // The required data type for the property
		Description string                 `yaml:"description,omitempty" json:"description,omitempty"` // The optional description for the property.
		Required    bool                   `yaml:"required,omitempty" json:"required,omitempty"`       // An optional key that declares a property as required ( true) or not ( false) Default: true
		Default     map[string]string      `yaml:"default,omitempty" json:"default,omitempty"`
		Status      Status                 `yaml:"status,omitempty" json:"status,omitempty"`
		Constraints Constraints            `yaml:"constraints,omitempty,flow" json:"constraints,omitempty"`
		EntrySchema map[string]interface{} `yaml:"entry_schema,omitempty" json:"entry_schema,omitempty"`
	}
	err = unmarshal(&test3)
	//fmt.Println(test2.Required)
	if err == nil {
		p.Value = test3.Value
		p.Type = test3.Type
		p.Description = test3.Description
		p.Required = test3.Required
		p.Default = ""
		p.Status = test3.Status
		p.Constraints = test3.Constraints
		p.EntrySchema = test3.EntrySchema
		return nil
	}
	var test4 struct {
		Value       PropertyAssignment     `yaml:"value,omitempty"`
		Type        string                 `yaml:"type" json:"type"`                                   // The required data type for the property
		Description string                 `yaml:"description,omitempty" json:"description,omitempty"` // The optional description for the property.
		Required    bool                   `yaml:"required,omitempty" json:"required,omitempty"`       // An optional key that declares a property as required ( true) or not ( false) Default: true
		Default     []string               `yaml:"default,omitempty" json:"default,omitempty"`
		Status      Status                 `yaml:"status,omitempty" json:"status,omitempty"`
		Constraints Constraints            `yaml:"constraints,omitempty,flow" json:"constraints,omitempty"`
		EntrySchema map[string]interface{} `yaml:"entry_schema,omitempty" json:"entry_schema,omitempty"`
	}
	err = unmarshal(&test4)
	//fmt.Println(test2.Required)
	if err == nil {
		p.Value = test4.Value
		p.Type = test4.Type
		p.Description = test4.Description
		p.Required = test4.Required
		p.Default = ""
		p.Status = test4.Status
		p.Constraints = test4.Constraints
		p.EntrySchema = test4.EntrySchema
		return nil
	}
	var res interface{}
	_ = unmarshal(&res)
	return fmt.Errorf("Cannot parse Property %v", res)
}

// PropertyAssignment supports Value evaluation
type PropertyAssignment struct {
	Assignment
}

func newPAValue(val interface{}) *PropertyAssignment {
	v := new(PropertyAssignment)
	v.Value = val
	return v
}

func newPA(def PropertyDefinition) *PropertyAssignment {
	if def.Value.Value != nil {
		return &def.Value
	}
	return newPAValue(def.Default)
}
