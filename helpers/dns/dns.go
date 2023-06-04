package dns

import "dyngo/helpers/ip"

type Record string

const A Record = "A"
const AAAA Record = "AAAA"

func GetRecordForInternetProtocol(protocol ip.InternetProtocol) Record {
	if protocol == ip.IPv4 {
		return A
	} else if protocol == ip.IPv6 {
		return AAAA
	} else {
		panic("Cannot identify record for protocol")
	}
}
