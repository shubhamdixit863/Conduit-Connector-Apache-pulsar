package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	apachePulsar "github.com/shubhamdixit863/Conduit-Connector-Apache-pulsar"
)

func main() {
	sdk.Serve(apachePulsar.Connector)
}
