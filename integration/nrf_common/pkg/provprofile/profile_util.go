package provprofile

import (
	"github.com/deckarep/golang-set"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"fmt"
	"strconv"
	"strings"
	"regexp"
	"gerrit.ericsson.se/udm/common/pkg/log"
)


const (
	gpsiMinLen = 2
	imsiMinLen = 5
)

const (

	//Imsi is for imsi prefix for supi pattern
	Imsi = "imsi"
	//Msisdn is for msisdn prefix for gpsi pattern
	Msisdn = "msisdn"
	//ImsiStartEnd is for imsi Start/End Pattern
	ImsiStartEnd = "ImsiStartEnd"
	//GpsiStartEnd is for gpsi Start/End Pattern
	GpsiStartEnd = "GpsiStartEnd"
	//ImsiPattern1 is for supported imsi Pattern format 1
	ImsiPattern1 = "ImsiPattern1"
	//ImsiPattern2 is for supported imsi Pattern format 2
	ImsiPattern2 = "ImsiPattern2"
	//ImsiPattern3 is for supported imsi Pattern format 3
	ImsiPattern3 = "ImsiPattern3"
	//ImsiPattern4 is for compile `^\^` + pfxType + `-\d+`
	ImsiPattern4 = "ImsiPattern4"
	//ImsiPattern5 is for compile `^\^` + pfxType + `-\d{` + minLenStr + `,15}`
	ImsiPattern5 = "ImsiPattern5"
	//GpsiPattern1 is for supported gpsi Pattern format 1
	GpsiPattern1 = "GpsiPattern1"
	//GpsiPattern2 is for supported gpsi Pattern format 2
	GpsiPattern2 = "GpsiPattern2"
	//GpsiPattern3 is for supported gpsi Pattern format 3
	GpsiPattern3 = "GpsiPattern3"
	//GpsiPattern4 is for compile `^\^` + pfxType + `-\d+`
	GpsiPattern4 = "GpsiPattern4"
	//GpsiPattern5 is for compile `^\^` + pfxType + `-\d{` + minLenStr + `,15}`
	GpsiPattern5 = "GpsiPattern5"
	//NbrRangePattern is for number range `\{\d{1,2}\}`
	NbrRangePattern = "NbrRangePattern"
)
//PrefixMapValue prefix search map value
type PrefixMapValue struct {
	Index     int
	ValueInfo string
}

//PrefixBody for save extracting Prefix and lenth
type PrefixBody struct {
	Prefix string
	Lenth  int
}

var (
	//Compile to compile partern into memory
	Compile map[string]*regexp.Regexp
)

