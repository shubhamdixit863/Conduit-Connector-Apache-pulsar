package source_test

import (
	"context"
	"testing"

	"github.com/matryer/is"

	apachePulsar "github.com/shubhamdixit863/Conduit-Connector-Apache-pulsar/source"
)

func TestTeardownSource_NoOpen(t *testing.T) {
	is := is.New(t)
	con := apachePulsar.NewSource()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}
