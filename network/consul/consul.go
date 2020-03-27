package consul

import (
	"fmt"
	_consul "github.com/hashicorp/consul/api"
	"github.com/tornadoyi/viking/http"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/task"
)


var (
	servers 	=  		map[string]*Server{}
)


func RegisterServer(cfg *_consul.Config, regCfg *AgentServiceRegistrationConfig) error {

	// check
	if regCfg == nil { return fmt.Errorf("Empty registration config")}
	registration := regCfg.AgentServiceRegistration()
	if _, ok := servers[registration.Name]; ok { return fmt.Errorf("Repeated registration %v", registration.Name)}

	// create client
	client, err := _consul.NewClient(cfg)
	if err != nil { return err}

	// register
	err = client.Agent().ServiceRegister(registration)
	if err != nil { return err }

	// start health checking server
	checkHandler := func (ctx *http.RequestCtx){
		fmt.Fprintf(ctx, "check")
	}
	t := task.Create(func() {
		if err := http.ListenAndServe(registration.Check.HTTP, checkHandler); err != nil {
			log.Error(err)
		}
	})
	t.Start()

	// save
	servers[registration.Name] = &Server{client, registration}
	return nil
}



type Server struct {
	client				*_consul.Client
	registration		*_consul.AgentServiceRegistration
}