package controller

import (
	"encoding/json"
	"log"
)

const ProjectBucket = "projects"

type Project struct {
	ID          string `json:"id"`
	Source      string `json:"source"`
	Job         string `json:"job"`
	URL         string `json:"url"`
	Status      string `json:"status"`
	BuildNumber string `json"buildnumber"`
}

func (c *Controller) GetProject(id string) (Project, error) {
	var project Project
	return project, c.store.Get(ProjectBucket, id, &project)
}

func (c *Controller) ListProjects() (*[]interface{}, error) {
	fn := func(v []byte) (interface{}, error) {
		var project Project
		return &project, json.Unmarshal(v, &project)
	}
	return c.store.List(ProjectBucket, fn)
}

func (c *Controller) CreateProject(project Project) error {

	log.Println("Entering CreateProject")
	fn := func(id string) interface{} {
		project.ID = id
		return project
	}
	return c.store.Create(ProjectBucket, fn)
}

func (c *Controller) UpdateProject(id string, payload Project) error {
	return c.store.Update(ProjectBucket, id, payload)
}

func (c *Controller) DeleteProject(id string) error {
	return c.store.Delete(ProjectBucket, id)
}
