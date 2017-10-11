package controller

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type Controller struct {
	store  *Store
	config Config
	state  *State
}

func New(config Config) (*Controller, error) {
	db, err := bolt.Open("buildwatcher.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	store := NewStore(db)
	c := &Controller{
		store:  store,
		state:  NewState(config, store),
		config: config,
	}
	return c, nil
}

func (c *Controller) CreateBuckets() error {
	buckets := []string{
		LightsBucket,
		ProjectBucket,
		UptimeBucket,
	}
	for _, bucket := range buckets {
		if err := c.store.CreateBucket(bucket); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) Start() error {
	if err := c.CreateBuckets(); err != nil {
		return err
	}
	c.logStartTime()
	c.state.Bootup()
	log.Println("Started Controller")
	return nil
}

func (c *Controller) Stop() error {
	c.state.TearDown()
	c.store.Close()
	c.logStopTime()
	log.Println("Stopped Controller")
	return nil
}
