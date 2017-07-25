package gocdprovider

import (
	"testing"
)

func TestProvider(t *testing.T) {
	provider := SchemaProvider()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
