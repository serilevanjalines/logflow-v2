package main

import "regexp"

func scrubPII(message string) string {
	emailPattern := regexp.MustCompile(`[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}`)
	ipPattern := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	cardPattern := regexp.MustCompile(`\b(?:\d[ -]?){13,19}\b`)
	phonePattern := regexp.MustCompile(`(?:\+\d{1,3}[- ]?)?[6-9]\d{9}`)

	message = emailPattern.ReplaceAllString(message, "[EMAIL]")
	message = ipPattern.ReplaceAllString(message, "[IP]")
	message = cardPattern.ReplaceAllString(message, "[CARD]")

	message = phonePattern.ReplaceAllString(message, "[PHONE]")

	return message
}
