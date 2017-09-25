package gocd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testMaterialEquality(t *testing.T) {
	s1 := Material{
		Type: "git",
		Attributes: MaterialAttributes{
			URL: "https://github.com/gocd/gocd",
		},
	}

	s2 := Material{
		Type: "git",
		Attributes: MaterialAttributes{
			Name: "gocd-src",
			URL:  "https://github.com/gocd/gocd",
		},
	}

	assert.True(t, s1.Equal(&s2))

}
