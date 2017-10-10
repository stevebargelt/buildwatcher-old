package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const JobBucket = "jobs"

type Job struct {
	ID        string `json:"id"`
	Minute    string `json:"minute"`
	Day       string `json:"day"`
	Hour      string `json:"hour"`
	Second    string `json:"second"`
	On        bool   `json:"on"`
	Value     int    `json:"value"`
	Name      string `json:"name"`
	ServerURL string `json:"serverurl"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	JobName   string `json:"jobname"`
}

func (c *Controller) GetJob(id string) (Job, error) {
	var job Job
	return job, c.store.Get(JobBucket, id, &job)
}

func (c *Controller) ListJobs() (*[]interface{}, error) {
	fn := func(v []byte) (interface{}, error) {
		var job Job
		return &job, json.Unmarshal(v, &job)
	}
	return c.store.List(JobBucket, fn)
}

func (c *Controller) CreateJob(job Job) error {

	log.Println("Entering CreateJob")
	fn := func(id string) interface{} {
		job.ID = id
		return job
	}
	if err := c.store.Create(JobBucket, fn); err != nil {
		return err
	}
	return c.addToCron(job)
}

func (c *Controller) UpdateJob(id string, payload Job) error {
	return c.store.Update(JobBucket, id, payload)
}

func (c *Controller) DeleteJob(id string) error {
	if err := c.store.Delete(JobBucket, id); err != nil {
		return err
	}
	return c.deleteFromCron(id)
}

func (c *Controller) loadAllJobs() error {
	jobs, err := c.ListJobs()
	if err != nil {
		return err
	}
	if jobs == nil {
		log.Printf("No jobs present")
		return nil
	}
	for _, rawJob := range *jobs {
		job, ok := rawJob.(*Job)
		if !ok {
			log.Println("ERROR: Failed to typecast to job")
			continue
		}
		if err := c.addToCron(*job); err != nil {
			log.Println("ERROR: Failed to add job in cron runner. Error:", err)
		}
	}
	return nil
}

func (c *Controller) addToCron(job Job) error {
	log.Println("Entering addToCron")
	cronSpec := strings.Join([]string{job.Second, job.Minute, job.Hour, job.Day, "*", "?"}, " ")
	log.Printf("Cronspec: %s\n", cronSpec)
	runner, err := c.Runner(job)
	if err != nil {
		return err
	}
	log.Printf("cronRunner.AddJob(%s, runner)\n", cronSpec)
	cronID, err := c.cronRunner.AddJob(cronSpec, runner)
	//cronID, err := c.cronRunner.AddJob("* 1 * * * *", runner)
	if err != nil {
		return err
	}
	log.Println("Successfully added cron entry. ID:", cronID)
	c.cronIDs[job.ID] = cronID
	return nil
}

func (c *Controller) deleteFromCron(jobID string) error {
	id, ok := c.cronIDs[jobID]
	if !ok {
		return fmt.Errorf("Cron ID not found for job ID:%s", jobID)
	}
	c.cronRunner.Remove(id)
	return nil
}
