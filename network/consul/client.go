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
	client, err := NewClient(name, cfg)
	if err != nil { return nil, err}

	// save
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
	agent				*Agent
	catalog				*Catalog		// unfinished
}

func NewClient(name string, cfg *_consul.Config) (*Client, error) {
	c, err := _consul.NewClient(cfg)
	if err != nil { return nil, err}
	client := &Client{
		name:   name,
		client: c,
	}
	client.agent = NewAgent(client)
	client.catalog = NewCatalog(client)
	return client, nil
}

func (h *Client) Name() string { return h.name }
func (h *Client) Agent() *Agent { return h.agent }
func (h *Client) Catalog() *Catalog { return h.catalog }




func listenHealthCheck(regCheckUrl string, handler func (ctx *http.RequestCtx)) error{

	u, err := url.Parse(regCheckUrl)
	if err != nil { return fmt.Errorf("Heath checking url parse error, %v", err ) }
	s := strings.Split(u.Host, ":")
	address := ":80"
	if len(s) >= 2 { address = fmt.Sprintf(":%v", s[1])}


	defaultHandler := func (ctx *http.RequestCtx){
		fmt.Fprintf(ctx, "check")
	}
	if handler == nil { handler = defaultHandler }

	t := task.NewTask(func() {
		if err := http.ListenAndServe(address, handler); err != nil {
			log.Errorw("Consul health service end with error", "error", err)
		}
	})
	t.SetSkipMonitor(true)
	t.Start()

	return nil
}

