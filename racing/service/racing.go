package service

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo will be used as repository access to races.
type RacesRepo interface {
	// List should return a list of races.
	List(req *racing.ListRacesRequest) ([]*racing.Race, error)
}

// RacingService implements the Racing interface.
type RacingService struct {
	racesRepo RacesRepo
}

// NewRacingService instantiates and returns a new RacingService.
func NewRacingService(racesRepo RacesRepo) *RacingService {
	return &RacingService{racesRepo}
}

func (s *RacingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in)
	if err != nil {
		// Does this error hint at a grpc code?
		if codeErr := (interface {
			Code() codes.Code
			Error() string
		})(nil); errors.As(err, &codeErr) {
			log.Printf("[ERR] ListRaces: %v", err)

			return nil, codeErrorToStatusError(codeErr.Code(), codeErr)
		}

		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}

// codeErrorToStatusError creates a grpc status error from an existing error.
func codeErrorToStatusError(code codes.Code, err error) error {
	var details string

	// Does the error provide details?
	if de, ok := err.(interface {
		Details() string
	}); ok {
		details = de.Details()
	} else {
		details = err.Error()
	}

	switch code {
	case codes.InvalidArgument:
		return invalidArgumentToStatusError(err, details)
	default:
		return status.Error(code, details)
	}
}

// codeErrorToStatusError creates an invalid argument grpc status error from an existing error.
func invalidArgumentToStatusError(err error, details string) error {
	var (
		field string
		msg   string
	)

	// Does the error specify a field?
	if fe, ok := err.(interface {
		Field() string
	}); ok {
		field = fe.Field()
	}

	if field != "" {
		msg = fmt.Sprintf("Field \"%s\" is invalid.", field)
	} else {
		msg = "The request is invalid."
	}

	s := status.New(codes.InvalidArgument, msg)

	// Add a field violation to the status
	if field != "" || details != "" {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: details,
		}

		br := &errdetails.BadRequest{}
		br.FieldViolations = append(br.FieldViolations, v)

		withDetails, err := s.WithDetails(br)
		if err != nil {
			log.Printf("[ERR] Unexpected error when attaching details to InvalidArgument status: %v", err)
		} else {
			s = withDetails
		}
	}

	return s.Err()
}
