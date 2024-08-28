package encoding

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"gerrit.ericsson.se/udm/common/pkg/log"
	kvdbType "gerrit.ericsson.se/udm/nrf_common/pkg/kvdbclient/encoding/schema/kvdbv1"
)

// GFSHCommand struct
type GFSHCommand struct {
        Command string
	ExecutionStatus string
	StatusCode int
	Output string
}

// EncodeGFSHCommand for GFSHCommand
func EncodeGFSHCommand(command GFSHCommand) ([]byte, error) {
	data := kvdbType.TGfshCommand{}

	data.Command = command.Command

	return json.Marshal(data)
}

// DecodeGFSHCommand in json format to golang struct.
func DecodeGFSHCommand(body io.ReadCloser) (*GFSHCommand, error) {
        var err error

        if body == nil {
                err = errors.New("Request body is nil")
                return nil, err
        }

	var data kvdbType.TGfshCommand
	buf, _ := ioutil.ReadAll(body)
        err = json.Unmarshal(buf, &data)
        if err != nil {
                log.Errorf("Fail to decode http request %v", err)
                return nil, err
        }

	result := &GFSHCommand{}
	result.ExecutionStatus = data.ExecutionStatus
	result.StatusCode = data.StatusCode
	result.Output = data.Output

	return result, err
}
