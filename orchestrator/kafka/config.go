// SPDX-FileCopyrightText: 2022 Free Mobile
// SPDX-License-Identifier: AGPL-3.0-only

package kafka

import "akvorado/common/kafka"

// Configuration describes the configuration for the Kafka configurator.
type Configuration struct {
	kafka.Configuration `mapstructure:",squash" yaml:",inline"`
	// TopicConfiguration describes the topic configuration.
	TopicConfiguration TopicConfiguration
}

// TopicConfiguration describes the configuration for a topic
type TopicConfiguration struct {
	// NumPartitions tells how many partitions should be used for the topic.
	NumPartitions int32 `validate:"min=1"`
	// ReplicationFactor tells the replication factor for the topic.
	ReplicationFactor int16 `validate:"min=1"`
	// ConfigEntries is a map to specify the topic overrides. Non-listed overrides will be removed
	ConfigEntries map[string]*string
}

// DefaultConfiguration represents the default configuration for the Kafka configurator.
func DefaultConfiguration() Configuration {
	return Configuration{
		Configuration: kafka.DefaultConfiguration(),
		TopicConfiguration: TopicConfiguration{
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
}