func provComplie(expr string) *regexp.Regexp {
	re, err := regexp.Compile(expr)
	if err == nil {
		return re
	}
	log.Warnf("Complie %s fail.", expr)
	return nil
}
//PreComplieRegexp to compile pattern into memory
func PreComplieRegexp() {
	Compile = make(map[string]*regexp.Regexp)
	re1 := provComplie("^[0-9]{3}$")
	Compile[constvalue.Mcc] = re1

	re2 := provComplie("^[0-9]{2,3}$")
	Compile[constvalue.Mnc] = re2

	re3 := provComplie("^[0-9]+$")
	Compile[constvalue.Start] = re3

	re4 := provComplie("^[0-9]+$")
	Compile[constvalue.End] = re4

	imsiMinLenStr := strconv.Itoa(imsiMinLen)
	imsiStartEndPattern := `^\d{` + imsiMinLenStr + `,15}$`
	re5 := provComplie(imsiStartEndPattern)
	Compile[ImsiStartEnd] = re5

	gpsiMinLenStr := strconv.Itoa(gpsiMinLen)
	gpsiStartEndPattern := `^\d{` + gpsiMinLenStr + `,15}$`
	re6 := provComplie(gpsiStartEndPattern)
	Compile[GpsiStartEnd] = re6

	imsiPattern1 := `^\^` + Imsi + `-\d{` + imsiMinLenStr + `,15}\\d\{\d{1,2}\}\$$`
	re7 := provComplie(imsiPattern1)
	Compile[ImsiPattern1] = re7

	imsiPattern2 := `^\^` + Imsi + `-\d{` + imsiMinLenStr + `,15}\\d\*\$$`
	re8 := provComplie(imsiPattern2)
	Compile[ImsiPattern2] = re8

	imsiPattern3 := `^\^` + Imsi + `-\d{` + imsiMinLenStr + `,15}$`
	re9 := provComplie(imsiPattern3)
	Compile[ImsiPattern3] = re9

	imsiPattern4 := `^\^` + Imsi + `-\d+`
	re10 := provComplie(imsiPattern4)
	Compile[ImsiPattern4] = re10

	imsiPattern5 := `^\^` + Imsi + `-\d{` + imsiMinLenStr + `,15}`
	re11 := provComplie(imsiPattern5)
	Compile[ImsiPattern5] = re11

	gpsiPattern1 := `^\^` + Msisdn + `-\d{` + gpsiMinLenStr + `,15}\\d\{\d{1,2}\}\$$`
	re12 := provComplie(gpsiPattern1)
	Compile[GpsiPattern1] = re12

	gpsiPattern2 := `^\^` + Msisdn + `-\d{` + gpsiMinLenStr + `,15}\\d\*\$$`
	re13 := provComplie(gpsiPattern2)
	Compile[GpsiPattern2] = re13

	gpsiPattern3 := `^\^` + Msisdn + `-\d{` + gpsiMinLenStr + `,15}$`
	re14 := provComplie(gpsiPattern3)
	Compile[GpsiPattern3] = re14

	gpsiPattern4 := `^\^` + Msisdn + `-\d+`
	re15 := provComplie(gpsiPattern4)
	Compile[GpsiPattern4] = re15

	gpsiPattern5 := `^\^` + Msisdn + `-\d{` + gpsiMinLenStr + `,15}`
	re16 := provComplie(gpsiPattern5)
	Compile[GpsiPattern5] = re16

	nbrRangePattern := `\{\d{1,2}\}`
	re17 := provComplie(nbrRangePattern)
	Compile[NbrRangePattern] = re17
}

func provAtoi(s string) int{
	i, e := strconv.Atoi(s)
	if e == nil {
		return i
	}
	log.Warnf("strconv Atoi fail :%s", s)
	return 0
}

