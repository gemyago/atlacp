// Code generated by apigen DO NOT EDIT.

package internal

import (
	"time"
	. "github.com/gemyago/atlacp/internal/api/http/v1routes/models"
)

// Below is to workaround unused imports.
var _ = time.Time{}

func NewErrorValidator() FieldValidator[*Error] {
	validateCode := NewSimpleFieldValidator[*interface{}](
		SkipNullValidator(EnsureNonDefault[interface{}]),
	)
	validateMessage := NewSimpleFieldValidator[string](
	)
	
	return func(bindingCtx *BindingContext, value *Error) {
		validateCode(bindingCtx.Fork("code"), value.Code)
		validateMessage(bindingCtx.Fork("message"), value.Message)
	}
}
