package concourse

import (
	"regexp"
	"strings"
)

var rePrefix = regexp.MustCompile(`^.*(://|@)`)
var reSuffix = regexp.MustCompile(`/$`)
var reFirstChar = regexp.MustCompile(`^[^\p{Ll}\p{Lt}\p{Lm}\p{Lo}]+`)
var reChars = regexp.MustCompile(`[^\p{Ll}\p{Lt}\p{Lm}\p{Lo}\d\-.]+`)

func harmonizeGitUrl(name string) string {
	harmonized := strings.ToLower(name)
	harmonized = rePrefix.ReplaceAllString(harmonized, "")
	harmonized = reSuffix.ReplaceAllString(harmonized, "")
	return harmonizeName(harmonized)
}

func harmonizeName(name string) string {
	harmonized := strings.ToLower(name)
	harmonized = reFirstChar.ReplaceAllString(harmonized, "")
	harmonized = reChars.ReplaceAllString(harmonized, "-")
	return harmonized
}
