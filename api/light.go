package api

import (
	"log"
	"net/http"

	"github.com/stevebargelt/buildwatcher/controller"
)

func (h *APIHandler) GetLight(w http.ResponseWriter, r *http.Request) {
	fn := func(id string) (interface{}, error) {
		return h.controller.GetLight(id)
	}
	h.jsonGetResponse(fn, w, r)
}

func (h *APIHandler) GetLights(w http.ResponseWriter, r *http.Request) {
	fn := func() (interface{}, error) {
		return h.controller.GetLights()
	}
	h.jsonListResponse(fn, w, r)
}

func (h *APIHandler) LightOn(w http.ResponseWriter, r *http.Request) {
	log.Println("Called api.light.LightOn")
	var l controller.Light
	fn := func(id string) error {
		l.ID = id
		return h.controller.LightOn(id)
	}
	h.jsonToggleResponse(&l, fn, w, r)
}

func (h *APIHandler) LightOff(w http.ResponseWriter, r *http.Request) {
	var l controller.Light
	fn := func(id string) error {
		return h.controller.LightOff(id)
	}
	h.jsonToggleResponse(&l, fn, w, r)
}

// func (h *APIHandler) AddLight(w http.ResponseWriter, r *http.Request) {
// 	var l controller.Light
// 	fn := func() error {
// 		return h.controller.AddLight(l)
// 	}
// 	h.jsonCreateResponse(&l, fn, w, r)
// }

// func (h *APIHandler) ConfigureLight(w http.ResponseWriter, r *http.Request) {
// 	fn := func(id string) error {
// 		return h.controller.ConfigureOutlet(id, a.On, a.Value)
// 	}
// 	h.jsonUpdateResponse(&a, fn, w, r)
// }
