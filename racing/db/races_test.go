package db

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

func TestRacesRepoList(t *testing.T) {
	t.Parallel()

	listColumns := []string{
		"id",
		"meeting_id",
		"name",
		"number",
		"visible",
		"advertised_start_time",
	}

	for _, tc := range []struct {
		name        string
		with        *RacesRepo
		give        *racing.ListRacesRequestFilter
		expect      []*racing.Race
		expectError string
	}{
		{
			name: "success",
			with: func() *RacesRepo {
				db, mock := newSQLMock(t)

				mock.ExpectQuery(getRaceQueries()[racesList]).
					WillReturnRows(
						mock.NewRows(listColumns).AddRow(1, 2, "3", 4, true, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
					)

				return NewRacesRepo(db)
			}(),
			give: &racing.ListRacesRequestFilter{},
			expect: []*racing.Race{
				{
					Id:                  1,
					MeetingId:           2,
					Name:                "3",
					Number:              4,
					Visible:             true,
					AdvertisedStartTime: timeToTimestampPB(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "success_no_filter",
			with: func() *RacesRepo {
				db, mock := newSQLMock(t)

				mock.ExpectQuery(getRaceQueries()[racesList]).
					WillReturnRows(
						mock.NewRows(listColumns).AddRow(1, 2, "3", 4, true, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
					)

				return NewRacesRepo(db)
			}(),
			give: nil,
			expect: []*racing.Race{
				{
					Id:                  1,
					MeetingId:           2,
					Name:                "3",
					Number:              4,
					Visible:             true,
					AdvertisedStartTime: timeToTimestampPB(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "success_no_results",
			with: func() *RacesRepo {
				db, mock := newSQLMock(t)

				mock.ExpectQuery(getRaceQueries()[racesList]).WillReturnRows(mock.NewRows(listColumns))

				return NewRacesRepo(db)
			}(),
			give: &racing.ListRacesRequestFilter{},
		},
		{
			name: "success_multiple_results",
			with: func() *RacesRepo {
				db, mock := newSQLMock(t)

				mock.ExpectQuery(getRaceQueries()[racesList]).
					WillReturnRows(
						mock.NewRows(listColumns).
							AddRow(1, 2, "3", 4, true, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)).
							AddRow(5, 6, "7", 8, false, time.Date(2001, time.February, 2, 0, 0, 0, 0, time.UTC)),
					)

				return NewRacesRepo(db)
			}(),
			give: &racing.ListRacesRequestFilter{
				MeetingIds: []int64{1},
			},
			expect: []*racing.Race{
				{
					Id:                  1,
					MeetingId:           2,
					Name:                "3",
					Number:              4,
					Visible:             true,
					AdvertisedStartTime: timeToTimestampPB(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:                  5,
					MeetingId:           6,
					Name:                "7",
					Number:              8,
					Visible:             false,
					AdvertisedStartTime: timeToTimestampPB(t, time.Date(2001, time.February, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "db_err",
			with: func() *RacesRepo {
				db, mock := newSQLMock(t)

				mock.ExpectQuery(getRaceQueries()[racesList]).WillReturnError(errors.New("TestError123"))

				return NewRacesRepo(db)
			}(),
			give:        &racing.ListRacesRequestFilter{},
			expectError: "TestError123",
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, actualErr := tc.with.List(tc.give)

			if tc.expect != nil {
				assert.Empty(t, cmp.Diff(tc.expect, actual, cmp.Options{protocmp.Transform(), protocmp.IgnoreUnknown()}), "expected vs actual")
			} else {
				assert.Nil(t, actual, "actual")
			}

			if tc.expectError != "" {
				assert.EqualError(t, actualErr, tc.expectError, "actualErr")
			} else {
				assert.NoError(t, actualErr, "actualErr")
			}
		})
	}
}

func timeToTimestampPB(t *testing.T, tt time.Time) *timestamppb.Timestamp {
	ts, err := ptypes.TimestampProto(tt)
	require.NoError(t, err, "TimeToTimestampProto")

	return ts
}

func newSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "sqlmock.New()")

	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	return db, mock
}
