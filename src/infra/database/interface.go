package database

import "entity/src/apperrors"

type DataBaseHandler interface {
	Insert(data interface{}, table string) (id interface{}, err *apperrors.Error)
	Get(id interface{}, table string, result interface{}) (err *apperrors.Error)
}
