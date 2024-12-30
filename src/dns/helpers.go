package dns

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func checkForErrors(c CommonResponse) error {
	switch { // TODO: find out if success can be true, and still have errors
	case c.Success == true && len(c.Errors) == 0:
		return nil
	case c.Success == true && len(c.Errors) > 0:
		return fmt.Errorf(
			"operation was successful but got following errors: %v",
			formatErrorsAsString(c.Errors),
		)
	case c.Success == false && len(c.Errors) == 0:
		return fmt.Errorf("operation failed but got no errors")
	default:
		return fmt.Errorf("operation failed: %v", formatErrorsAsString(c.Errors))
	}
}

func formatErrorsAsString(xe []Error) string {
	var s string
	for i := 0; i < len(xe); i++ {
		e := xe[i]
		if i > 0 {
			s += " | " // add a separator for readability
		}
		s += fmt.Sprintf("%v: %v", e.Code, e.Message)
	}
	return s
}

func formatErrors(xe []Error) error {
	return errors.New(formatErrorsAsString(xe))
}

func newRequestWithToken(method, url, token string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Bearer "+token)
	return r, nil
}
