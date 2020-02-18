package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type demoservice struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Desc string `json:"price"`
}

type ServiceConfiguration struct {
	Categories []string `json:"categories"`
}

func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "demoservice"
	registration.Name = "demoservice"
	port, err := strconv.Atoi(port()[1:len(port())])
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
}

func Configuration(w http.ResponseWriter, r *http.Request) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	kvpair, _, err := consul.KV().Get("demoservice-configuration", nil)
	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	if kvpair.Value == nil {
		fmt.Fprintf(w, "Configuration empty")
		return
	}
	val := string(kvpair.Value)
	fmt.Fprintf(w, "%s", val)

}

func Services(w http.ResponseWriter, r *http.Request) {
	demoservices := []demoservice{
		{
			ID:   1,
			Name: "demoService",
			Desc: "demo service",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&demoservices)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `demoservice demoservice is good`)
}

func port() string {
	p := os.Getenv("PRODUCT_SERVICE_PORT")
	h := os.Getenv("PRODUCT_SERVICE_HOST")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8100"
	}
	return fmt.Sprintf("%s:%s", h, p)
}

func main() {
	registerServiceWithConsul()
	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc("/demoservices", Services)
	http.HandleFunc("/demoservice/config", Configuration)
	fmt.Printf("demoservice demoservice is up on port: %s", port())
	http.ListenAndServe(port(), nil)
}

// todo : now call the prom discovery manager and try to get the service in it
// todo : remove the service as in deregister it from main and see if its really updated in discovery interface
