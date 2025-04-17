package model

import (
	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/rpc/params"
)

type (
	DetailedStatus = params.DetailedStatus
	UnitStatus     = params.UnitStatus
)

type Facility struct {
	Name   string
	Status *params.ApplicationStatus
}

type Action struct {
	Name string
	Spec *action.ActionSpec
}

type FacilityMetadata struct {
	ConfigYAML string
}
