package service

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/testing/protocmp"

	"git.neds.sh/matty/entain/racing/grpctest"
	"git.neds.sh/matty/entain/racing/proto/racing"
)

func TestRacingServiceListRaces(t *testing.T) {
	for _, tc := range []struct {
		name        string
		with        *RacingService
		giveContext context.Context
		giveRequest *racing.ListRacesRequest
		expect      *racing.ListRacesResponse
		errAssert   func(*testing.T, error) bool
	}{
		{
			name: "success_no_results",
			with: NewRacingService(&mockRacesRepo{
				list: func(filter *racing.ListRacesRequest) ([]*racing.Race, error) {
					return nil, nil
				},
			}),
			giveContext: context.Background(),
			giveRequest: &racing.ListRacesRequest{},
			expect:      &racing.ListRacesResponse{},
		},
		{
			name: "repo_err_invalid_argument",
			with: NewRacingService(&mockRacesRepo{
				list: func(filter *racing.ListRacesRequest) ([]*racing.Race, error) {
					return nil, errors.Wrap(&mockCodeError{code: codes.InvalidArgument, field: "TestField", details: "TestDetails"}, "wrapped")
				},
			}),
			giveContext: context.Background(),
			giveRequest: &racing.ListRacesRequest{},
			errAssert: grpctest.NewGRPCErrorAsserter(
				codes.InvalidArgument,
				"Field \"TestField\" is invalid.",
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequest_FieldViolation{
						{
							Field:       "TestField",
							Description: "TestDetails",
						},
					},
				},
			),
		},
		{
			name: "repo_err_invalid_argument_no_field",
			with: NewRacingService(&mockRacesRepo{
				list: func(filter *racing.ListRacesRequest) ([]*racing.Race, error) {
					return nil, errors.Wrap(&mockCodeError{code: codes.InvalidArgument, field: "", details: "TestDetails"}, "wrapped")
				},
			}),
			giveContext: context.Background(),
			giveRequest: &racing.ListRacesRequest{},
			errAssert: grpctest.NewGRPCErrorAsserter(
				codes.InvalidArgument,
				"The request is invalid.",
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequest_FieldViolation{
						{
							Description: "TestDetails",
						},
					},
				},
			),
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(tc.giveContext, time.Second*2)
			t.Cleanup(cancel)

			actual, actualErr := tc.with.ListRaces(ctx, tc.giveRequest)

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

func TestRacingServiceGetRace(t *testing.T) {
	for _, tc := range []struct {
		name        string
		with        *RacingService
		giveContext context.Context
		giveRequest *racing.GetRaceRequest
		expect      *racing.Race
		errAssert   func(*testing.T, error) bool
	}{
		{
			name: "success",
			with: NewRacingService(&mockRacesRepo{
				get: func(_ int64) (*racing.Race, error) {
					return &racing.Race{}, nil
				},
			}),
			giveContext: context.Background(),
			giveRequest: &racing.GetRaceRequest{},
			expect:      &racing.Race{},
		},
		{
			name: "repo_err_not_found",
			with: NewRacingService(&mockRacesRepo{
				get: func(_ int64) (*racing.Race, error) {
					return nil, errors.Wrap(&mockCodeError{code: codes.NotFound, details: "Test error message."}, "wrapped")
				},
			}),
			giveContext: context.Background(),
			giveRequest: &racing.GetRaceRequest{},
			errAssert: grpctest.NewGRPCErrorAsserter(
				codes.NotFound,
				"Test error message.",
				nil,
			),
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(tc.giveContext, time.Second*2)
			t.Cleanup(cancel)

			actual, actualErr := tc.with.GetRace(ctx, tc.giveRequest)

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

type mockRacesRepo struct {
	list func(*racing.ListRacesRequest) ([]*racing.Race, error)
	get  func(int64) (*racing.Race, error)
}

func (r *mockRacesRepo) List(req *racing.ListRacesRequest) ([]*racing.Race, error) {
	if r.list != nil {
		return r.list(req)
	}

	panic("mockRacesRepo: unexpected call to List().")
}

func (r *mockRacesRepo) Get(id int64) (*racing.Race, error) {
	if r.get != nil {
		return r.get(id)
	}

	panic("mockRacesRepo: unexpected call to Get().")
}

type mockCodeError struct {
	code      codes.Code
	errString string
	field     string
	details   string
}

func (e *mockCodeError) Code() codes.Code {
	return e.code
}

func (e *mockCodeError) Field() string {
	return e.field
}

func (e *mockCodeError) Details() string {
	return e.details
}

func (e *mockCodeError) Error() string {
	return e.errString
}
