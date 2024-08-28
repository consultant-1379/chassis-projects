package profileop

import (
	"com/dbproxy/nfmessage/nfprofile"
)

// NFProfileQueryFilter is to make up the filter to query NF Profiles
type NFProfileQueryFilter struct {
	ExpiredTimeStart    uint64
	ExpiredTimeEnd      uint64
	LastUpdateTimeStart uint64
	LastUpdateTimeEnd   uint64
	Provisioned         int
	ProvVersion         *nfprofile.ProvVersion
	QueryList           []QueryKeyValue
}

// QueryKeyValue is to make up the filter to query NF Profiles(e.g. value.helper.nfType = 'AUSF')
type QueryKeyValue struct {
	Key   string
	Value string
}
