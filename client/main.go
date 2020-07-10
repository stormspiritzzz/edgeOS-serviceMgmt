package main

import (
	eservicemgmt "edgeOS/edgeService"
	"log"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	req := eservicemgmt.DeployReq{AppUUID: uuid.NewV4().String(),
		ServUUID:      uuid.NewV4().String(),
		MajorManifest: "testtest",
		MinorManifest: "testtest",
	}
	Deploy(&req)

	// req := eservicemgmt.DestroyReq{ServUUID: uuid.NewV4().String()}
	// Destroy(&req)

	// req := eservicemgmt.DiscoverReq{AppUUID: uuid.NewV4().String()}
	// Discover(&req)

}

func Deploy(req *eservicemgmt.DeployReq) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := eservicemgmt.NewEdgeServiceMgmtClient(conn)

	r, err := c.Deploy(context.Background(), req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("####### get server Greeting response: %s", r.ErrStr)
}

func Destroy(req *eservicemgmt.DestroyReq) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := eservicemgmt.NewEdgeServiceMgmtClient(conn)

	_, err1 := c.Destroy(context.Background(), req)
	if err1 != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("####### get server Greeting response: %s", "")
}

func Discover(req *eservicemgmt.DiscoverReq) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := eservicemgmt.NewEdgeServiceMgmtClient(conn)

	r, err := c.Discover(context.Background(), req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("####### get server Greeting response: %s", r.String())
}
