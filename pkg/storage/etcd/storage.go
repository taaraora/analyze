package etcd

import (
	"context"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/sirupsen/logrus"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"

	"github.com/supergiant/analyze/pkg/storage"
)

type Storage struct {
	cfg    clientv3.Config
	client *clientv3.Client
	logger logrus.FieldLogger
}

type msg []byte

func (m msg) Payload() []byte {
	return m
}

type watchEvent struct {
	msg
	err       error
	eventType storage.WatchEventType
}

func (we *watchEvent) Err() error {
	return we.err
}

func (we *watchEvent) Type() storage.WatchEventType {
	return we.eventType
}

func (e *Storage) Get(ctx context.Context, prefix string, key string) (storage.Message, error) {
	res, err := e.client.Get(ctx, prefix+key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from the etcd")
	}
	if res.Count == 0 {
		return nil, storage.ErrNotFound
	}
	return msg(res.Kvs[0].Value), nil
}

func (e *Storage) Put(ctx context.Context, prefix string, key string, value storage.Message) error {
	_, err := e.client.Put(ctx, prefix+key, string(value.Payload()))
	return errors.Wrap(err, "failed to write to the etcd")
}

func (e *Storage) Delete(ctx context.Context, prefix string, key string) error {
	_, err := e.client.Delete(ctx, prefix+key, clientv3.WithPrefix())
	return errors.Wrap(err, "failed to delete kv from the etcd")
}

func (e *Storage) GetAll(ctx context.Context, prefix string) ([]storage.Message, error) {
	result := make([]storage.Message, 0)

	r, err := e.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return result, errors.Wrap(err, "failed to read from the etcd")
	}
	for _, v := range r.Kvs {
		result = append(result, msg(v.Value))
	}
	return result, nil
}

func NewETCDStorage(cfg clientv3.Config, logger logrus.FieldLogger) (storage.Interface, error) {
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to the etcd")
	}

	return &Storage{
		cfg:    cfg,
		client: client,
		logger: logger,
	}, nil
}

func (e *Storage) Close() error {
	return e.client.Close()
}

// TODO: etcd does not ensure linearizability for watch operations. revisit this logic in future
func (e *Storage) WatchRange(ctx context.Context, key string) <-chan storage.WatchEvent {
	w := clientv3.NewWatcher(e.client)
	watchChan := w.Watch(ctx, key, clientv3.WithPrefix() /*clientv3.WithProgressNotify()*/)
	results := make(chan storage.WatchEvent)
	values, err := e.GetAll(ctx, key)
	if err != nil {
		e.logger.Errorf("got error at loading initial values, for k: %v, error: %v", key, err)
		results <- &watchEvent{
			eventType: storage.Error,
			err:       errors.Errorf("got error at loading initial values, for k: %v, error: %v", key, err),
		}
	}

	initialSyncDone := make(chan struct{})
	go func() {
		for _, v := range values {
			results <- &watchEvent{
				msg:       v.Payload(),
				err:       nil,
				eventType: storage.Added,
			}
		}
		initialSyncDone <- struct{}{}
	}()

	go func() {
		<-initialSyncDone
		for v := range watchChan {
			e.logger.Infof("got watch message: %+v", v)
			we := &watchEvent{
				eventType: storage.Unknown,
			}

			if err := v.Err(); err != nil {
				e.logger.Errorf("got watch error: %v", err)
				we.eventType = storage.Error
				we.err = err
				results <- we
				continue
			}
			if v.Canceled {
				results <- we
				break
			}

			for _, event := range v.Events {
				we = &watchEvent{
					eventType: storage.Unknown,
				}

				if event.IsCreate() {
					we.eventType = storage.Added
				}
				if event.IsModify() {
					we.eventType = storage.Modified
				}
				if event.Type == mvccpb.DELETE {
					we.eventType = storage.Deleted
				}
				if event.Kv == nil {
					e.logger.Errorf("got nil kv, for kv: %v", event)
					continue
				}
				if len(event.Kv.Key) == 0 {
					e.logger.Errorf("got empty key, for kv: %v", event)
					continue
				}
				we.msg = event.Kv.Value
				results <- we
			}
		}
		close(results)
	}()

	return results
}
