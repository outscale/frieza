package outscale_oapi

import (
	"fmt"

	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

func extractApiError(err error) (bool, *osc.ErrorResponse) {
	genericError, ok := err.(*osc.ErrorResponse)
	if ok {
		return true, genericError
	}
	return false, nil
}

func getErrorInfo(err error) error {
	if ok, apiError := extractApiError(err); ok {
		return fmt.Errorf(
			"%v",
			apiError.GetCode(),
		)
	}

	return err
}
