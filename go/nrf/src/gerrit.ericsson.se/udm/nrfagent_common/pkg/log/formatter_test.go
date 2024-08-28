package log

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"testing"
	//"fmt"
)

func init() {
	SetServiceID("ausf_service")
	SetNF("ausf")

	log.SetFormatter(&JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyTime:  "timestamp",
			FieldKeyLevel: "level",
			FieldKeyMsg:   "message",
		}})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func getLevelString(s string) string {
	re := regexp.MustCompile("\"level\":\"([a-z|A-Z]+)\"")
	r := re.FindStringSubmatch(s)
	if r == nil {
		return ""
	}
	//fmt.Println(r)
	return r[1]
}

func TestLogLevelString(t *testing.T) {
	type args struct {
		f func(args ...interface{})
		s string
	}
	tests := []struct {
		name  string
		args  args
		wantR string
	}{
		{"debug", args{log.Debug, "1"}, "DEBUG"},
		{"warn", args{log.Warn, "2"}, "WARN"},
		{"info", args{log.Info, "2"}, "INFO"},
		{"error", args{log.Error, "2"}, "ERROR"},
		//{"fatal", args{log.Fatal, "2"}, "FATAL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			log.SetOutput(w)
			defer log.SetOutput(os.Stdout)

			tt.args.f(tt.args.s)

			e := getLevelString(w.String())
			w.Reset()

			if e != tt.wantR {
				t.Errorf("shall be %v instead of %v", e, tt.wantR)
			}
		})
	}
}
