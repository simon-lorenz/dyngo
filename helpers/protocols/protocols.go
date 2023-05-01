package protocols

type InternetProtocol struct {
	Version string
}

var IPv4 InternetProtocol = InternetProtocol{Version: "IPv4"}
var IPv6 InternetProtocol = InternetProtocol{Version: "IPv6"}
