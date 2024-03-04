package secret

import (
	"encoding/json"
	"fmt"
	"testing"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
)

func TestSecret(t *testing.T) {
	s := String("secret")
	require.Equal(t, "secret", s.Value(), "secret")
	require.Equal(t, "REDACTED", fmt.Sprintf("%v", s))
	require.Equal(t, "REDACTED", s.String())
	require.Equal(t, "REDACTED", s.GoString())

	b, err := json.Marshal(s)
	require.NoError(t, err)
	require.Equal(t, `"REDACTED"`, string(b))
}

func TestSecret_TOMLMarshal(t *testing.T) {
	toMarshal := map[string]String{
		"some_key": String("some_secret_value"),
	}
	data, err := toml.Marshal(toMarshal)
	require.NoError(t, err)
	// If `toml` changes its default marshal formatting it is fine to update this
	// test to match - we only care that the secret is rendered in plain text.
	require.Equal(t, "some_key = 'some_secret_value'\n", string(data))
}