func processStartEnd(start string, end string, minLen int) (mapset.Set, error) {
	var imsiPrefixBody PrefixBody
	imsiPrefixSet := mapset.NewSet()

	trimSuffixLen := 0
	startLen := len(start)
	endLen := len(end)
	if startLen != endLen {
		err := fmt.Errorf("Range Start lenth is not equal to End lenth,invalid Start is %s, End is %s", start, end)
		return nil, err
	}

	if start[:minLen] != end[:minLen] {
		err := fmt.Errorf("The first %d digits must be the same between Range start and end,invalid Start is %s, End is %s", minLen, start, end)
		return nil, err
	}

	for i := 1; i < startLen; i++ {
		if string(start[startLen-i]) == "0" && string(end[startLen-i]) == "9" {
			trimSuffixLen++
		} else {
			break
		}
	}
	imsiPrefixBody.Lenth = startLen

	rawStartPrefix := string(start[0:(startLen - trimSuffixLen)])
	rawEndPrefix := string(end[0:(startLen - trimSuffixLen)])

	rawSamePrefixLen := 0
	rawStartPrefixLen := len(rawStartPrefix)

	for j := 0; j < rawStartPrefixLen; j++ {
		if rawStartPrefix[j] == rawEndPrefix[j] {
			rawSamePrefixLen++
		} else {
			break
		}
	}
	prefixDiffLen := rawStartPrefixLen - rawSamePrefixLen

	startDiffNum := string(rawStartPrefix[rawSamePrefixLen:])
	endDiffNum := string(rawEndPrefix[rawSamePrefixLen:])
	commonImsiPrefix := string(rawStartPrefix[:rawSamePrefixLen])

	if prefixDiffLen == 0 {
		imsiPrefixBody.Prefix = commonImsiPrefix
		imsiPrefixSet.Add(imsiPrefixBody)
	} else if prefixDiffLen == 1 {
		startItem := provAtoi(string(startDiffNum[:prefixDiffLen]))
		endItem := provAtoi(string(endDiffNum[:prefixDiffLen]))
		diffNum := endItem - startItem
		for j := 0; j <= diffNum; j++ {
			element := startItem + j
			imsiPrefixBody.Prefix = commonImsiPrefix + strconv.Itoa(element)
			imsiPrefixSet.Add(imsiPrefixBody)
		}
	} else {
		startTailAggregationFlag := true
		endTailAggregationFlag := true
		for i := 0; i < prefixDiffLen; i++ {
			if i != prefixDiffLen-1 {
				endItem := provAtoi(string(endDiffNum[:i+1]))
				startNum := provAtoi(string(startDiffNum[i]))
				endNum := provAtoi(string(endDiffNum[i]))

				if i == 0 {
					for j := startNum + 1; j < endNum; j++ {
						element := strconv.Itoa(j)
						imsiPrefixBody.Prefix = commonImsiPrefix + element
						imsiPrefixSet.Add(imsiPrefixBody)
					}

					if strings.Count(string(startDiffNum[i+1:]), "0") != prefixDiffLen-i-1 {
						startTailAggregationFlag = false
					}

					if strings.Count(string(endDiffNum[i+1:]), "9") != prefixDiffLen-i-1 {
						endTailAggregationFlag = false
					}

					if startTailAggregationFlag {
						element := strconv.Itoa(startNum)
						imsiPrefixBody.Prefix = commonImsiPrefix + element
						imsiPrefixSet.Add(imsiPrefixBody)
					}
					if endTailAggregationFlag {
						element := strconv.Itoa(endNum)
						imsiPrefixBody.Prefix = commonImsiPrefix + element
						imsiPrefixSet.Add(imsiPrefixBody)
					}
				} else {
					if !startTailAggregationFlag {
						for j := startNum + 1; j <= 9; j++ {
							element := string(startDiffNum[:i]) + strconv.Itoa(j)
							imsiPrefixBody.Prefix = commonImsiPrefix + element
							imsiPrefixSet.Add(imsiPrefixBody)
						}

						if strings.Count(string(startDiffNum[i+1:]), "0") == prefixDiffLen-i-1 {
							startTailAggregationFlag = true
						} else {
							startTailAggregationFlag = false
						}

						if startTailAggregationFlag {
							element := string(startDiffNum[:i+1])
							imsiPrefixBody.Prefix = commonImsiPrefix + element
							imsiPrefixSet.Add(imsiPrefixBody)
						}
					}
					if !endTailAggregationFlag {
						for j := endItem - endNum; j < endItem; j++ {
							element := strconv.Itoa(j)
							imsiPrefixBody.Prefix = commonImsiPrefix + element
							imsiPrefixSet.Add(imsiPrefixBody)
						}

						if strings.Count(string(endDiffNum[i+1:]), "9") == prefixDiffLen-i-1 {
							endTailAggregationFlag = true
						} else {
							endTailAggregationFlag = false
						}

						if endTailAggregationFlag {
							element := string(endDiffNum[:i+1])
							imsiPrefixBody.Prefix = commonImsiPrefix + element
							imsiPrefixSet.Add(imsiPrefixBody)
						}
					}

				}
			} else {
				if !startTailAggregationFlag {
					startItem := provAtoi(string(startDiffNum[i]))
					for j := startItem; j <= 9; j++ {
						startImsiPrefix := string(startDiffNum[:i]) + strconv.Itoa(j)
						imsiPrefixBody.Prefix = commonImsiPrefix + startImsiPrefix
						imsiPrefixSet.Add(imsiPrefixBody)
					}
				}
				if !endTailAggregationFlag {
					endItem  :=  provAtoi(string(endDiffNum[i]))
					for j := endItem; j >= 0; j-- {
						endDiffNumStr := endDiffNum[:i] + strconv.Itoa(j)
						imsiPrefixBody.Prefix = commonImsiPrefix + endDiffNumStr
						imsiPrefixSet.Add(imsiPrefixBody)
					}
				}
			}
		}
	}

	return imsiPrefixSet, nil

}

