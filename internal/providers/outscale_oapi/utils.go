package outscale_oapi

import (
	"fmt"
	"net/http"

	osc "github.com/outscale/osc-sdk-go/v2"
)

func extractApiError(err error) (bool, *osc.ErrorResponse) {
	genericError, ok := err.(osc.GenericOpenAPIError)
	if ok {
		errorsResponse, ok := genericError.Model().(osc.ErrorResponse)
		if ok {
			return true, &errorsResponse
		}
		return false, nil
	}
	return false, nil
}

func getErrorInfo(err error, httpRes *http.Response) error {
	if ok, apiError := extractApiError(err); ok {
		return fmt.Errorf("%v %v %v %v", httpRes.Status, apiError.GetErrors()[0].GetCode(), apiError.GetErrors()[0].GetType(), apiError.GetErrors()[0].GetDetails())
	}
	if httpRes != nil {
		return fmt.Errorf("http status %d", httpRes.Status)
	}

	return err
}
