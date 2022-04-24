package services

import "dyngo/config"

type DynDnsService interface {
	UpdateIPv4(string)
	UpdateIPv6(string)
	GetHosts() []config.HostConfiguration
	GetName() string
}
