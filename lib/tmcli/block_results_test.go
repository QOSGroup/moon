// Copyright 2018 The QOS Authors

package tmcli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockResult(t *testing.T) {
	opt, err := NewOption()
	assert.Nil(t, err)

	res, err := NewClient(opt).BlockResults.Retrieve(nil, &BlockResultsOption{Height: 1})
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
