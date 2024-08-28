package utils

import (
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

//MaskPrint of personal ids like supi/gpsi
func MaskPrint(in string) string {
	length := len(in)
	switch length {
	case 1:
		return "x"
	case 2:
		return "xx"
	case 3:
		return "xxx"
	case 4:
		return "xxxx"
	case 5, 6, 7, 8:
		return "xxxx" + in[4:]
	default:
		return in[:(length - 8)] + "xxxx" + in[(length - 4):]
	}
}

//MaskPrintbyKey for discovery parameters
func MaskPrintbyKey(key, value string)string {
	if key == constvalue.SearchDataSupi || key == constvalue.SearchDataGpsi {
		return MaskPrint(value)
	}
	return value
}