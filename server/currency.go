package server

import (
	"context"
	"io"
	"time"

	"grpc-coba/data"
	protos "grpc-coba/protos/currency"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	rates        *data.ExchangeRates
	log          hclog.Logger
	subscription map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {

	go func() {
		ru := r.MonitorRates(5 * time.Second)
		for range ru {
			l.Info("Got updated rates")
		}
	}()

	return &Currency{r, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	go func() {
		for {
			rr, err := src.Recv()
			if err == io.EOF {
				c.log.Info("Client has closed connection")
				break
			}

			if err != nil {
				c.log.Error("Unable to read from client", "error", err)
			}

			c.log.Info("Handle client request", "request", rr)

			rrs, ok := c.subscription[src]
			if !ok {
				rrs = []*protos.RateRequest{}
			}

			rrs = append(rrs, rr)
			c.subscription[src] = rrs
		}
	}()

	for {
		err := src.Send(&protos.RateResponse{Rate: 12.1})
		if err != nil {
			return err
		}

		time.Sleep(5 * time.Second)
	}
}
