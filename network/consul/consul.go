package consul

import (
	"fmt"
	_consul "github.com/hashicorp/consul/api"
	"github.com/tornadoyi/viking/http"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/task"
)




func RegisterServer(regCfg *AgentServiceRegistrationConfig) error {

	// create client
	client, err := _consul.NewClient(_consul.DefaultConfig())
	if err != nil { return err}

	// register
	registration := regCfg.AgentServiceRegistration()
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

	return nil
}



type Server struct {
	client				*_consul.Client
	registration		*_consul.AgentServiceRegistration
}