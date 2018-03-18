package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/kwiesmueller/ubisoft-backend-interview/pkg/feedback"

	// using postgresql in the implementation
	_ "github.com/lib/pq"
	"github.com/playnet-public/libs/log"
	"go.uber.org/zap"
)

// Connection implementing the feedback.Repository interface
type Connection struct {
	*log.Logger
	*sql.DB
}

// New database connection taking a sql connect string
func New(log *log.Logger) *Connection {
	log = log.WithFields(zap.String("component", "database"))
	return &Connection{
		Logger: log,
	}
}

// Open the db connection
func (c *Connection) Open(con string) error {
	c.Info("connecting db")
	db, err := sql.Open("postgres", con)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		c.Error("open connection error", zap.Error(err))
		return err
	}
	c.DB = db
	return nil
}

// Add feedback entry to DB
func (c *Connection) Add(entry feedback.Entry) (err error) {
	c.Debug("adding entry",
		zap.String("session", entry.SessionID),
		zap.String("user", entry.UserID),
	)
	defer func() {
		err = handleError(err)
	}()

	query := "INSERT INTO entries(session_id, user_id, rating, comment) VALUES ($1, $2, $3, $4)"
	statement, err := c.Prepare(query)
	if err != nil {
		c.Error("statement error",
			zap.String("session", entry.SessionID),
			zap.String("user", entry.UserID),
			zap.Error(err),
		)
		return err
	}

	_, err = statement.Exec(entry.SessionID, entry.UserID, entry.Rating, entry.Comment)
	if err != nil {
		c.Error("exec error",
			zap.String("session", entry.SessionID),
			zap.String("user", entry.UserID),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// GetLatest n entries from the database
func (c *Connection) GetLatest(n uint) (entries []feedback.Entry, err error) {
	c.Debug("reading entries",
		zap.Uint("limit", n),
	)
	defer c.Debug("finished reading entries",
		zap.Uint("limit", n),
		zap.Int("entries", len(entries)),
	)

	query := `SELECT id, session_id, user_id, rating, comment FROM entries
	ORDER BY id DESC LIMIT $1`
	entries, err = c.getEntries(query, n)
	if err != nil {
		c.Error("get entries failed",
			zap.Uint("limit", n),
			zap.Error(err),
		)
	}

	return entries, err
}

// GetLatestFiltered n entries by rating from the database
func (c *Connection) GetLatestFiltered(n uint, filter int) (entries []feedback.Entry, err error) {
	c.Debug("reading entries",
		zap.Uint("limit", n),
		zap.Int("filter", filter),
	)
	defer c.Debug("finished reading entries",
		zap.Uint("limit", n),
		zap.Int("filter", filter),
		zap.Int("entries", len(entries)),
	)

	query := `SELECT id, session_id, user_id, rating, comment FROM entries WHERE rating = $2
	ORDER BY id DESC LIMIT $1`
	entries, err = c.getEntries(query, n, filter)
	if err != nil {
		c.Error("get entries failed",
			zap.Uint("limit", n),
			zap.Int("filter", filter),
			zap.Error(err),
		)
	}

	return entries, err
}

func (c *Connection) getEntries(query string, args ...interface{}) ([]feedback.Entry, error) {
	statement, err := c.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "statement error")
	}
	rows, err := statement.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entry := feedback.Entry{}
	var entries []feedback.Entry
	var count int8

	for rows.Next() {
		err := rows.Scan(
			&entry.ID,
			&entry.SessionID,
			&entry.UserID,
			&entry.Rating,
			&entry.Comment,
		)
		if err != nil {
			return nil, errors.Wrap(err, "row scan error")
		}
		entries = append(entries, entry)
		count++
	}
	fmt.Println(entries)
	return entries, nil
}

func handleError(err error) error {
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return feedback.ErrDuplicateEntry
		}
	}
	return err
}
