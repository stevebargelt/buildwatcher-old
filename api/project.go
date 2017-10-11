package api

import (
	"net/http"

	"github.com/stevebargelt/buildwatcher/controller"
)

func (h *APIHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	fn := func(id string) (interface{}, error) {
		return h.controller.GetProject(id)
	}
	h.jsonGetResponse(fn, w, r)
}

func (h *APIHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	fn := func() (interface{}, error) {
		return h.controller.ListProjects()
	}
	h.jsonListResponse(fn, w, r)
}

func (h *APIHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var p controller.Project
	fn := func() error {
		return h.controller.CreateProject(p)
	}
	h.jsonCreateResponse(&p, fn, w, r)
}

func (h *APIHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var p controller.Project
	fn := func(id string) error {
		p.ID = id
		return h.controller.UpdateProject(id, p)
	}
	h.jsonUpdateResponse(&p, fn, w, r)
}

func (h *APIHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	h.jsonDeleteResponse(h.controller.DeleteProject, w, r)
}
