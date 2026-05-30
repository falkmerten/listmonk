package models

import (
	"encoding/json"
	"testing"

	"pgregory.net/rapid"
)

// Feature: aws-ses-integration, Property 1: Credentials round-trip
// Validates: Requirements 1.1
func TestSESSettingsRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		original := Settings{}
		original.SES.Enabled = rapid.Bool().Draw(t, "enabled")
		original.SES.Region = rapid.StringN(1, 32, -1).Draw(t, "region")
		original.SES.AccessKeyID = rapid.StringN(1, 64, -1).Draw(t, "access_key_id")
		original.SES.SecretAccessKey = rapid.StringN(1, 64, -1).Draw(t, "secret_access_key")
		original.SES.RoleARN = rapid.StringN(0, 128, -1).Draw(t, "role_arn")
		original.SES.ConfigSet = rapid.StringN(0, 64, -1).Draw(t, "config_set")
		original.SES.MaxMsgRetries = rapid.IntRange(0, 10).Draw(t, "max_msg_retries")
		original.SES.PricePerMessage = float64(rapid.IntRange(1, 1000).Draw(t, "price_cents")) / 10000.0

		data, err := json.Marshal(original.SES)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}

		var restored struct {
			Enabled         bool    `json:"enabled"`
			Region          string  `json:"region"`
			AccessKeyID     string  `json:"access_key_id"`
			SecretAccessKey string  `json:"secret_access_key,omitempty"`
			RoleARN         string  `json:"role_arn"`
			ConfigSet       string  `json:"config_set"`
			MaxMsgRetries   int     `json:"max_msg_retries"`
			PricePerMessage float64 `json:"price_per_message"`
		}
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}

		if restored.Enabled != original.SES.Enabled {
			t.Errorf("Enabled mismatch: got %v, want %v", restored.Enabled, original.SES.Enabled)
		}
		if restored.Region != original.SES.Region {
			t.Errorf("Region mismatch: got %q, want %q", restored.Region, original.SES.Region)
		}
		if restored.AccessKeyID != original.SES.AccessKeyID {
			t.Errorf("AccessKeyID mismatch: got %q, want %q", restored.AccessKeyID, original.SES.AccessKeyID)
		}
		if restored.SecretAccessKey != original.SES.SecretAccessKey {
			t.Errorf("SecretAccessKey mismatch: got %q, want %q", restored.SecretAccessKey, original.SES.SecretAccessKey)
		}
		if restored.RoleARN != original.SES.RoleARN {
			t.Errorf("RoleARN mismatch: got %q, want %q", restored.RoleARN, original.SES.RoleARN)
		}
		if restored.ConfigSet != original.SES.ConfigSet {
			t.Errorf("ConfigSet mismatch: got %q, want %q", restored.ConfigSet, original.SES.ConfigSet)
		}
		if restored.MaxMsgRetries != original.SES.MaxMsgRetries {
			t.Errorf("MaxMsgRetries mismatch: got %d, want %d", restored.MaxMsgRetries, original.SES.MaxMsgRetries)
		}
		if restored.PricePerMessage != original.SES.PricePerMessage {
			t.Errorf("PricePerMessage mismatch: got %f, want %f", restored.PricePerMessage, original.SES.PricePerMessage)
		}
	})
}
