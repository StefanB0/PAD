package lib

import "flag"

func GetFlags() (port string, sd string) {
	portf := flag.String("port", "8083", "Port to listen on")
	sdf := flag.String("serviceDiscovery", "http://localhost:8500", "Address of service discovery")

	flag.Parse()
	return "localhost:" + *portf, *sdf
}
