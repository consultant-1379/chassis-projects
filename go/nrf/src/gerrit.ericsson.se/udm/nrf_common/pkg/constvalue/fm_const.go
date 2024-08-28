package constvalue

const (
	//HomeNRF is for Home NRF
	HomeNRF = "Home NRF"
	// PlmnNRF is for PLMN NRF
	PlmnNRF = "PLMN NRF"
	// RegionNRF is for "Region NRF"
	RegionNRF = "Region NRF"

	// HomeNRFInfoFormat is the format to mark a home NRF
	HomeNRFInfoFormat = "{mcc:%s, mnc:%s, Addr:%s}"

	// PlmnNRFInfoFormat is the format to mark a PLMN NRF
	PlmnNRFInfoFormat = "{Addr:%s}"

	// RegionNRFInfoFormat is the format to mark a Region NRF
	RegionNRFInfoFormat = "{Addr:%s}"
)
