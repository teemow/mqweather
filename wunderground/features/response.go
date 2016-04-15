package features

import "fmt"

type Response struct {
	Version        string       `json:"version"`
	TermsOfService string       `json:"termsofService"`
	Error          *ErrResponse `json:"error,omitempty"`
}

type ErrResponse struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (r *Response) HasError() error {
	if r.Error == nil {
		return nil
	}

	return fmt.Errorf("%s: %s", r.Error.Type, r.Error.Description)
}
