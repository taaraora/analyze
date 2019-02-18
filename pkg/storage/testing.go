package storage

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type mockStorage struct {
	data     map[string][]byte
	isBroken bool
}

type mockMsg []byte

func (m mockMsg) Payload() []byte {
	return m
}

var errBroken = errors.New("internal storage error")

func GetMockStorage(t *testing.T, data map[string]string) Interface {
	t.Helper()
	return NewMockStorage(data, false)
}

func GetMockBrokenStorage(t *testing.T) Interface {
	t.Helper()
	return NewMockStorage(nil, true)
}

func NewMockStorage(data map[string]string, isBroken bool) Interface {
	result := map[string][]byte{}
	for key, value := range data {
		result[key] = []byte(value)
	}

	return &mockStorage{
		data:     result,
		isBroken: isBroken,
	}
}

func (s *mockStorage) GetAll(ctx context.Context, prefix string) ([]Message, error) {
	if s.isBroken {
		return nil, errBroken
	}
	result := []Message{}
	for key, _ := range s.data {
		if strings.Contains(key, prefix) {
			result = append(result, mockMsg(s.data[key]))
		}
	}

	return result, nil
}

func (s *mockStorage) Get(ctx context.Context, prefix string, key string) (Message, error) {
	if s.isBroken {
		return nil, errBroken
	}
	v, ok := s.data[prefix+key]
	if !ok {
		return nil, ErrNotFound
	}

	return mockMsg(v), nil
}

func (s *mockStorage) Put(ctx context.Context, prefix string, key string, value Message) error {
	if s.isBroken {
		return errBroken
	}
	s.data[prefix+key] = value.Payload()

	return nil
}

func (s *mockStorage) Delete(ctx context.Context, prefix string, key string) error {
	if s.isBroken {
		return errBroken
	}
	delete(s.data, prefix+key)
	return nil
}

func (s *mockStorage) Close() error {
	if s.isBroken {
		return errBroken
	}
	return nil
}

func (s *mockStorage) WatchRange(ctx context.Context, key string) <-chan WatchEvent {
	panic("not implemented")
}
