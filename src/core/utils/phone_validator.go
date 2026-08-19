package utils

import "regexp"

// phoneNumberPattern — O'zbekiston telefon raqami: "998" va undan keyin 9 ta
// raqam. Paket darajasida, chunki funksiya har bir send-otp so'rovida
// chaqiriladi va MustCompile ni qayta-qayta ishga tushirish keraksiz.
var phoneNumberPattern = regexp.MustCompile(`^998[0-9]{9}$`)

// IsValidPhoneNumber raqam kutilgan formatda ekanini bildiradi. Normalizatsiya
// qilmaydi: "+", probel yoki defis bo'lsa rad etadi.
func IsValidPhoneNumber(phone string) bool {
	return phoneNumberPattern.MatchString(phone)
}
