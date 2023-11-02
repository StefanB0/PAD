package lib

import "flag"

func GetFlags() (map[string]string) {
	dockerf := flag.Bool("docker", false, "Whether to use docker or not")
	host := flag.String("host", "localhost", "Host address")

	flag.Parse()

	flags := make(map[string]string)

	if *dockerf {
		flags["service_discovery"] = "http://discovery:8500"
	} else {
		flags["service_discovery"] = "http://localhost:8500"
	}

	flags["host"] = *host

	return flags
}
