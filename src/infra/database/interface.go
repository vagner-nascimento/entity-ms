package database

import (
	"entity/src/apperrors"
)

type DataBaseHandler interface {
	Insert(data interface{}, table string) (interface{}, *apperrors.Error)
	Get(id interface{}, table string, result interface{}, filters ...map[string]interface{}) *apperrors.Error
	Update(id interface{}, data interface{}, table string, result interface{}, filters ...map[string]interface{}) *apperrors.Error
}
