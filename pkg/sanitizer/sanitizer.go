package sanitizer

import "github.com/microcosm-cc/bluemonday"

type Sanitizer struct {
	sanitizer *bluemonday.Policy
}

func NewSanitizer(sz *bluemonday.Policy) *Sanitizer {
	customSanitizer := Sanitizer{sanitizer: sz}
	return &customSanitizer
}
