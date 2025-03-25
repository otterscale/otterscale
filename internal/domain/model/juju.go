package model

import "github.com/juju/juju/api/base"

type Environment struct {
	*base.UserModel
	Statuses []base.ModelStatus
}
