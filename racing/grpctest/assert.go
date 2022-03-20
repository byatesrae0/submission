package grpctest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

// NewGRPCErrorAsserter creates a new asserter that asserts various parts of a gRPC
// error.
func NewGRPCErrorAsserter(
	expectedErrCode codes.Code,
	expectedErrMsg string,
	expectedErrDetails proto.Message,
) func(*testing.T, error) bool {
	return func(t *testing.T, err error) bool {
		t.Helper()

		if expectedErrCode == codes.OK {
			return assert.NoError(t, err)
		}

		st, ok := status.FromError(err)
		require.True(t, ok, "is status error")

		result := assert.Equal(t, expectedErrCode.String(), st.Code().String(), "status error code")

		result = assert.Equal(t, expectedErrMsg, st.Message(), "status error message") && result

		if expectedErrDetails == nil {
			return assert.Len(t, st.Details(), 0, "status error details") && result
		}

		if ok := assert.Len(t, st.Details(), 1, "status error details"); !ok {
			return false
		}

		return assert.Empty(t, cmp.Diff(expectedErrDetails, st.Details()[0], cmp.Options{protocmp.Transform(), protocmp.IgnoreUnknown()}), "status error details diff") && result
	}
}
