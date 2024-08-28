package cache

import (
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

type SNssai struct {
	Sst int32  `json:"sst"`
	Sd  string `json:"sd,omitempty"`
}

//type PlmnId struct {
//	Mcc string `json:"mcc"`
//	Mnc string `json:"mnc"`
//}

type Tai struct {
	Plmn structs.PlmnID `json:"plmnID,omitempty"`
	Tac  string         `json:"tac,omitempty"`
}

type Ecgi struct {
	Plmn        structs.PlmnID `json:"plmnID,omitempty"`
	EutraCellId string         `json:"eutraCellId,omitempty"`
}

type Ncgi struct {
	Plmn     structs.PlmnID `json:"plmnID,omitempty"`
	NrCellId string         `json:"nrCellId,omitempty"`
}
