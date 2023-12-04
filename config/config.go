package config

// Config contains shared config parameters, common to the source and
// destination. If you don't need shared parameters you can entirely remove this
// file.
type Config struct {
	// GlobalConfigParam is named global_config_param_name and needs to be
	// provided by the user.
	ConfigPulsarUrl              string `json:"config_pulsar_url" validate:"required"`
	ConfigPulsarJWT              string `json:"config_pulsar_jwt" validate:"required"`
	ConfigPulsarTopic            string `json:"config_pulsar_topic" validate:"required"`
	ConfigPulsarSubscriptionName string `json:"config_pulsar_subscription_name" validate:"required"`
}

const (
	ConfigPulsarUrl = "config_pulsar_url"

	ConfigPulsarJWT = "config_pulsar_jwt"

	ConfigPulsarTopic = "config_pulsar_topic"

	ConfigPulsarSubscriptionName = "config_pulsar_subscription_name"
)
