package kafkaconsumer

import (
	"os"
	"strconv"
)

// environment variables declaration
const (
	BootstrapServerEnvVar = "BOOTSTRAP_SERVERS"
	TopicEnvVar           = "TOPIC"
	GroupIDEnvVar         = "GROUP_ID"
	SaslEnvVar            = "SASL"
	SaslUserEnvVar        = "SASL_USER"
	SaslPasswordEnvVar    = "SASL_PASSWORD"
)

// default values for environment variables
const (
	BootstrapServersDefault = "localhost:9092"
	TopicDefault            = "my-topic"
	GroupIDDefault          = "my-group"
	SaslEnabledDefault      = "disabled"
	SaslUserDefault         = "user"
	SaslPasswordDefault     = "password"
)

// ConsumerConfig defines the producer configuration
type ConsumerConfig struct {
	BootstrapServers string
	Topic            string
	GroupID          string
	SaslEnabled      bool
	SaslUser         string
	SaslPassword     string
}

func NewConsumerConfig() *ConsumerConfig {

	saslEnabled := false
	if lookupStringEnv(SaslEnvVar, SaslEnabledDefault) == "enabled" {
		saslEnabled = true
	}

	config := ConsumerConfig{
		BootstrapServers: lookupStringEnv(BootstrapServerEnvVar, BootstrapServersDefault),
		Topic:            lookupStringEnv(TopicEnvVar, TopicDefault),
		GroupID:          lookupStringEnv(GroupIDEnvVar, GroupIDDefault),
		SaslEnabled:      saslEnabled,
		SaslUser:         lookupStringEnv(SaslUserEnvVar, SaslUserDefault),
		SaslPassword:     lookupStringEnv(SaslPasswordEnvVar, SaslPasswordDefault),
	}
	return &config
}

func lookupStringEnv(envVar string, defaultValue string) string {
	envVarValue, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	return envVarValue
}

func lookupInt64Env(envVar string, defaultValue int64) int64 {
	envVarValue, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	int64Val, _ := strconv.ParseInt(envVarValue, 10, 64)
	return int64Val
}
