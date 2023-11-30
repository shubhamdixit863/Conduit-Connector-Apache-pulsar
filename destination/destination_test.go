package destination_test

import (
	"context"
	"testing"

	"github.com/matryer/is"

	apachePulsar "github.com/shubhamdixit863/Conduit-Connector-Apache-pulsar/destination"
)

func TestTeardown_NoOpen(t *testing.T) {
	is := is.New(t)
	con := apachePulsar.NewDestination()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}
