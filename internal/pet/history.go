package pet

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Record struct {
	At        time.Time `json:"at"`
	User      string    `json:"user"`
	Action    Action    `json:"action"`
	Food      float64   `json:"food"`
	Water     float64   `json:"water"`
	Affection float64   `json:"affection"`
}

func (s *State) Record(a Action, user string, now time.Time) Record {
	return Record{
		At:        now.UTC().Truncate(time.Second),
		User:      user,
		Action:    a,
		Food:      round2(s.Food),
		Water:     round2(s.Water),
		Affection: round2(s.Affection),
	}
}

func AppendHistory(path string, r Record) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(b, '\n'))
	return err
}

func RecentlyActed(path, user string, now time.Time) (bool, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		var r Record
		if json.Unmarshal(sc.Bytes(), &r) != nil {
			continue
		}
		if strings.EqualFold(r.User, user) && now.Sub(r.At) < RateLimit {
			return true, nil
		}
	}
	return false, sc.Err()
}
