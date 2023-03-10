package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScryptConfig_EncodeDecode(t *testing.T) {
	encodedConfig := ScryptLevel[0].Encode()
	config, err := DecodeScryptConfig(encodedConfig)
	assert.NoError(t, err)

	assert.Equal(t, ScryptLevel[0].N, config.N)
	assert.Equal(t, ScryptLevel[0].R, config.R)
	assert.Equal(t, ScryptLevel[0].P, config.P)
}
