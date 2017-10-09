package controller

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
	"gopkg.in/robfig/cron.v2"
)

type Controller struct {
	store      *Store
	cronRunner *cron.Cron
	cronIDs    map[string]cron.EntryID
	config     Config
	state      *State
}

func New(config Config) (*Controller, error) {
	db, err := bolt.Open("buildwatcher.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	store := NewStore(db)
	c := &Controller{
		store:      store,
		state:      NewState(config, store),
		cronRunner: cron.New(),
		cronIDs:    make(map[string]cron.EntryID),
		config:     config,
	}
	return c, nil
}

func (c *Controller) CreateBuckets() error {
	buckets := []string{
		LightsBucket,
		// JobBucket,
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
	c.cronRunner.Start()
	// if err := c.loadAllJobs(); err != nil {
	// 	return err
	// }
	log.Println("Started Controller")
	return nil
}

func (c *Controller) Stop() error {
	c.cronRunner.Stop()
	c.state.TearDown()
	c.store.Close()
	c.logStopTime()
	log.Println("Stopped Controller")
	return nil
}
