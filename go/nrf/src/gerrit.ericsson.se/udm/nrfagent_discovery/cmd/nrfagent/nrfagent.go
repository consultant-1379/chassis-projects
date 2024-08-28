package main

import (
	"gerrit.ericsson.se/udm/nrfagent_discovery/app"
)

func init() {

}

func main() {
	serv := app.NewServer()
	defer serv.Stop()
	serv.Run()

	select {
	case <-serv.Terminate:
		break
	}
}
