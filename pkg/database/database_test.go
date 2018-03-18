package database

import (
	"testing"

	"github.com/lib/pq"

	"github.com/kwiesmueller/ubisoft-backend-interview/pkg/feedback"
	"github.com/playnet-public/libs/log"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestConnection_Open(t *testing.T) {
	con := New(log.NewNop())
	err := con.Open("")
	if err == nil {
		t.Fatal("connection error mishandled")
	}
}

func TestConnection_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	con := New(log.NewNop())
	con.DB = db

	tests := []struct {
		name            string
		input           feedback.Entry
		expectedPrepare *sqlmock.ExpectedPrepare
		expectedExec    *sqlmock.ExpectedExec
		err             bool
	}{
		{
			"basicEntry",
			feedback.Entry{
				SessionID: "abc123",
				UserID:    "123abc",
				Rating:    1,
				Comment:   "test",
			},
			mock.ExpectPrepare("INSERT INTO entries(.+) VALUES (.+)"),
			mock.ExpectExec("INSERT INTO entries(.+) VALUES (.+)").WithArgs(
				"abc123", "123abc", 1, "test",
			).WillReturnResult(sqlmock.NewResult(0, 0)),
			false,
		},
		{
			"errorPrepareEntry",
			feedback.Entry{

				SessionID: "err",
				UserID:    "err",
				Rating:    1,
				Comment:   "err",
			},
			nil,
			nil,
			true,
		},
		{
			"errorExecEntry",
			feedback.Entry{
				SessionID: "abc123",
				UserID:    "123abc",
				Rating:    1,
				Comment:   "test",
			},
			mock.ExpectPrepare("INSERT INTO entries(.+) VALUES (.+)"),
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err = sqlmock.New()
			if err != nil {
				t.Error(err)
			}
			mock.MatchExpectationsInOrder(false)

			err := con.Add(tt.input)
			if (err == nil) == tt.err {
				t.Fatalf("Add() == %v want %v", err, tt.err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal("expectations not met for insert", err)
			}
		})
	}
}

func TestConnection_GetLatest(t *testing.T) {
	tests := []struct {
		name            string
		input           uint
		result          []feedback.Entry
		expectedPrepare bool
		expectedQuery   bool
		err             bool
	}{
		{
			"basicEntry",
			1,
			[]feedback.Entry{
				{
					SessionID: "abc123",
					UserID:    "123abc",
					Rating:    1,
					Comment:   "test",
				},
			},
			true,
			true,
			false,
		},
		{
			"errorPrepareEntry",
			1,
			[]feedback.Entry{
				{
					SessionID: "abc123",
					UserID:    "123abc",
					Rating:    1,
					Comment:   "test",
				},
			},
			false,
			false,
			true,
		},
		{
			"errorQueryEntry",
			1,
			[]feedback.Entry{
				{
					SessionID: "abc123",
					UserID:    "123abc",
					Rating:    1,
					Comment:   "test",
				},
			},
			true,
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			con := New(log.NewNop())
			con.DB = db

			mock.MatchExpectationsInOrder(false)

			query := `SELECT id, session_id, user_id, rating, comment FROM entries
			ORDER BY id DESC LIMIT (.+)`
			rows := sqlmock.NewRows([]string{"id", "session_id", "user_id", "rating", "comment"})
			for i, e := range tt.result {
				rows = rows.AddRow(i+1, e.SessionID, e.UserID, e.Rating, e.Comment)
			}
			if tt.expectedPrepare {
				mock.ExpectPrepare(query)
			}
			if tt.expectedQuery {
				mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			}

			entries, err := con.GetLatest(tt.input)
			if (err == nil) == tt.err {
				t.Fatalf("Add() == %v want %v", err, tt.err)
			}
			if count := len(entries); count > int(tt.input) {
				t.Fatal("too many entries returned", count)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal("expectations not met for insert", err)
			}
		})
	}
}

func TestConnection_GetLatestFiltered(t *testing.T) {
	tests := []struct {
		name            string
		limit           uint
		filter          int
		result          []feedback.Entry
		expectedPrepare bool
		expectedQuery   bool
		err             bool
	}{
		{
			"basicEntry",
			1,
			1,
			[]feedback.Entry{
				{
					SessionID: "abc123",
					UserID:    "123abc",
					Rating:    1,
					Comment:   "test",
				},
			},
			true,
			true,
			false,
		},
		{
			"errorPrepareEntry",
			1,
			1,
			[]feedback.Entry{
				{
					SessionID: "abc123",
					UserID:    "123abc",
					Rating:    1,
					Comment:   "test",
				},
			},
			false,
			false,
			true,
		},
		{
			"errorQueryEntry",
			1,
			1,
			[]feedback.Entry{
				{
					SessionID: "abc123",
					UserID:    "123abc",
					Rating:    1,
					Comment:   "test",
				},
			},
			true,
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			con := New(log.NewNop())
			con.DB = db

			mock.MatchExpectationsInOrder(false)

			query := `SELECT id, session_id, user_id, rating, comment FROM entries WHERE rating = (.+)
			ORDER BY id DESC LIMIT (.+)`
			rows := sqlmock.NewRows([]string{"id", "session_id", "user_id", "rating", "comment"})
			for i, e := range tt.result {
				rows = rows.AddRow(i+1, e.SessionID, e.UserID, e.Rating, e.Comment)
			}
			if tt.expectedPrepare {
				mock.ExpectPrepare(query)
			}
			if tt.expectedQuery {
				mock.ExpectQuery(query).WithArgs(1, 2).WillReturnRows(rows)
			}

			entries, err := con.GetLatestFiltered(tt.limit, tt.filter)
			if (err == nil) == tt.err {
				t.Fatalf("Add() == %v want %v", err, tt.err)
			}
			if count := len(entries); count > int(tt.limit) {
				t.Fatal("too many entries returned", count)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal("expectations not met for insert", err)
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	errs := []struct {
		in  error
		out error
	}{
		{
			&pq.Error{
				Code: "23505",
			},
			feedback.ErrDuplicateEntry,
		},
	}
	for _, err := range errs {
		t.Run(err.in.Error(), func(t *testing.T) {
			if got := handleError(err.in); got != err.out {
				t.Fatalf("handleError() = %v want %v", got, err.out)
			}
		})
	}

}
