package dbmgmt

const (
	// DbGetSuccess indicates getting data successfully from dbproxy
	DbGetSuccess = 2000
	// DbDeleteSuccess indicates deleting data successfully from dbproxy
	DbDeleteSuccess = 2000
	// DbPutSuccess indicates putting data successfully to dbproxy
	DbPutSuccess = 2001
	// DbPutDuplicateSub indicates putting duplication subscription Id to dbproxy, old value is replaced.
	DbPutDuplicateSub = 2003

	// DbDataNotExist indicates data dones't exist in dbproxy when getting and deleting.
	DbDataNotExist = 2002

	// DbDataDbProxyFailed indicates connection between app and dbporxy is failed.
	DbDataDbProxyFailed = 9998

	// DbInvalidData indicates that data is malformed from the response of DB proxy
	DbInvalidData = 9999
)
