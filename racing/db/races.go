package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3" // Imported for side effects
	"github.com/pkg/errors"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) *RacesRepo {
	return &RacesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *RacesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

// List will return a collection of races.
func (r *RacesRepo) List(req *racing.ListRacesRequest) ([]*racing.Race, error) {
	query := getRaceQueries()[racesList]

	query, args := r.applyFilter(query, req.Filter)

	query, err := r.applyOrderBy(query, req.OrderBy)
	if err != nil {
		return nil, errors.Wrap(err, "order by")
	}

	log.Printf("[DBG] Races query(%v): %s", args, query)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return scanRaces(rows)
}

func (r *RacesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	if filter.VisibileOnly {
		clauses = append(clauses, "visible = 1")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

// applyOrderBy applies an ordering by single field.
// orderBy is expected to be in the format "field [ASC|DESC]", e.g "advertised_start_time DESC" or "meeting_id"
func (r *RacesRepo) applyOrderBy(query string, orderBy string) (string, error) {
	var (
		field     string
		direction string
		fieldName = "orderBy"
	)

	if orderBy == "" {
		return query, nil
	}

	fieldAndOrderSplit := strings.Split(orderBy, " ")

	if len(fieldAndOrderSplit) == 2 {
		field = fieldAndOrderSplit[0]
		direction = fieldAndOrderSplit[1]
	} else if len(fieldAndOrderSplit) == 1 {
		field = fieldAndOrderSplit[0]
	} else {
		return "", &invalidArgumentError{field: fieldName, details: "orderBy is invalid, must be in the format \"field [ASC|DESC]\"."}
	}

	if field == "" {
		return "", &invalidArgumentError{field: fieldName, details: "orderBy field is required."}
	}

	switch strings.ToLower(field) {
	case "id":
		field = "id"
	case "meeting_id":
		field = "meeting_id"
	case "name":
		field = "name"
	case "number":
		field = "number"
	case "visible":
		field = "visible"
	case "advertised_start_time":
		field = "advertised_start_time"
	default:
		return "", &invalidArgumentError{field: fieldName, details: "orderBy field is invalid, must be either \"id\", \"meeting_id\", \"name\", \"number\", \"visible\" or \"advertised_start_time\"."}
	}

	query = fmt.Sprintf(" %s ORDER BY %s", query, field)

	if direction != "" {
		switch strings.ToUpper(direction) {
		case "ASC":
			query = fmt.Sprintf(" %s ASC", query)
		case "DESC":
			query = fmt.Sprintf(" %s DESC", query)
		default:
			return "", &invalidArgumentError{field: fieldName, details: "orderBy direction invalid, must be either \"ASC\" or \"DESC\"."}
		}
	}

	return query, nil
}

func scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &race.Status); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		races = append(races, &race)
	}

	return races, nil
}
