package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type cache struct {
	Context context.Context
}

func Instance(ctx context.Context) *cache {
	return &cache{Context: ctx}
}

func (c *cache) getKey(s *Stream) *datastore.Key {
	uid := fmt.Sprintf("%v", s.StreamID)
	return datastore.NewKey(c.Context, "Streams", uid, 0, nil)
}

func (c *cache) IsStreamInStore(s *Stream) bool {
	var single Stream

	if err := datastore.Get(c.Context, c.getKey(s), &single); err != nil {
		if err != datastore.ErrNoSuchEntity {
			log.Println(err.Error())
		}

		return false
	}

	return true
}

func (c *cache) SaveStreamToStore(s *Stream) error {
	if _, err := datastore.Put(c.Context, c.getKey(s), s); err != nil {
		return err
	}

	return nil
}
