package main

import (
	"fmt"
	"os"
	_ "path/filepath"

	"gerrit.ericsson.se/udm/common/tool/pkg/swagger"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("swagger <spec file> <dest dir> <service> <version>")
		fmt.Println("example: swagger ausf.json ./output ausf v1")
	} else {
		spec := os.Args[1]
		dst := os.Args[2]
		service := os.Args[3]
		version := os.Args[4]
		fmt.Printf("  input:  %v\n", spec)
		fmt.Printf("  output: %v\n", dst+"/"+service+version+"/struct.go")

		_ = swagger.DecodeSpecFile(spec, dst, service, version)
	}
}
