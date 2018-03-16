package feedback

import (
	"testing"

	"github.com/playnet-public/libs/log"
)

func TestNew(t *testing.T) {
	log := log.New("nop", "", true)
	if New(log, nil) == nil {
		t.Errorf("New() == nil")
	}
}

func TestAdd(t *testing.T) {
	log := log.New("nop", "", true)
	svc := New(log, &mockRepository{})
	if svc == nil {
		t.Errorf("New() == nil")
	}

	e := Entry{Rating: 0}
	err := svc.Add(e)
	if err != errInvalidRating {
		t.Fatal("Add() should return error")
	}
	e = Entry{Rating: 6}
	err = svc.Add(e)
	if err != errInvalidRating {
		t.Fatal("Add() should return error")
	}
	e = Entry{Rating: 5}
	err = svc.Add(e)
	if err != nil {
		t.Fatal("Add() should not return error")
	}
}

func TestLatest(t *testing.T) {
	log := log.New("nop", "", true)
	svc := New(log, &mockRepository{})
	if svc == nil {
		t.Errorf("New() == nil")
	}

	e, err := svc.GetLatest(0)
	if err != nil {
		t.Fatal("GetLatest() should not return error")
	}
	if len(e) > 0 {
		t.Fatal("GetLatest() should not return more than", 0)
	}
}
