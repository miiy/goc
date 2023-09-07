package validate

import (
	"regexp"
)

// CheckPhone check phone number
func CheckPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

// CheckEmail check email
func CheckEmail(email string) bool {
	pattern := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
