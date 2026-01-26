package queries

import (
	"os"
	"path/filepath"
)

type Store struct {
	BaseDir string
}

func New(baseDir string) *Store {
	return &Store{BaseDir: baseDir}
}

func (s *Store) Load(name string) (string, error) {
	p := filepath.Join(s.BaseDir, name)
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
