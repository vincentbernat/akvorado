// SPDX-FileCopyrightText: 2022 Free Mobile
// SPDX-License-Identifier: AGPL-3.0-only

package snmp

import (
	"testing"
	"time"

	"akvorado/common/helpers"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func TestDefaultConfiguration(t *testing.T) {
	if err := helpers.Validate.Struct(DefaultConfiguration()); err != nil {
		t.Fatalf("validate.Struct() error:\n%+v", err)
	}
}

func TestConfigurationUnmarshallerHook(t *testing.T) {
	cases := []struct {
		Description string
		Input       gin.H
		Output      Configuration
	}{
		{
			Description: "nil",
			Input:       nil,
		}, {
			Description: "empty",
			Input:       gin.H{},
		}, {
			Description: "no communities, no default community",
			Input: gin.H{
				"cache-refresh":  "10s",
				"poller-retries": 10,
			},
			Output: Configuration{
				CacheRefresh:  10 * time.Second,
				PollerRetries: 10,
			},
		}, {
			Description: "communities, no default community",
			Input: gin.H{
				"communities": gin.H{
					"203.0.113.0/25":   "public",
					"203.0.113.128/25": "private",
				},
			},
			Output: Configuration{
				Communities: helpers.MustNewSubnetMap(map[string]string{
					"::ffff:203.0.113.0/121":   "public",
					"::ffff:203.0.113.128/121": "private",
				}),
			},
		}, {
			Description: "no communities, default community",
			Input: gin.H{
				"default-community": "private",
			},
			Output: Configuration{
				Communities: helpers.MustNewSubnetMap(map[string]string{
					"::/0": "private",
				}),
			},
		}, {
			Description: "communities, default community",
			Input: gin.H{
				"default-community": "private",
				"communities": gin.H{
					"203.0.113.0/25":   "public",
					"203.0.113.128/25": "private",
				},
			},
			Output: Configuration{
				Communities: helpers.MustNewSubnetMap(map[string]string{
					"::/0":                     "private",
					"::ffff:203.0.113.0/121":   "public",
					"::ffff:203.0.113.128/121": "private",
				}),
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			var got Configuration
			decoder, err := mapstructure.NewDecoder(helpers.GetMapStructureDecoderConfig(&got))
			if err != nil {
				t.Fatalf("NewDecoder() error:\n%+v", err)
			}
			err = decoder.Decode(tc.Input)
			if err != nil {
				t.Fatalf("Decode() error:\n%+v", err)
			} else if diff := helpers.Diff(got, tc.Output); diff != "" {
				t.Fatalf("Decode() (-got, +want):\n%s", diff)
			}
		})
	}
}
