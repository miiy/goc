package strings

import "fmt"

// MaskIdCard Hide sensitive information on ID Card
func MaskIdCard(id string) string {
	if len(id) != 18 {
		return id
	}
	return fmt.Sprintf("%s********%s", id[:6], id[14:])
}

// MaskPhone  Hide sensitive information on phone number
func MaskPhone(phone string) string {
	if len(phone) == 11 {
		return fmt.Sprintf("%s****%s", phone[:3], phone[len(phone)-4:])
	} else {
		return phone
	}
}
