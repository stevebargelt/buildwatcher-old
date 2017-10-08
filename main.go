package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver

	"github.com/bndr/gojenkins"
	"github.com/gorilla/mux"
)

// Light is the structure for a single LED Light attached to your GPIO
type Light struct {
	ID    string          `json:"id"`
	GPIO  int             `json:"integer,gpio"`
	Color string          `json:"color"`
	Desc  string          `json:"desc"`
	State string          `json:"state"`
	dpin  embd.DigitalPin `json:"-"`
}

var lights map[string]*Light

func LightOn(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	light := lights[vars["id"]]

	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&light)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

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

func LightOff(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	light := lights[vars["id"]]

	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&light)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

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

func AddLight(w http.ResponseWriter, req *http.Request) {

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
	defer light.dpin.Close()

	if err := light.dpin.SetDirection(embd.Out); err != nil {
		panic(err)
	}

	retjs, err := json.Marshal(light)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(retjs))

}

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

func main() {

	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)

	logger.Print("Hello, log file!")
	fmt.Print(&buf)

	fmt.Println("Starting")
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	// r := raspi.NewAdaptor()
	//leds := []*LedDriver
	// ledRed := gpio.NewLedDriver(r, "12")    //GPIO18
	// ledYellow := gpio.NewLedDriver(r, "18") //GPIO24
	// ledGreen := gpio.NewLedDriver(r, "13")  //GPIO27
	//buzzer := gpio.NewBuzzerDriver(r, "16") //GPIO23

	var port = 9002

	router := mux.NewRouter()
	logger.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), router))
	logger.Print("Listening on port", port)

	router.HandleFunc("/api/lights/", AddLight).Methods("POST")
	router.HandleFunc("/api/lights/{id}/on", LightOn).Methods("POST")
	router.HandleFunc("/api/lights/{id}/off", LightOff).Methods("POST")
	//router.HandleFunc("/api/lights/{color}", GetLightInfo).Methods("Get")
	router.HandleFunc("/api/lights", GetLights).Methods("Get")

	fmt.Print("Connecting to Jenkins... ")
	jenkins, err := gojenkins.CreateJenkins(nil,
		"https://abs.harebrained-apps.com",
		"stevebargelt",
		"steel2000").Init()
	if err != nil {
		panic("Something Went Wrong")
	}

	jobName := "myretail-aspdotnetcore"
	builds, err := jenkins.GetAllBuildIds(jobName)

	for _, build := range builds {
		buildID := build.Number
		data, err := jenkins.GetBuild(jobName, buildID)
		if err != nil {
			panic(err)
		}

		if "SUCCESS" == data.GetResult() {
			fmt.Printf("%v: This build succeeded\n", buildID)
		} else {
			fmt.Printf("%v: This build failed\n", buildID)
		}
	}
	if err != nil {
		panic("Job Does Not Exist")
	}

	// master := gobot.NewMaster()
	// a := api.NewAPI(master)
	// a.Port = "3001"
	// a.Start()

	// work := func() {
	// 	gobot.Every(1*time.Second, func() {
	// 		ledRed.Toggle()
	// 		ledYellow.Toggle()
	// 		ledGreen.Toggle()
	// 		// buzzer.Toggle()
	// 		// for _, val := range song {
	// 		// 	buzzer.Tone(val.tone, val.duration)
	// 		// 	time.Sleep(10 * time.Millisecond)
	// 		// }
	// 	})
	// }

	// robot := gobot.NewRobot("blinkBot",
	// 	[]gobot.Connection{r},
	// 	[]gobot.Device{ledRed, ledYellow, ledGreen, buzzer},
	// 	work,
	// )

	// master.AddRobot(robot)
	// master.Start()

	//robot.Start()

	//fmt.Println("success. Connected.")
	//fmt.Print("Adding jenkins job... ")

}
