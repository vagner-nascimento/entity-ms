package database

import (
	"entity/src/apperrors"
)

type DataBaseHandler interface {
	Insert(data interface{}, target string) (interface{}, *apperrors.Error)
	Get(id interface{}, target string, result interface{}, filters ...map[string]interface{}) *apperrors.Error
	Update(id interface{}, data interface{}, target string, result interface{}, filters ...map[string]interface{}) *apperrors.Error
	Search(parms map[string][]string, target string) (interface{}, *apperrors.Error)
}
