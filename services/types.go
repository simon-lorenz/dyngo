package services

type DynDnsService interface {
	UpdateAllDomains()

	GetDomains() []DynDnsDomain
	GetName() string

	SetTargetIPv4(string)
	SetTargetIPv6(string)
}

type DynDnsDomain struct {
	domain      string
	V4          bool
	V6          bool
	currentIpV4 string
	currentIPv6 string
}
