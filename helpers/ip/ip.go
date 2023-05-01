package ip

type InternetProtocol struct {
	Version string
}

type IPAddress struct {
	Content  string
	Protocol InternetProtocol
}

var IPv4 InternetProtocol = InternetProtocol{Version: "IPv4"}
var IPv6 InternetProtocol = InternetProtocol{Version: "IPv6"}
