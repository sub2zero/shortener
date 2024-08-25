package store

import "errors"

var (
	NotFoundErr = errors.New("not found")
)

type UrlStore struct {
	list map[string]ShortUrls
}

func NewUrlStore() *UrlStore {
	list := make(map[string]ShortUrls)
	return &UrlStore{
		list,
	}
}

func (m UrlStore) Add(name string, ShortUrls ShortUrls) error {
	m.list[name] = ShortUrls
	return nil
}

func (m UrlStore) Get(name string) (ShortUrls, error) {

	if val, ok := m.list[name]; ok {
		return val, nil
	}

	return ShortUrls{}, NotFoundErr
}

func (m UrlStore) List() (map[string]ShortUrls, error) {
	return m.list, nil
}

func (m UrlStore) Update(name string, ShortUrls ShortUrls) error {

	if _, ok := m.list[name]; ok {
		m.list[name] = ShortUrls
		return nil
	}

	return NotFoundErr
}

func (m UrlStore) Remove(name string) error {
	delete(m.list, name)
	return nil
}
