package probe

import (
	"net/http"
)

func HealthCheck_Handler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}
