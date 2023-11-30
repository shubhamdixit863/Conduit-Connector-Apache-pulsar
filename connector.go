package connectorname

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/shubhamdixit863/Conduit-Connector-Apache-pulsar/destination"
	"github.com/shubhamdixit863/Conduit-Connector-Apache-pulsar/source"
)

// Connector combines all constructors for each plugin in one struct.
var Connector = sdk.Connector{
	NewSpecification: Specification,
	NewSource:        source.NewSource,
	NewDestination:   destination.NewDestination,
}
