package destination

//go:generate paramgen -output=paramgen_dest.go DestinationConfig

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/shubhamdixit863/Conduit-Connector-Apache-pulsar/config"
)

type Destination struct {
	sdk.UnimplementedDestination
	producer     pulsar.Producer
	pulsarClient pulsar.Client
	config       DestinationConfig
}

type DestinationConfig struct {
	// Config includes parameters that are the same in the source and destination.
	config.Config
	// DestinationConfigParam must be either yes or no (defaults to yes).
	DestinationConfigParam string `validate:"inclusion=yes|no" default:"yes"`
}

func NewDestination() sdk.Destination {
	// Create Destination and wrap it in the default middleware.
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

func (d *Destination) Parameters() map[string]sdk.Parameter {
	// Parameters is a map of named Parameters that describe how to configure
	// the Destination. Parameters can be generated from DestinationConfig with
	// paramgen.
	return d.config.Parameters()
}

func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	// Configure is the first function to be called in a connector. It provides
	// the connector with the configuration that can be validated and stored.
	// In case the configuration is not valid it should return an error.
	// Testing if your connector can reach the configured data source should be
	// done in Open, not in Configure.
	// The SDK will validate the configuration and populate default values
	// before calling Configure. If you need to do more complex validations you
	// can do them manually here.

	sdk.Logger(ctx).Info().Msg("Configuring Destination...")
	err := sdk.Util.ParseConfig(cfg, &d.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return nil
}

func (d *Destination) Open(ctx context.Context) error {
	// Open is called after Configure to signal the plugin it can prepare to
	// start writing records. If needed, the plugin should open connections in
	// this function.
	// create pulsar client
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:            d.config.ConfigPulsarUrl,
		Authentication: pulsar.NewAuthenticationToken(d.config.ConfigPulsarJWT),
	})
	if err != nil {
		return fmt.Errorf("could not create pulsar client: %v", err)
	}
	d.pulsarClient = client
	// create producer
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: d.config.ConfigPulsarTopic,
	})
	if err != nil {
		return fmt.Errorf("could not create pulsar producer: %v", err)
	}
	d.producer = producer
	return nil
}

func (d *Destination) Write(ctx context.Context, records []sdk.Record) (int, error) {
	// Write writes len(r) records from r to the destination right away without
	// caching. It should return the number of records written from r
	// (0 <= n <= len(r)) and any error encountered that caused the write to
	// stop early. Write must return a non-nil error if it returns n < len(r).
	countMsgs := 0
	for i := 0; i < len(records); i++ {
		_, err := d.producer.Send(ctx, &pulsar.ProducerMessage{
			Payload: records[i].Bytes(),
		})
		if err != nil {
			return countMsgs, err
		}
		countMsgs++
	}
	return countMsgs, nil
}

func (d *Destination) Teardown(ctx context.Context) error {
	// Teardown signals to the plugin that all records were written and there
	// will be no more calls to any other function. After Teardown returns, the
	// plugin should be ready for a graceful shutdown.
	d.producer.Close()
	d.pulsarClient.Close()
	return nil
}
