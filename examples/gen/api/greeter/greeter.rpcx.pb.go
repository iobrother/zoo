// Code generated by protoc-gen-rpcx. DO NOT EDIT.
// versions:
// - protoc-gen-rpcx v0.3.0
// - protoc          (unknown)
// source: greeter/greeter.proto

package greeter

import (
	context "context"
	client "github.com/smallnest/rpcx/client"
	protocol "github.com/smallnest/rpcx/protocol"
	server "github.com/smallnest/rpcx/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = context.TODO
var _ = server.NewServer
var _ = client.NewClient
var _ = protocol.NewMessage

//================== interface skeleton ===================
type GreeterAble interface {
	// GreeterAble can be used for interface verification.

	// SayHello is server rpc method as defined
	SayHello(ctx context.Context, args *HelloRequest, reply *HelloReply) (err error)
}

//================== server skeleton ===================
type GreeterImpl struct{}

// ServeForGreeter starts a server only registers one service.
// You can register more services and only start one server.
// It blocks until the application exits.
func ServeForGreeter(addr string) error {
	s := server.NewServer()
	s.RegisterName("Greeter", new(GreeterImpl), "")
	return s.Serve("tcp", addr)
}

// SayHello is server rpc method as defined
func (s *GreeterImpl) SayHello(ctx context.Context, args *HelloRequest, reply *HelloReply) (err error) {
	// TODO: add business logics

	// TODO: setting return values
	*reply = HelloReply{}

	return nil
}

//================== client stub ===================
// Greeter is a client wrapped XClient.
type GreeterClient struct {
	xclient client.XClient
}

// NewGreeterClient wraps a XClient as GreeterClient.
// You can pass a shared XClient object created by NewXClientForGreeter.
func NewGreeterClient(xclient client.XClient) *GreeterClient {
	return &GreeterClient{xclient: xclient}
}

// NewXClientForGreeter creates a XClient.
// You can configure this client with more options such as etcd registry, serialize type, select algorithm and fail mode.
func NewXClientForGreeter(addr string) (client.XClient, error) {
	d, err := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	if err != nil {
		return nil, err
	}

	opt := client.DefaultOption
	opt.SerializeType = protocol.ProtoBuffer

	xclient := client.NewXClient("Greeter", client.Failtry, client.RoundRobin, d, opt)

	return xclient, nil
}

// SayHello is client rpc method as defined
func (c *GreeterClient) SayHello(ctx context.Context, args *HelloRequest) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	err = c.xclient.Call(ctx, "SayHello", args, reply)
	return reply, err
}

//================== oneclient stub ===================
// GreeterOneClient is a client wrapped oneClient.
type GreeterOneClient struct {
	serviceName string
	oneclient   *client.OneClient
}

// NewGreeterOneClient wraps a OneClient as GreeterOneClient.
// You can pass a shared OneClient object created by NewOneClientForGreeter.
func NewGreeterOneClient(oneclient *client.OneClient) *GreeterOneClient {
	return &GreeterOneClient{
		serviceName: "Greeter",
		oneclient:   oneclient,
	}
}

// ======================================================

// SayHello is client rpc method as defined
func (c *GreeterOneClient) SayHello(ctx context.Context, args *HelloRequest) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	err = c.oneclient.Call(ctx, c.serviceName, "SayHello", args, reply)
	return reply, err
}
