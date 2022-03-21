package service

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.neds.sh/matty/entain/sports/grpctest"
	"git.neds.sh/matty/entain/sports/proto/sports"
)

func TestSportsServiceListEvents(t *testing.T) {
	for _, tc := range []struct {
		name        string
		with        *SportsService
		giveContext context.Context
		giveRequest *sports.ListEventsRequest
		expect      *sports.ListEventsResponse
	}{
		{
			name:        "success_no_results",
			with:        NewSportsService(),
			giveContext: context.Background(),
			giveRequest: &sports.ListEventsRequest{},
			expect: &sports.ListEventsResponse{
				Events: []*sports.Event{
					{
						Id:                  1,
						Name:                "Event 1",
						AdvertisedStartTime: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					},
					{
						Id:                  2,
						Name:                "Event 2",
						AdvertisedStartTime: timestamppb.New(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)),
					},
					{
						Id:                  3,
						Name:                "Event 3",
						AdvertisedStartTime: timestamppb.New(time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC)),
					},
				},
			},
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(tc.giveContext, time.Second*2)
			t.Cleanup(cancel)

			actual, actualErr := tc.with.ListEvents(ctx, tc.giveRequest)

			if tc.expect != nil {
				assert.Empty(t, cmp.Diff(tc.expect, actual, cmp.Options{protocmp.Transform(), protocmp.IgnoreUnknown()}), "expected vs actual")
			} else {
				assert.Nil(t, actual, "actual")
			}

			assert.NoError(t, actualErr, "actualErr")
		})
	}
}

func TestSportsServiceGetEvent(t *testing.T) {
	for _, tc := range []struct {
		name        string
		with        *SportsService
		giveContext context.Context
		giveRequest *sports.GetEventRequest
		expect      *sports.Event
		errAssert   func(*testing.T, error) bool
	}{
		{
			name:        "success",
			with:        NewSportsService(),
			giveContext: context.Background(),
			giveRequest: &sports.GetEventRequest{Id: 1},
			expect: &sports.Event{
				Id:                  1,
				Name:                "Event 1",
				AdvertisedStartTime: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name:        "repo_err_not_found",
			with:        NewSportsService(),
			giveContext: context.Background(),
			giveRequest: &sports.GetEventRequest{Id: 123},
			errAssert: grpctest.NewGRPCErrorAsserter(
				codes.NotFound,
				"Event with ID 123 does not exist.",
				nil,
			),
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(tc.giveContext, time.Second*2)
			t.Cleanup(cancel)

			actual, actualErr := tc.with.GetEvent(ctx, tc.giveRequest)

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
