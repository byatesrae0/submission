package db

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"git.neds.sh/matty/entain/racing/grpctest"
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
		give        *racing.ListRacesRequest
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
			give: &racing.ListRacesRequest{Filter: &racing.ListRacesRequestFilter{}},
			expect: []*racing.Race{
				{
					Id:                  1,
					MeetingId:           2,
					Name:                "3",
					Number:              4,
					Visible:             true,
					AdvertisedStartTime: grpctest.TimeToTimestampPB(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
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
			give: &racing.ListRacesRequest{},
			expect: []*racing.Race{
				{
					Id:                  1,
					MeetingId:           2,
					Name:                "3",
					Number:              4,
					Visible:             true,
					AdvertisedStartTime: grpctest.TimeToTimestampPB(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
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
			give: &racing.ListRacesRequest{},
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
			give: &racing.ListRacesRequest{
				Filter: &racing.ListRacesRequestFilter{
					MeetingIds:   []int64{1},
					VisibileOnly: true,
				},
			},
			expect: []*racing.Race{
				{
					Id:                  1,
					MeetingId:           2,
					Name:                "3",
					Number:              4,
					Visible:             true,
					AdvertisedStartTime: grpctest.TimeToTimestampPB(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:                  5,
					MeetingId:           6,
					Name:                "7",
					Number:              8,
					Visible:             false,
					AdvertisedStartTime: grpctest.TimeToTimestampPB(t, time.Date(2001, time.February, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			name: "success_ordering",
			with: func() *RacesRepo {
				db, mock := newSQLMock(t)

				mock.ExpectQuery(getRaceQueries()[racesList]).
					WillReturnRows(
						mock.NewRows(listColumns),
					)

				return NewRacesRepo(db)
			}(),
			give: &racing.ListRacesRequest{OrderBy: "id ASC"},
		},
		{
			name: "success_ordering_no_direction",
			with: func() *RacesRepo {
				db, mock := newSQLMock(t)

				mock.ExpectQuery(getRaceQueries()[racesList]).
					WillReturnRows(
						mock.NewRows(listColumns),
					)

				return NewRacesRepo(db)
			}(),
			give: &racing.ListRacesRequest{OrderBy: "meeting_id"},
		},
		{
			name: "db_err",
			with: func() *RacesRepo {
				db, mock := newSQLMock(t)

				mock.ExpectQuery(getRaceQueries()[racesList]).WillReturnError(errors.New("TestError123"))

				return NewRacesRepo(db)
			}(),
			give:        &racing.ListRacesRequest{},
			expectError: "TestError123",
		},
		{
			name: "invalid_orderby_field",
			with: func() *RacesRepo {
				db, _ := newSQLMock(t)

				return NewRacesRepo(db)
			}(),
			give:        &racing.ListRacesRequest{OrderBy: "meeting_iid"},
			expectError: "order by: invalid argument \"orderBy\", orderBy field is invalid, must be either \"id\", \"meeting_id\", \"name\", \"number\", \"visible\" or \"advertised_start_time\".",
		},
		{
			name: "invalid_orderby_no_field",
			with: func() *RacesRepo {
				db, _ := newSQLMock(t)

				return NewRacesRepo(db)
			}(),
			give:        &racing.ListRacesRequest{OrderBy: " DESC"},
			expectError: "order by: invalid argument \"orderBy\", orderBy field is required.",
		},
		{
			name: "invalid_orderby_format",
			with: func() *RacesRepo {
				db, _ := newSQLMock(t)

				return NewRacesRepo(db)
			}(),
			give:        &racing.ListRacesRequest{OrderBy: "A B C"},
			expectError: "order by: invalid argument \"orderBy\", orderBy is invalid, must be in the format \"field [ASC|DESC]\".",
		},
		{
			name: "invalid_orderby_direction",
			with: func() *RacesRepo {
				db, _ := newSQLMock(t)

				return NewRacesRepo(db)
			}(),
			give:        &racing.ListRacesRequest{OrderBy: "advertised_start_time DESsC"},
			expectError: "order by: invalid argument \"orderBy\", orderBy direction invalid, must be either \"ASC\" or \"DESC\".",
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

func newSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "sqlmock.New()")

	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	return db, mock
}
