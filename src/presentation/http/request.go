package http

import (
	"entity/src/apperrors"
	"entity/src/model"
	"entity/src/utils"
	"strings"
)

func getPathParam(path string, skip int) string {
	params := strings.Split(path, "/")
	return params[len(params)-skip]
}

// Validations
const queryAllowedFilters = "lt,gt,lte,gte,after,before"

// Validates the the entity query and returns a slice with the validation errors. If is valid, the response is a empty slice
func validateEntityQuery(q map[string][]string) (res []apperrors.Error) {
	allowFils := model.GetEntityBsonFields()
	allowFils = append(allowFils, "limit")
	allowFils = append(allowFils, "page")
	allowFilts := strings.Split(queryAllowedFilters, ",")
	found := func(sl []string, val string) bool {
		for _, v := range sl {
			if v == val {
				return true
			}
		}

		return false
	}

	for k, _ := range q {
		fild := strings.Split(k, "[")[0]
		filt := utils.StringBetweenBrackets(k)

		if !found(allowFils, fild) {
			res = append(res, apperrors.NewValidationError("field is not allowed", &fild, nil))
		}

		if filt != "" && !found(allowFilts, filt) {
			res = append(res, apperrors.NewValidationError("filter is not allowed", nil, &filt))
		}
	}

	return
}
