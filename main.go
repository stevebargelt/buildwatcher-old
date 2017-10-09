package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver

	"github.com/gorilla/mux"
)

// Light is the structure for a single LED Light attached to your GPIO
type Light struct {
	ID    string          `json:"id"`
	GPIO  int             `json:"gpio"`
	Color string          `json:"color"`
	Desc  string          `json:"desc"`
	State string          `json:"state"`
	dpin  embd.DigitalPin `json:"-"`
}

type apiServer struct {
	http.Server
	shutdownReq chan bool
	reqCount    uint32
}

var lights map[string]*Light

//NewServer : creates a new API Server
func NewServer() *apiServer {
	//create server
	s := &apiServer{
		Server: http.Server{
			Addr:         ":9002",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		shutdownReq: make(chan bool),
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/lights", AddLight).Methods("POST")
	router.HandleFunc("/api/lights/{id}/on", LightOn).Methods("POST")
	router.HandleFunc("/api/lights/{id}/off", LightOff).Methods("POST")
	//router.HandleFunc("/api/lights/{color}", GetLightInfo).Methods("Get")
	router.HandleFunc("/api/lights", GetLights).Methods("Get")
	router.HandleFunc("/api/shutdown", s.ShutdownHandler)

	s.Handler = router

	return s
}

//LightOn : Turns a light on through GPIO
func LightOn(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	light := lights[vars["id"]]
	fmt.Println("Turning on... Light ID = ", light.ID)
	light.State = "on"

	if err := light.dpin.Write(embd.High); err != nil {
		panic(err)
	}

	retjs, err := json.Marshal(light)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(rw, string(retjs))

}

//LightOff : Turns a light off through GPIO
func LightOff(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	light := lights[vars["id"]]
	fmt.Println("Turning off... Light ID = ", light.ID)
	light.State = "off"

	if err := light.dpin.Write(embd.Low); err != nil {
		panic(err)
	}

	retjs, err := json.Marshal(light)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(rw, string(retjs))
}

//AddLight - adds a light to the system
func AddLight(w http.ResponseWriter, req *http.Request) {

	fmt.Println("Add light called")

	light := new(Light)

	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&light)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lights[light.ID] = light

	light.dpin, err = embd.NewDigitalPin(light.GPIO)
	if err != nil {
		panic(err)
	}

	if err := light.dpin.SetDirection(embd.Out); err != nil {
		log.Println("light.dpin.SetDirection(embd.Out) failed - just a warning")
	}

	retjs, err := json.Marshal(light)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(retjs))

}

//GetLights - gets all the lights and returns JSON
func GetLights(rw http.ResponseWriter, req *http.Request) {

	v := make([]*Light, 0, len(lights))

	for _, value := range lights {
		v = append(v, value)
	}
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(rw, string(js))
}

func (s *apiServer) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}

	log.Printf("Turning off and cleaning up all lights ...")
	for _, light := range lights {
		if err := light.dpin.Write(embd.Low); err != nil {
			panic(err)
		}
		light.dpin.Close()
	}

	log.Printf("Stoping http server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}

func (s *apiServer) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutdown server"))

	//Do nothing if shutdown request already issued
	//if s.reqCount == 0 then set to 1, return true otherwise false
	if !atomic.CompareAndSwapUint32(&s.reqCount, 0, 1) {
		log.Printf("Shutdown through API call in progress...")
		return
	}

	go func() {
		s.shutdownReq <- true
	}()
}

func main() {

	lights = make(map[string]*Light)

	//create your file with desired read/write permissions
	f, err := os.OpenFile("buildwatcher.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()
	//set output of logs to f
	log.SetOutput(f)
	//test case

	log.Print("Starting embd...")
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()
	log.Print("embd started...")

	server := NewServer()

	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen and serve: %v", err)
		}
		done <- true
	}()

	//wait shutdown
	server.WaitShutdown()

	<-done
	log.Printf("DONE!")

}

// Red    "12" //GPIO18
// Yellow "18" //GPIO24
// Green  "13" //GPIO27
// buzzer "16" //GPIO23
