package service

import (
	"log"
	"net"
	"time"

	"gerrit.ericsson.se/HSS/5G/cmproxy/service/cmadapter"
	"gerrit.ericsson.se/HSS/5G/cmproxy/service/cmadapter/adpcm"
	"gerrit.ericsson.se/HSS/5G/cmproxy/statistics"
	cmgrpc "gerrit.ericsson.se/HSS/5G/protocols/sidecar/cm"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type data struct {
	Data string `json:"data"`
}

// Server : CM Server
type Server struct {
	cmadapter  cmadapter.Adapter
	grpcServer *grpc.Server
}

// Read : get value by key
func (s *Server) Read(ctx context.Context, in *cmgrpc.CmRequest) (*cmgrpc.CmResponse, error) {

	var err error
	res := cmgrpc.CmResponse{}

	statistics.Statistics.NumberOfAppReads = statistics.Statistics.NumberOfAppReads + 1
	statistics.Statistics.LastAppUpdate = time.Now().String()

	res.Value, err = s.cmadapter.GetValue(in.Key)
	return &res, err
}

// NewServer get new server
func NewServer(cmUri, notifEndpoint string) *Server {

	statistics.Statistics.NumberOfAppReads = 0

	adapter := adpcm.NewAdpCm(cmUri, notifEndpoint)
	if adapter == nil {
		return nil
	}

	s := &Server{cmadapter: adapter, grpcServer: grpc.NewServer()}

	cmgrpc.RegisterCmServiceServer(s.grpcServer, s)

	return s
}

// Start : Start the server
func (s *Server) Start(cmPort string) error {

	log.Println("Listening on", cmPort)

	lis, err := net.Listen("tcp", ":"+cmPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = s.grpcServer.Serve(lis)

	return err

}

// Stop Stop the server
func (s *Server) Stop() {
	s.grpcServer.Stop()
}

// GetCB .
func (s *Server) GetCB() func([]byte) {
	return s.cmadapter.MonitorToReLoad
}
