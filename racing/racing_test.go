package main

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/testing/protocmp"

	"git.neds.sh/matty/entain/racing/grpctest"
	"git.neds.sh/matty/entain/racing/proto/racing"
)

func TestRacingListRaces(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Parallel()

	for _, tc := range []struct {
		name        string
		giveContext context.Context
		giveRequest *racing.ListRacesRequest
		expect      *racing.ListRacesResponse
		errAssert   func(*testing.T, error) bool
	}{
		{
			name:        "success",
			giveContext: context.Background(),
			giveRequest: &racing.ListRacesRequest{
				OrderBy: "number desc",
				Filter: &racing.ListRacesRequestFilter{
					MeetingIds:   []int64{9},
					VisibileOnly: true,
				},
			},
			// WARNING: This test relies on data seeded in the db package.
			expect: &racing.ListRacesResponse{
				Races: []*racing.Race{
					{Id: 69, MeetingId: 9, Name: "Michigan buffalo", Number: 10, Visible: true, AdvertisedStartTime: grpctest.TimeToTimestampPB(t, time.Unix(1614558171, 0))},
					{Id: 37, MeetingId: 9, Name: "Mississippi ducks", Number: 8, Visible: true, AdvertisedStartTime: grpctest.TimeToTimestampPB(t, time.Unix(1614720480, 0))},
					{Id: 6, MeetingId: 9, Name: "Nebraska giants", Number: 4, Visible: true, AdvertisedStartTime: grpctest.TimeToTimestampPB(t, time.Unix(1614534492, 0))},
				},
			},
		},
		{
			name:        "invalid_orderby",
			giveContext: context.Background(),
			giveRequest: &racing.ListRacesRequest{
				OrderBy: "number des",
				Filter: &racing.ListRacesRequestFilter{
					MeetingIds:   []int64{9},
					VisibileOnly: true,
				},
			},
			errAssert: grpctest.NewGRPCErrorAsserter(
				codes.InvalidArgument,
				"Field \"orderBy\" is invalid.",
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequest_FieldViolation{
						{
							Field:       "orderBy",
							Description: "orderBy direction invalid, must be either \"ASC\" or \"DESC\".",
						},
					},
				},
			),
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(tc.giveContext, time.Second*3)
			t.Cleanup(cancel)

			actual, actualErr := racingCli.ListRaces(ctx, tc.giveRequest)

			if tc.expect != nil {
				assert.Empty(t, cmp.Diff(tc.expect, actual, cmp.Options{protocmp.Transform(), protocmp.IgnoreUnknown()}), "expected vs actual")
			} else {
				assert.Nil(t, actual, "actual")
			}

			if tc.errAssert != nil {
				tc.errAssert(t, actualErr)
			} else {
				assert.NoError(t, actualErr, "actualErr")
			}
		})
	}
}
