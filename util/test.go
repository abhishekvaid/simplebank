package util

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// Generic function to match the body of any object
func RequireBodyMatch[T any](t *testing.T, buf *bytes.Buffer, want T) {
	t.Helper()

	// Read from the buffer
	bytes, err := io.ReadAll(buf)
	require.NoError(t, err)

	// Unmarshal into the type T
	var got T
	err = json.Unmarshal(bytes, &got)
	require.NoError(t, err)

	// Assert equality between the unmarshalled object and the expected one
	require.Equal(t, want, got)
}
