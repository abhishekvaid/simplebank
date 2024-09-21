package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// Generic function to match the body of any object
func MatchResponseBodyWith[T any](t *testing.T, buf *bytes.Buffer, want T) {
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

func CompareResponseBodyJSON(t *testing.T, rr *httptest.ResponseRecorder, expected interface{}) {
	// Marshal the expected object into JSON bytes
	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	// Unmarshal the expected JSON into a map or slice (generic type)
	var expectedMap interface{}
	err = json.Unmarshal(expectedJSON, &expectedMap)
	require.NoError(t, err)

	// Unmarshal the response body into a map or slice (generic type)
	var actualMap interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &actualMap)
	require.NoError(t, err)

	// Compare the actual and expected maps/slices
	require.Equal(t, expectedMap, actualMap)
}
