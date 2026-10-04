package main

import "regexp"

func scrubPII(message string) string{
	emailPattern := regexp.MustCompile(`[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}`)
	ipPattern := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

	message = emailPattern.ReplaceAllString(message, "[EMAIL]")
	message = ipPattern.ReplaceAllString(message, "[IP]")

	return message
}