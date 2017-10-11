package controller

import (
	"fmt"
	"log"

	"github.com/kidoman/embd"
)

const LightsBucket = "lights"

// Light is the structure for a single LED Light attached to your GPIO
type Light struct {
	ID    string          `json:"id"`
	GPIO  int             `json:"gpio"`
	Color string          `json:"color"`
	Desc  string          `json:"desc"`
	State string          `json:"state"`
	Dpin  embd.DigitalPin `json:"-"`
}

func (c *Controller) GetLight(name string) (Light, error) {
	light, ok := c.config.Lights[name]
	if !ok {
		return light, fmt.Errorf("No light named '%s' present", name)
	}
	return light, nil
}

//GetLights - gets all the lights and returns JSON
func (c *Controller) GetLights() (*[]interface{}, error) {
	list := []interface{}{}
	for _, l := range c.config.Lights {
		l1 := l
		list = append(list, &l1)
	}
	return &list, nil
}

func (c *Controller) ConfigureLight(id string, on bool, value int) error {

	l, ok := c.config.Lights[id]
	if !ok {
		return fmt.Errorf("Light named: '%s' does noy exist", id)
	}

	if c.config.DevMode {
		log.Println("Dev mode on. Skipping:", id, "On:", on, "Value:", value)
		return nil
	}

	return c.doSwitching(l.GPIO, on)

}

//LightOn : Turns a light on through GPIO
func (c *Controller) LightOn(id string) error {

	log.Printf("Called controller.light.LightOn for %v", id)
	l, ok := c.config.Lights[id]
	if !ok {
		return fmt.Errorf("Light named: '%s' does not exist", id)
	}

	l.State = "on"
	return c.doSwitching(l.GPIO, true)

}

//LightOff : Turns a light off through GPIO
func (c *Controller) LightOff(id string) error {

	l, ok := c.config.Lights[id]
	if !ok {
		return fmt.Errorf("Light named: '%s' does not exist", id)
	}

	l.State = "off"
	return c.doSwitching(l.GPIO, false)

}

// //AddLight - adds a light to the system
// func (c *Controller) AddLight(w http.ResponseWriter, req *http.Request) {

// 	light := new(Light)

// 	dec := json.NewDecoder(req.Body)
// 	err := dec.Decode(&light)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	Lights[light.ID] = light

// 	light.dpin, err = embd.NewDigitalPin(light.GPIO)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := light.dpin.SetDirection(embd.Out); err != nil {
// 		log.Println("light.dpin.SetDirection(embd.Out) failed - just a warning")
// 	}

// 	retjs, err := json.Marshal(light)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprint(w, string(retjs))

// }
