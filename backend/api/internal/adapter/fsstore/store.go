package fsstore

import "context"

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}

func (f *Store) SaveFileBody(ctx context.Context, body []byte) (url string, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *Store) DeleteFile(ctx context.Context, url string) error {
	//TODO implement me
	panic("implement me")
}
