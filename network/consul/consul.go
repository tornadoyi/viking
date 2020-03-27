package consul

import (
	"fmt"
	_consul "github.com/hashicorp/consul/api"
	"github.com/tornadoyi/viking/http"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/task"
	"strings"
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
	s := strings.Split(registration.Check.HTTP, ":")
	if len(s) < 2 { return fmt.Errorf("Invalid address format for http heath checking, address: %v",registration.Check.HTTP ) }
	address := strings.Split(s[1],"/")[0]

	checkHandler := func (ctx *http.RequestCtx){
		fmt.Fprintf(ctx, "check")
	}
	t := task.Create(func() {
		if err := http.ListenAndServe(address, checkHandler); err != nil {
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