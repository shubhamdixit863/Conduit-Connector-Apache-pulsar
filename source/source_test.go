package source_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	check "github.com/matryer/is"

	"github.com/shubhamdixit863/Conduit-Connector-Apache-pulsar/config"
	apachePulsar "github.com/shubhamdixit863/Conduit-Connector-Apache-pulsar/source"
)

func TestTeardownSource_NoOpen(t *testing.T) {
	setEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	is := check.New(t)
	con := apachePulsar.NewSource()
	integrationConfig, err := parseIntegrationConfig()
	is.NoErr(err)
	err = con.Configure(ctx, integrationConfig)
	is.NoErr(err)
	err = con.Open(ctx, nil)
	is.NoErr(err)
	read, err := con.Read(ctx)
	is.True(len(read.Bytes()) >= 0)
	is.NoErr(err)
	err = con.Teardown(context.Background())
	is.NoErr(err)
}

func setEnv() {

	const (
		pulsarURL        string = "pulsar+ssl://pandio--starter-147.us-east-1.aws.pulsar.pandio.com:6651"
		pulsarJWT        string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwYW5kaW8uY29tOnNhOmNsaXVzcTQyM2c5cnZoY3MyZTAwIn0.FXb5wqsXR22EeFRC4lND6pZRyEjMQDyxEv-CtlHVCaI3ZogcZwLUj2OAi4t-MLoXE_0PwJNO_UGqTWmomjxzSEA3dUbj4KkUqJ_Eg8LhGnVuYZGkSsDGnU5-U19arMwyvHVzmEWgRTV2766NIxOc9fU8P7jWUdQ9sigMv0vRlrfVsRUCAhCepOEWvshv4UMzf4CcFwx3Dm9C6u65Tydjqmi7wLWNH-RdH2WbsmndrOAsBkTZubkU5DTXAijWtPhIDkMB5xbAGwcTTZTwcHshGdWfwT6HyLJcV2_5I2waOt232z0x3ewkzgfnZe7TTMTMmpyXC-fAfd9RqlTAShQScQ"
		pulsarTopic      string = "persistent://self-pcyffh/default/go-simple"
		subscriptionName string = "sample"
	)

	// Set environment variables
	os.Setenv("PULSAR_URL", pulsarURL)
	os.Setenv("PULSAR_JWT", pulsarJWT)
	os.Setenv("PULSAR_TOPIC", pulsarTopic)
	os.Setenv("PULSAR_SUB", subscriptionName)

}
func parseIntegrationConfig() (map[string]string, error) {
	pulsarURL := os.Getenv("PULSAR_URL")

	if pulsarURL == "" {
		return map[string]string{}, errors.New("pulsar URL env var must be set")
	}

	pulsarTopic := os.Getenv("PULSAR_TOPIC")
	if pulsarTopic == "" {
		return map[string]string{}, errors.New("PULSAR_TOPIC env var must be set")
	}

	pulsarJWT := os.Getenv("PULSAR_JWT")
	if pulsarJWT == "" {
		return map[string]string{}, errors.New("PULSAR JWT var must be set")
	}

	pulsarSubName := os.Getenv("PULSAR_SUB")
	if pulsarSubName == "" {
		return map[string]string{}, errors.New("PULSAR Subscription var must be set")
	}

	return map[string]string{
		config.ConfigPulsarUrl:              pulsarURL,
		config.ConfigPulsarTopic:            pulsarTopic,
		config.ConfigPulsarJWT:              pulsarJWT,
		config.ConfigPulsarSubscriptionName: pulsarSubName,
	}, nil
}
