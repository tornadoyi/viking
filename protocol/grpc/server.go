package grpc

import (
	"errors"
	"fmt"
	"sync"
	_grpc "google.golang.org/grpc"
)

var (
	servers = make(map[string]*Server, 0)
	mutex sync.Mutex
)


type Server = _grpc.Server

func CreateServer(name string, opt ...ServerOption) (*Server, error) {
	defer mutex.Unlock()
	mutex.Lock()

	if _, ok := servers[name]; ok { return nil, errors.New(fmt.Sprintf("Repteated server %v", name))}

	server := _grpc.NewServer(opt...)
	servers[name] = server

	return server, nil
}


func GetServer(name string) (*Server, bool) {
	defer mutex.Unlock()
	mutex.Lock()
	server, ok := servers[name]
	return server, ok
}

func RemoveServer(name string) {
	defer mutex.Unlock()
	mutex.Lock()
	if _, ok := servers[name]; !ok { return }
	delete(servers, name)
}

