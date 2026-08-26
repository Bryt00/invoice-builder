package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRegexp = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_` + "`" + `|~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
var NameRegexp = regexp.MustCompile(`^[a-zA-Z\s\-\.\'\,\"]{2,100}$`)
var PhoneRegexp = regexp.MustCompile(`^\+?[0-9\s\-\(\)\.]{7,20}$`)

func IsName(value string) bool {
	return NameRegexp.MatchString(strings.TrimSpace(value))
}

func IsPhone(value string) bool {
	val := strings.TrimSpace(value)
	if val == "" {
		return true
	}
	return PhoneRegexp.MatchString(val)
}

func IsEmail(value string) bool {
	return EmailRegexp.MatchString(strings.TrimSpace(value))
}

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func New() *Validator {
	return &Validator{
		FieldErrors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(msg string) {
	v.NonFieldErrors = append(v.NonFieldErrors, msg)
}
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}
func Matches(value string, regex *regexp.Regexp) bool {
	return regex.MatchString(value)
}

func PermittedInt(value int, permittedValues ...int) bool {
	return slices.Contains(permittedValues, value)
}
