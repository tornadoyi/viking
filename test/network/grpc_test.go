package network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tornadoyi/viking/network/grpc"
	"github.com/tornadoyi/viking/network/grpc/resolver"
	"gopkg.in/yaml.v2"
	"sync"
	"testing"
	"time"
)


func TestServerCreator(t *testing.T) {

	createServer := func(index int) *grpc.Server{
		cfgFormat := `
network:	tcp
address:	:%v
`
		cfgContent := fmt.Sprintf(cfgFormat, 9000+index)
		type ServerConfig struct {
			Network				string								`yaml:"network"`
			Address				string								`yaml:"address"`
			Option				*grpc.ServerConfig					`yaml:"option"`
		}

		var cfg ServerConfig
		err := yaml.Unmarshal([]byte(cfgContent), &cfg)
		if err != nil { t.Fatal(err)}
		server, err := grpc.CreateServer(fmt.Sprintf("test_%v", index), cfg.Network, cfg.Address, cfg.Option.ServerOptions()...)
		if err != nil { t.Fatal(err)}

		service := &EchoService{Index: index}
		if err := server.RegisterService(service, RegisterEchoServer); err != nil { t.Fatal(err) }

		return server
	}

	for i:=0; i<10; i++{
		server := createServer(i)
		go func() {
			server.Serve()
		}()
	}


}



func TestResolverCreator(t *testing.T) {
	cfgContent := `
- scheme: scheme_1
  address:
    echo_1:
    - localhost:9000
    - localhost:9001
    - localhost:9002
    - localhost:9003
    - localhost:9004
    - localhost:9005
    - localhost:9006
    - localhost:9007
    - localhost:9008
    - localhost:9009
- scheme: scheme_2
  address:
    echo_2:
    - localhost:8001
    - localhost:8002
`
	var rsvCfg []resolver.ResolverBuilderConfig
	if err := yaml.Unmarshal([]byte(cfgContent), &rsvCfg); err != nil { t.Fatal(err) }

	for _, r := range rsvCfg {
		resolver.Register(r.ResolverBuilder())
	}
}


func TestClientCreator(t *testing.T) {
	cfgContent := `
address:	"scheme_1:///echo_1"
option:
  insecure: true
  balance_name:	round_robin

`
	type ClientConfig struct {
		Address				string								`yaml:"address"`
		Option				*grpc.DialConfig					`yaml:"option"`
	}

	var cfg ClientConfig
	err := yaml.Unmarshal([]byte(cfgContent), &cfg)
	if err != nil { t.Fatal(err)}

	_, err = grpc.CreateClient("test", cfg.Address, NewEchoClient, cfg.Option.DialOptions()...)
	if err != nil { t.Fatal(err) }

}



func TestPerformance(t *testing.T) {
	client, ok := grpc.GetClient("test")
	if !ok { t.Fatal("get client failed") }
	c := client.Service().(EchoClient)

	// statistics
	receives := make(map[int]int, 0)
	mutex := &sync.Mutex{}

	const (
		numClient			=	100
		sendPerClient		= 	100
	)

	sendOne := func(index int) error{
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		js := map[string]interface{}{
			"client": index,
		}
		msg, _ := json.Marshal(js)
		r, err := c.UnaryEcho(ctx, &EchoRequest{Message: string(msg)})
		if err != nil { return err}

		var rjs map[string]interface{}
		if err := json.Unmarshal([]byte(r.Message), &rjs); err != nil { return err}
		rIndex, ok := rjs["server"]
		if !ok {return err}

		sIndex := int(rIndex.(float64))

		mutex.Lock()
		v, _ := receives[sIndex];
		receives[sIndex] = v + 1
		mutex.Unlock()

		return nil
	}

	group := sync.WaitGroup{}
	group.Add(numClient)
	for i:=0; i<numClient; i++{
		go func(index int) {
			for j:=0; j<sendPerClient; j++{
				err := sendOne(index)
				if err != nil {t.Error(err)}
			}
			group.Done()
		}(i)
	}

	group.Wait()

	for index, count := range receives{
		t.Log(fmt.Sprintf("receive %v messages from server %v", count, index))
	}

}


type EchoService struct {
	EchoServer
	Index				int
}

func (h* EchoService) UnaryEcho(ctx context.Context, request *EchoRequest) (*EchoResponse, error){
	js := map[string]interface{}{
		"server": h.Index,
	}
	bs, err := json.Marshal(js)
	if err != nil { return nil, errors.New(fmt.Sprintf("%v", err)) }

	return &EchoResponse{Message: string(bs)}, nil

}
