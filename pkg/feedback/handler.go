package feedback

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Handler for the service endpoints
func (s *Service) Handler() *mux.Router {
	m := mux.NewRouter()
	m.Path("/list").Methods("GET").HandlerFunc(s.MakeHandler(s.getEntries))
	m.Path("/{sessionID}").Methods("POST").HandlerFunc(s.MakeHandler(s.addEntry))
	return m
}

type handler func(http.ResponseWriter, *http.Request) error

// MakeHandler with logging
func (s *Service) MakeHandler(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			s.Error("request error", zap.Error(err))
		}
	}
}

func (s *Service) getEntries(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() { s.deferError(w, err) }()
	var entries []Entry

	limit := uint(15)
	limitParam := r.URL.Query().Get("limit")
	if len(limitParam) > 0 {
		u64, err := strconv.ParseUint(limitParam, 10, 32)
		if err != nil {
			return errors.Wrap(err, "invalid limit value")
		}
		limit = uint(u64)
	}

	filter := r.URL.Query().Get("filter")
	if len(filter) > 0 {
		entries, err = s.getFiltered(limit, filter)
	} else {
		entries, err = s.GetLatest(limit)
	}
	if err != nil {
		return err
	}
	if entries == nil {
		entries = []Entry{}
	}
	return writeJSON(w, entries)
}

func (s *Service) getFiltered(limit uint, filter string) (entries []Entry, err error) {
	f, err := strconv.Atoi(filter)
	if err != nil {
		return nil, errors.Wrap(err, "invalid filter value")
	}
	entries, err = s.GetLatestFiltered(limit, f)
	if err != nil {
		return nil, err
	}
	return entries, err
}

func (s *Service) addEntry(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() { s.deferError(w, err) }()

	var entry Entry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		return err
	}

	vars := mux.Vars(r)
	if len(vars["sessionID"]) < 1 {
		return ErrNoSession
	}

	entry.SessionID = vars["sessionID"]
	entry.UserID = r.Header.Get("Ubi-UserId")

	if len(entry.UserID) < 1 {
		return ErrNoUserID
	}

	if err := s.Add(entry); err != nil {
		return err
	}

	return writeJSON(w, struct{}{})
}

func (s *Service) deferError(w http.ResponseWriter, err error) {
	if err != nil {
		s.Warn("failed adding entry", zap.Error(err))
		if err := writeError(w, err); err != nil {
			s.Error("write error", zap.Error(err))
		}
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	_, err = w.Write(data)
	return nil
}

type errorResponse struct {
	Err string `json:"error"`
}

func writeError(w http.ResponseWriter, e error) error {
	data, err := json.Marshal(errorResponse{e.Error()})
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	http.Error(w, string(data), http.StatusInternalServerError)
	return nil
}
