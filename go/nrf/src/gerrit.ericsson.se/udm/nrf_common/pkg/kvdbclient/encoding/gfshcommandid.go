package encoding

import (
        "encoding/json"
        "errors"
        "io"
        "io/ioutil"
        "gerrit.ericsson.se/udm/common/pkg/log"
        kvdbType "gerrit.ericsson.se/udm/nrf_common/pkg/kvdbclient/encoding/schema/kvdbv1"
)

// GFSHCommandID struct
type GFSHCommandID struct {
	CommandId string
}

// DecodeGFSHCommandID in json format to golang struct.
func DecodeGFSHCommandID(body io.ReadCloser) (*GFSHCommandID, error) {
        var err error

        if body == nil {
                err = errors.New("Request body is nil")
                return nil, err
        }

        var data kvdbType.TGfshCommandId
        buf, _ := ioutil.ReadAll(body)
        err = json.Unmarshal(buf, &data)
        if err != nil {
                log.Errorf("Fail to decode http request %v", err)
                return nil, err
        }

        result := &GFSHCommandID{}
        result.CommandId = data.CommandId

        return result, err
}
