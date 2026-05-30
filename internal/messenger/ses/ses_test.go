package ses

import (
	"log"
	"os"
	"testing"
)

func TestNew_EmptyRegionReturnsError(t *testing.T) {
	_, err := New("ses", Options{
		Region:          "",
		AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}, log.New(os.Stderr, "", 0))

	if err == nil {
		t.Fatal("expected error when region is empty, got nil")
	}
}

func TestNew_ValidOptionsSucceeds(t *testing.T) {
	m, err := New("ses", Options{
		Region:          "us-east-1",
		AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}, log.New(os.Stderr, "", 0))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil SESMessenger")
	}
	if m.Name() != "ses" {
		t.Fatalf("expected name 'ses', got %q", m.Name())
	}
}

func TestNew_DefaultPricePerMessage(t *testing.T) {
	m, err := New("ses", Options{
		Region:          "us-east-1",
		AccessKeyID:     "key",
		SecretAccessKey: "secret",
	}, log.New(os.Stderr, "", 0))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.opt.PricePerMessage != 0.0001 {
		t.Fatalf("expected default PricePerMessage 0.0001, got %v", m.opt.PricePerMessage)
	}
}
