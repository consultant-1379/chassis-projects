package http2

import (
	"net/http"
)

// SetT1 set http2.transport setting by http.transport
func (t *Transport) SetT1(transport *http.Transport) {
	t.t1 = transport
}

// SetT1 get http2.transport setting by http.transport
func (t *Transport) GetT1() *http.Transport {
	return t.t1
}
