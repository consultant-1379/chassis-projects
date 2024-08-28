package statistics

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Configuration struct {
	Name                  string `json:"name"`
	LastAppUpdate         string `json:"lastAppUpdate"`
	LastAdpUpdate         string `json:"lastAdpUpdate"`
	NumberOfAppReads      int    `json:"numberOfAppReads"`
	NumberOfAdpReads      int    `json:"numberOfAdpReads"`
	NumberOfKafkaMessages int    `json:"numberOfKafkaMessages"`
}

var Statistics = &Configuration{}

func getConfiguration() *Configuration {
	return Statistics
}

type handler struct {
}

func (s handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		config := getConfiguration()
		b, err := json.Marshal(config)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(b); err != nil {
				log.Println("Write message failed")
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

func StartHttpServer(name string, listeningPort int) {
	Statistics.Name = name

	mux := http.NewServeMux()
	mux.Handle("/api/v1/status", handler{})

	h2s := &http2.Server{}

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(listeningPort),
		Handler: h2c.NewHandler(mux, h2s),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server start failed")
	}
}
