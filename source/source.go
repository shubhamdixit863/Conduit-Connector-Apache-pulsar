package source

//go:generate paramgen -output=paramgen_src.go SourceConfig

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"

	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/shubhamdixit863/Conduit-Connector-Apache-pulsar/config"
)

type Source struct {
	sdk.UnimplementedSource
	consumer         pulsar.Consumer
	pulsarClient     pulsar.Client
	config           SourceConfig
	readMsg          map[string]pulsar.MessageID // this is used to store the msgs that are read such that we can acknowledge them in later stage
	lastPositionRead sdk.Position                //nolint:unused // this is just an example
}

type SourceConfig struct {
	// Config includes parameters that are the same in the source and destination.
	config.Config
}

func NewSource() sdk.Source {
	// Create Source and wrap it in the default middleware.
	return sdk.SourceWithMiddleware(&Source{}, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	// Parameters is a map of named Parameters that describe how to configure
	// the Source. Parameters can be generated from SourceConfig with paramgen.
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	// Configure is the first function to be called in a connector. It provides
	// the connector with the configuration that can be validated and stored.
	// In case the configuration is not valid it should return an error.
	// Testing if your connector can reach the configured data source should be
	// done in Open, not in Configure.
	// The SDK will validate the configuration and populate default values
	// before calling Configure. If you need to do more complex validations you
	// can do them manually here.

	sdk.Logger(ctx).Info().Msg("Configuring Source...")
	err := sdk.Util.ParseConfig(cfg, &s.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	s.config.ConfigPulsarUrl = cfg[config.ConfigPulsarUrl]
	s.config.ConfigPulsarTopic = cfg[config.ConfigPulsarTopic]
	s.config.ConfigPulsarSubscriptionName = cfg[config.ConfigPulsarSubscriptionName]
	s.config.ConfigPulsarJWT = cfg[config.ConfigPulsarJWT]
	return nil
}

func (s *Source) Open(ctx context.Context, pos sdk.Position) error {
	// Open is called after Configure to signal the plugin it can prepare to
	// start producing records. If needed, the plugin should open connections in
	// this function. The position parameter will contain the position of the
	// last record that was successfully processed, Source should therefore
	// start producing records after this position. The context passed to Open
	// will be cancelled once the plugin receives a stop signal from Conduit.
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:            s.config.ConfigPulsarUrl,
		Authentication: pulsar.NewAuthenticationToken(s.config.ConfigPulsarJWT),
	})
	s.pulsarClient = client
	if err != nil {
		return err
	}
	// create consumer
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            s.config.ConfigPulsarTopic,
		SubscriptionName: s.config.ConfigPulsarSubscriptionName,
		Type:             pulsar.Shared,
	})
	if err != nil {
		return err
	}
	s.consumer = consumer

	return nil
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	log.Println("read calledd-------------")

	// Read returns a new Record and is supposed to block until there is either
	// a new record or the context gets cancelled. It can also return the error
	// ErrBackoffRetry to signal to the SDK it should call Read again with a
	// backoff retry.
	// If Read receives a cancelled context or the context is cancelled while
	// Read is running it must stop retrieving new records from the source
	// system and start returning records that have already been buffered. If
	// there are no buffered records left Read must return the context error to
	// signal a graceful stop. If Read returns ErrBackoffRetry while the context
	// is cancelled it will also signal that there are no records left and Read
	// won't be called again.
	// After Read returns an error the function won't be called again (except if
	// the error is ErrBackoffRetry, as mentioned above).
	// Read can be called concurrently with Ack.
	msgChan := make(chan pulsar.Message)
	errChan := make(chan error)
	go func() {
		msg, err := s.consumer.Receive(ctx)
		if err != nil {
			errChan <- err
		}
		msgChan <- msg
	}()
	select {
	case msg := <-msgChan:
		s.readMsg[msg.Key()] = msg.ID()
		return sdk.Util.Source.NewRecordCreate(sdk.Position(msg.Key()), nil, sdk.RawData(msg.ID().String()), sdk.RawData(msg.Payload())), nil
	case err := <-errChan:
		return sdk.Record{}, err
	case <-ctx.Done():
		return sdk.Record{}, nil // if the context is cancelled or timeout return

	}
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	// Ack signals to the implementation that the record with the supplied
	// position was successfully processed. This method might be called after
	// the context of Read is already cancelled, since there might be
	// outstanding acks that need to be delivered. When Teardown is called it is
	// guaranteed there won't be any more calls to Ack.
	// Ack can be called concurrently with Read.
	err := s.consumer.AckID(s.readMsg[string(position)])
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) Teardown(ctx context.Context) error {
	// Teardown signals to the plugin that there will be no more calls to any
	// other function. After Teardown returns, the plugin should be ready for a
	// graceful shutdown.
	s.consumer.Close()
	s.pulsarClient.Close()
	return nil
}
