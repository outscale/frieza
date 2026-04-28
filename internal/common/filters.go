package common

import (
	"slices"
)

type FilterKind string

const (
	FilterKindExclude FilterKind = "exclude"
	FilterKindOnly    FilterKind = "only"
)

type ResourceFilterEnvelope struct {
	Kind  FilterKind   `json:"kind"`
	Types []ObjectType `json:"types,omitempty"`
}

func NewResourceFilterExclude(excluded []ObjectType) *ResourceFilterEnvelope {
	return &ResourceFilterEnvelope{
		Kind:  FilterKindExclude,
		Types: excluded,
	}
}

func NewResourceFilterOnly(included []ObjectType) *ResourceFilterEnvelope {
	return &ResourceFilterEnvelope{
		Kind:  FilterKindOnly,
		Types: included,
	}
}

func (f *ResourceFilterEnvelope) Select(typeName ObjectType) bool {
	switch f.Kind {
	case FilterKindExclude:
		return !slices.Contains(f.Types, typeName)
	case FilterKindOnly:
		return slices.Contains(f.Types, typeName)
	default:
		return false
	}
}
