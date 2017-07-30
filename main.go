package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	wol "github.com/sabhiram/go-wol"
)

type WolService interface {
	Wake(context.Context, string) error
}

type wolService struct {
	broadcastAddress string
	iface            string
}

func (w wolService) Wake(_ context.Context, mac string) error {
	// TODO: validation
	if err := wol.SendMagicPacket(mac, w.broadcastAddress, w.iface); err != nil {
		return err
	}
	return nil
}

// ErrEmpty is returned when input mac address is empty
var ErrEmpty = errors.New("Empty Mac Address")

type wakeRequest struct {
	Mac string `json:"mac"`
}

type wakeResponse struct {
	Result string `json:"result"`
}

func makeWakeEndpoint(svc WolService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(wakeRequest)
		if err := svc.Wake(ctx, req.Mac); err != nil {
			return wakeResponse{err.Error()}, nil
		}
		return wakeResponse{"success"}, nil
	}
}

func main() {
	listenAddress := flag.String("listen-address", "127.0.0.1", "Listen Address")
	broadcastAddress := flag.String("broadcast-address", "10.0.0.255", "Broadcast Address")
	iface := flag.String("iface", "eth0", "Network address sent the magick packet")
	flag.Parse()

	svc := wolService{*broadcastAddress, *iface}

	wakeHandler := httptransport.NewServer(
		makeWakeEndpoint(svc),
		decodeWakeRequest,
		encodeResponse,
	)

	http.Handle("/wake", wakeHandler)
	log.Fatal(http.ListenAndServe(*listenAddress+":3000", nil))
}

func decodeWakeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request wakeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