func processPattern(pattern string, pfxType string, minLen int) (mapset.Set, error) {
	imsiPrefixSet := mapset.NewSet()

	var imsiPrefixBody PrefixBody

	var locStringOne, locStringTwo, locStringThree string
	var rePattern4, rePattern5 *regexp.Regexp
	if pfxType == Imsi {
		re := Compile[ImsiPattern1]
		locStringOne = re.FindString(pattern)

		re = Compile[ImsiPattern2]
		locStringTwo = re.FindString(pattern)

		re = Compile[ImsiPattern3]
		locStringThree = re.FindString(pattern)

		rePattern4 = Compile[ImsiPattern4]
		rePattern5 = Compile[ImsiPattern5]
	} else if pfxType == Msisdn {
		re := Compile[GpsiPattern1]
		locStringOne = re.FindString(pattern)

		re = Compile[GpsiPattern2]
		locStringTwo = re.FindString(pattern)

		re = Compile[GpsiPattern3]
		locStringThree = re.FindString(pattern)

		rePattern4 = Compile[GpsiPattern4]
		rePattern5 = Compile[GpsiPattern5]
	}
	if locStringOne != "" {
		re := Compile[NbrRangePattern]
		tailLenthNumRawData := re.FindString(pattern)

		tailLenthNumRawData = strings.Trim(tailLenthNumRawData, "{")
		tailLenthNum, err := strconv.Atoi(strings.Trim(tailLenthNumRawData, "}"))
		if err != nil {
			err := fmt.Errorf(`Get Range pattern regexp number failed,invalid pattern is %s`, pattern)
			return nil, err
		}

		imsiRawData := rePattern4.FindString(pattern)
		imsiPrefix := strings.Split(imsiRawData, "-")[1]

		imsiPrefixBody.Prefix = imsiPrefix
		imsiPrefixBody.Lenth = len(imsiPrefix) + tailLenthNum

		if imsiPrefixBody.Lenth > 15 || imsiPrefixBody.Lenth < minLen {
			err := fmt.Errorf(`Range Pattern lenth range from %d-15, invalid pattern is %s`, minLen, pattern)
			return nil, err
		}

		imsiPrefixSet.Add(imsiPrefixBody)
	} else if locStringTwo != "" {
		imsiRawData := rePattern5.FindString(pattern)
		imsiPrefix := strings.Split(imsiRawData, "-")[1]

		imsiPrefixBody.Prefix = imsiPrefix
		imsiPrefixBody.Lenth = 0
		imsiPrefixSet.Add(imsiPrefixBody)
	} else if locStringThree != "" {
		imsiRawData := rePattern4.FindString(pattern)
		imsiPrefix := strings.Split(imsiRawData, "-")[1]

		imsiPrefixBody.Prefix = imsiPrefix
		imsiPrefixBody.Lenth = len(imsiPrefix)
		if imsiPrefixBody.Lenth > 15 || imsiPrefixBody.Lenth < minLen {
			err := fmt.Errorf(`Range Pattern lenth range from %d-15, invalid pattern is %s`, minLen, pattern)
			return nil, err
		}
		imsiPrefixSet.Add(imsiPrefixBody)
	} else {
		if pfxType == Imsi {
			err := fmt.Errorf(`SupiRange Pattern format is invalid, invalid pattern is %s. The correct imsi lenth range from 5-15 and the valid pattern format as follows: ^imsi-46000\d{n}$,(n is number),^imsi-46000\d*$,^imsi-46000`, pattern)
			return nil, err
		} else if pfxType == Msisdn {
			err := fmt.Errorf(`GpsiRange Pattern format is invalid, invalid pattern is %s. The correct gpsi lenth range from 2-15 and the valid pattern format as follows: ^msisdn-86\d{n}$,(n is number),^msisdn-86\d*$,^msisdn-46000`, pattern)
			return nil, err
		}
	}

	return imsiPrefixSet, nil
}
