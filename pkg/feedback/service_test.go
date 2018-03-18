package feedback

import (
	"reflect"
	"testing"

	"github.com/playnet-public/libs/log"
)

func TestNew(t *testing.T) {
	log := log.NewNop()
	if New(log, nil) == nil {
		t.Errorf("New() == nil")
	}
}

func TestService_Add(t *testing.T) {
	log := log.NewNop()
	svc := New(log, newMockRepository(nil, nil, nil))
	if svc == nil {
		t.Errorf("New() == nil")
	}

	e := Entry{Rating: 0}
	err := svc.Add(e)
	if err != ErrInvalidRating {
		t.Fatal("Add() should return error")
	}
	e = Entry{Rating: 6}
	err = svc.Add(e)
	if err != ErrInvalidRating {
		t.Fatal("Add() should return error")
	}
	e = Entry{Rating: 5}
	err = svc.Add(e)
	if err != nil {
		t.Fatal("Add() should not return error")
	}
}

func TestLatest(t *testing.T) {
	log := log.NewNop()
	svc := New(log, newMockRepository(nil, nil, nil))
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

func TestService_GetLatestFiltered(t *testing.T) {
	log := log.New("nop", "", true)

	type args struct {
		n      uint
		filter int
	}
	tests := []struct {
		name    string
		args    args
		getFunc func(uint, int) ([]Entry, error)
		want    []Entry
		wantErr bool
	}{
		{
			"basic",
			args{1, 0},
			func(n uint, f int) ([]Entry, error) {
				return []Entry{}, nil
			},
			[]Entry{},
			false,
		},
		{
			"filterBy1",
			args{1, 1},
			func(n uint, f int) ([]Entry, error) {
				return []Entry{
					{"1", "1", "1", 1, ""},
					{"3", "3", "1", 1, ""},
					{"5", "5", "1", 1, ""},
				}, nil
			},
			[]Entry{
				{"1", "1", "1", 1, ""},
				{"3", "3", "1", 1, ""},
				{"5", "5", "1", 1, ""},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(log, newMockRepository(nil, nil, tt.getFunc))
			got, err := svc.GetLatestFiltered(tt.args.n, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetLatestFiltered() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetLatestFiltered() = %v, want %v", got, tt.want)
			}
		})
	}
}
