package lib

import "flag"

func GetFlags() (port string, sd string) {
	portf := flag.String("port", "8082", "Port to listen on")
	sdf := flag.String("serviceDiscovery", "http://localhost:8500", "Address of service discovery")

	flag.Parse()
	return ":" + *portf, *sdf
}
