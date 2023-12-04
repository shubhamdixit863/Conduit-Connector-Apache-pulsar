// Copyright © 2022 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
	"testing"
)

var exampleConfig = map[string]string{
	"config_pulsar_url":               "http://localhost",
	"config_pulsar_jwt":               "secret-key-321",
	"config_pulsar_topic":             "topic",
	"config_pulsar_subscription_name": "sub",
}

func TestParseConfig(t *testing.T) {
	is := is.New(t)
	var got Config
	err := sdk.Util.ParseConfig(exampleConfig, &got)
	want := Config{
		ConfigPulsarUrl:              "http://localhost",
		ConfigPulsarJWT:              "secret-key-321",
		ConfigPulsarTopic:            "topic",
		ConfigPulsarSubscriptionName: "sub",
	}
	is.NoErr(err)
	is.Equal(want, got)
}
