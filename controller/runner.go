package controller

import (
	"fmt"
	"log"

	"github.com/bndr/gojenkins"
)

type JobRunner struct {
	c   *Controller
	job Job
}

func (r *JobRunner) Run() {

	log.Println("Entering Run")
	log.Println("Connecting to Jenkins... ")
	jenkins, err := gojenkins.CreateJenkins(r.job.ServerURL,
		r.job.UserName,
		r.job.Password).Init()
	if err != nil {
		panic(err)
	}

	jenkinsJob, err := jenkins.GetJob(r.job.JobName)
	if err != nil {
		panic(err)
	}

	lastBuild, err := jenkinsJob.GetLastBuild()
	if "SUCCESS" == lastBuild.GetResult() {
		fmt.Printf("Last build (ID:%v) succeeded.\n", lastBuild.GetBuildNumber())
		r.job.Value = 1
	} else {
		fmt.Printf("Last build (ID:%v) result = %s.\n", lastBuild.GetBuildNumber(), lastBuild.GetResult())
		r.c.LightOff("green")
		r.c.LightOff("yellow")
		r.c.LightOn("red")
		r.job.Value = 0
	}
}

func (c *Controller) Runner(job Job) (*JobRunner, error) {

	return &JobRunner{
		c:   c,
		job: job,
	}, nil
}
