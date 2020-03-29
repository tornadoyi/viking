package consul

import (
	"fmt"
	_consul "github.com/hashicorp/consul/api"
	"github.com/tornadoyi/viking/http"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/task"
	"net/url"
	"strings"
	"sync"
	"time"
)


var (
	clients 	=  		map[string]*Client{}
	mutex		=		sync.RWMutex{}
)


func CreateClient(name string, cfg *_consul.Config) (*Client, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// check
	if _, ok := clients[name]; ok { return nil, fmt.Errorf("Repeated client %v", name)}

	// create client
	c, err := _consul.NewClient(cfg)
	if err != nil { return nil, err}

	// save
	client := &Client{name, c, nil, nil,
		make(map[string]*_consul.AgentService), &sync.RWMutex{}}
	clients[name] = client
	return client, nil
}


func GetClient(name string) (*Client, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	c, ok := clients[name]
	return c, ok
}




type Client struct {
	name				string
	client				*_consul.Client
	registration		*_consul.AgentServiceRegistration
	timer				*time.Timer
	services			map[string]*AgentService
	mutex				*sync.RWMutex
}

func (h *Client) RegisterServer(regCfg *AgentServiceRegistrationConfig) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// check
	if h.registration != nil { return fmt.Errorf("Repeated server register by client %v", h.name) }
	if regCfg == nil { return fmt.Errorf("Empty registration config")}
	registration := regCfg.AgentServiceRegistration()
	h.registration = registration

	// register
	err := h.client.Agent().ServiceRegister(registration)
	if err != nil { return err }

	// start health checking server
	u, err := url.Parse(registration.Check.HTTP)
	if err != nil { return fmt.Errorf("Heath checking url parse error, %v", err ) }
	s := strings.Split(u.Host, ":")
	address := ":80"
	if len(s) >= 2 { address = fmt.Sprintf(":%v", s[1])}

	checkHandler := func (ctx *http.RequestCtx){
		fmt.Fprintf(ctx, "check")
	}
	t := task.NewTask(func() {
		if err := http.ListenAndServe(address, checkHandler); err != nil {
			log.Error(err)
		}
	})
	t.Start()

	return nil
}


func (h *Client) SetFetchInterval(interval time.Duration) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.timer != nil { h.timer.Stop() }
	h.timer = time.AfterFunc(interval, func() {
		h.mutex.Lock()
		defer h.mutex.Unlock()
		defer h.timer.Reset(interval)
		services, err := h.client.Agent().Services()
		if err != nil { log.Error(err); return }
		h.services = services
	})
}

func (h *Client) FetchServices() (map[string]*AgentService, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	services, err := h.client.Agent().Services()
	if err != nil { return nil, err}
	h.services = services
	return services, nil
}



type AgentService = _consul.AgentService