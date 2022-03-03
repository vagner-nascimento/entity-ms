package model

import "entity/src/apperrors"

type Model interface {
	Validate() (bool, []apperrors.Error)
	String() string
}
