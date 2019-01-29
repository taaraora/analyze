package storage

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type mockStorage struct {
	data map[string][]byte
	isBroken bool
}

var errBroken = errors.New("internal storage error")

func GetMockStorage(t *testing.T, data map[string]string)Interface {
	t.Helper()
	return 	NewMockStorage(data, false)
}

func GetMockBrokenStorage(t *testing.T)Interface {
	t.Helper()
	return 	NewMockStorage(nil, true)
}


func NewMockStorage(data map[string]string, isBroken bool) Interface {
	result := map[string][]byte{}
	for key, value := range data {
		result[key] = []byte(value)
	}

	return &mockStorage{
		data: result,
		isBroken: isBroken,
	}
}

func (s *mockStorage)GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	if s.isBroken {
		return nil, errBroken
	}
	result := [][]byte{}
	for key, _ := range s.data {
		if strings.Contains(key, prefix){
			result = append(result, s.data[key])
		}
	}

	return result, nil
}

func (s *mockStorage)Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	if s.isBroken {
		return nil, errBroken
	}
	v, ok := s.data[prefix+key]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}

func (s *mockStorage)Put(ctx context.Context, prefix string, key string, value []byte) error {
	if s.isBroken {
		return errBroken
	}
	s.data[prefix+key] = value

	return nil
}

func (s *mockStorage)Delete(ctx context.Context, prefix string, key string) error {
	if s.isBroken {
		return errBroken
	}
	delete(s.data, prefix+key)
	return nil
}

func (s *mockStorage)Close() error {
	if s.isBroken {
		return errBroken
	}
	return nil
}
