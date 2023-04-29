package helpers

import "strings"

func ExtractSubdomain(domain string) (string, string) {
	parts := strings.Split(domain, ".")
	return parts[0], strings.Join(parts[1:], ".")
}

func JoinDomainParts(parts ...string) string {
	return strings.Join(parts, ".")
}
