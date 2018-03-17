package database

import (
	"database/sql"

	"github.com/lib/pq"

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
func (c *Connection) GetLatest(n uint) ([]feedback.Entry, error) {
	c.Debug("reading entries",
		zap.Uint("count", n),
	)

	query := `SELECT id, session_id, user_id, rating, comment FROM entries
	ORDER BY id DESC LIMIT $1`
	statement, err := c.Prepare(query)
	if err != nil {
		c.Error("statement error",
			zap.Uint("count", n),
			zap.Error(err),
		)
		return nil, err
	}

	rows, err := statement.Query(n)
	if err != nil {
		c.Error("query error",
			zap.Uint("count", n),
			zap.Error(err),
		)
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
			c.Error("row scan error",
				zap.Uint("count", n),
				zap.Error(err),
			)
			return nil, err
		}
		entries = append(entries, entry)
		count++
	}
	c.Debug("finished reading entries",
		zap.Uint("count", n),
		zap.Int8("entries", count),
	)

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
