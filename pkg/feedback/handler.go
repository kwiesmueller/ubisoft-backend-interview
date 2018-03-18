package feedback

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Handler for the service endpoints
func (s *Service) Handler() *mux.Router {
	m := mux.NewRouter()
	m.Path("/").Methods("GET").HandlerFunc(s.MakeHandler(s.getEntries))
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

func (s *Service) getEntries(w http.ResponseWriter, r *http.Request) error {
	entries, err := s.GetLatest(15)
	if err != nil {
		if err := writeError(w, err); err != nil {
			s.Error("write error", zap.Error(err))
		}
		return err
	}
	return writeJSON(w, entries)
}

func (s *Service) addEntry(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if err != nil {
			s.Warn("failed adding entry", zap.Error(err))
			if err := writeError(w, err); err != nil {
				s.Error("write error", zap.Error(err))
			}
		}
	}()
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
	fmt.Println("userid", entry.UserID)

	if len(entry.UserID) < 1 {
		return ErrNoUserID
	}

	if err := s.Add(entry); err != nil {
		return err
	}

	return writeJSON(w, struct{}{})
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
