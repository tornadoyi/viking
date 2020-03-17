package grpc


import (
	"errors"
	"fmt"
	_grpc "google.golang.org/grpc"
	"sync"
)

var (
	clients = make(map[string]*Client)
	clientMutex sync.Mutex
)


type Client = _grpc.ClientConn



func CreateClient(name string, address string,  opt ...DialOption) (*Client, error) {
	defer clientMutex.Unlock()
	clientMutex.Lock()

	if _, ok := clients[name]; ok { return nil, errors.New(fmt.Sprintf("Repteated client %v", name))}
	client, err := _grpc.Dial(address, opt...)
	if err != nil { return nil, err}
	clients[name] = client

	return client, nil
}


func GetClinet(name string) (*Client, bool) {
	defer clientMutex.Unlock()
	clientMutex.Lock()
	client, ok := clients[name]
	return client, ok
}

func RemoveClient(name string) {
	defer clientMutex.Unlock()
	clientMutex.Lock()
	if _, ok := clients[name]; !ok { return }
	delete(clients, name)
}

