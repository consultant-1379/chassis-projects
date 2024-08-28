package constvalue

const (

	//---------management service----------
	// MgmtPlmnNRFInfoFormat is the format to mark a home NRF
	MgmtHomeNRFInfoFormat = "NRF Management Service Failed to Connect to NRF in Home PLMN; mcc:%s, mnc:%s, address:%s"

	// MgmtPlmnNRFInfoFormat is the format to mark a PLMN NRF
	MgmtPlmnNRFInfoFormat = "NRF Management Service Failed to Connect to PLMN NRF; address:%s, port:%s"

	// MgmtReplicationConnectionFailureFormat is the format to alarm ReplicationConnectionFailure
	MgmtReplicationConnectionFailureFormat = "NRF Management Service Data Replication Connection to Peer NRF is Broken; NRFInstanceId:%s, FQDN:%s"

	// MgmtRegistrationFailureFormat is the formato to alarm nrfMngtNrfRegistrationFailure
	MgmtRegistrationFailureFormat = "NRF Management Service Failed to Register to PLMN NRF(s); %s"

	//---------discovery service----------
	// DiscRegionNRFInfoFormat is the format to mark a Region NRF
	DiscRegionNRFInfoFormat = "NRF Discovery Service Failed to Connect to Region NRF; %s"

	// DiscHomeNRFInfoFormat is the format to mark a home NRF
	DiscHomeNRFInfoFormat = "NRF Discovery Service Failed to Connect to NRF in Home PLMN; mcc:%s, mnc:%s, address:%s"

	// DiscPlmnNRFInfoFormat is the format to mark a PLMN NRF
	DiscPlmnNRFInfoFormat = "NRF Discovery Service Failed to Connect to PLMN NRF; %s"
)
