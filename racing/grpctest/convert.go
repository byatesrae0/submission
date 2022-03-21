package grpctest

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TimeToTimestampPB converts a time to a pb timestamp
func TimeToTimestampPB(t *testing.T, tt time.Time) *timestamppb.Timestamp {
	t.Helper()

	ts, err := ptypes.TimestampProto(tt)
	require.NoError(t, err, "TimeToTimestampProto")

	return ts
}
