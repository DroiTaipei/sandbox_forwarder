# Example
- The simple guide for using `github.com/DroiTaipei/droipkg/grpc/connpool`


## client.go

```
package main


import (
    "golang.org/x/net/context"
    "errors"
    "github.com/DroiTaipei/droipkg/grpc/connpool"
    "time"
    pb "github.com/DroiTaipei/droipkg/grpc/examples/helloworld"
    "log"

)

func main() {
    server := "127.0.0.1:8888"
    maxConn := 1
    pools := connpool.NewPools()

    err := pools.Connect(server, maxConn)
    if err != nil {
        log.Fatal( errors.New("can't connect to " + server))
    }

    p, err := pools.GetRoundRobin()
    if err != nil {
        log.Fatal( err)
    }

    conn, err := p.Get()
    if err != nil {
        log.Fatal( err)
    }

    c := pb.NewGreeterClient(conn)

    p.Put(conn)

    req := pb.HelloRequest{
        Name: "John Doe",
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    resp, err := c.SayHello(ctx, &req)
    if err != nil {
       log.Fatal( err)
    }
    log.Println(resp.GetName())

}

```

## server.go

```
package main

import (
    "google.golang.org/grpc"
    "golang.org/x/net/context"
    "net"
    pb "github.com/DroiTaipei/droipkg/grpc/examples/helloworld"
    "log"
)

type server struct {
    name string
}

func (s *server) Serve() error {

    lis, err := net.Listen("tcp", "127.0.0.1:8888")
    if err != nil {
        return err
    }
    println("Server Listen on", lis.Addr().String())
    gs := grpc.NewServer()
    pb.RegisterGreeterServer(gs, s)

    err = gs.Serve(lis)
    if err != nil {
        return err
    }

    return nil
}

func (s *server) SayHello(ctx context.Context, c *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{Name:"Jane Doe"}, nil
}

func main() {
    s := server{}
    log.Fatal(s.Serve())
}
```