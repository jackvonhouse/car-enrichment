package transport

import (
	"net/http"
	"strconv"

	"github.com/jackvonhouse/car-enrichment/internal/errors"
	errpkg "github.com/jackvonhouse/car-enrichment/pkg/errors"
)

func StringToInt(valueStr string) (int, error) {
	valueInt, err := strconv.Atoi(valueStr)

	if err != nil {
		return 0, err
	}

	return valueInt, nil
}

var DefaultErrorHttpCodes = map[uint32]int{
	errors.ErrInternal.TypeId:      http.StatusInternalServerError,
	errors.ErrAlreadyExists.TypeId: http.StatusConflict,
	errors.ErrNotFound.TypeId:      http.StatusNotFound,
	errors.ErrInvalid.TypeId:       http.StatusBadRequest,
	errors.ErrFailed.TypeId:        http.StatusBadRequest,
}

func ErrorToHttpResponse(
	err error, codes map[uint32]int,
) (int, string) {

	if errpkg.Has(err, errors.ErrInternal) {
		return http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)
	}

	code := codes[errpkg.TypeId(err)]

	if code == 0 {
		return http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)
	}

	return code, err.Error()
}
