package main

import (
	eservicemgmt "edgeOS/edgeService"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// req := eservicemgmt.DeployReq{AppUUID: uuid.NewV4().String(),
	// 	ServUUID:      uuid.NewV4().String(),
	// 	MajorManifest: "testtest",
	// 	MinorManifest: "testtest",
	// }
	req := eservicemgmt.DeployReq{AppUUID: "2b104870-b306-4e17-9741-b33af22e5991",
		ServUUID:      "59f55ab8-2210-4914-a2d5-4f102e6ce140",
		MajorManifest: "testtest",
		MinorManifest: "testtest",
	}
	Deploy(&req)

	// req := eservicemgmt.DestroyReq{ServUUID: "59f55ab8-2210-4914-a2d5-4f102e6ce140"}
	// Destroy(&req)

	// req := eservicemgmt.DiscoverReq{AppUUID: "2b104870-b306-4e17-9741-b33af22e5991"}
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
		log.Fatalf("deploy error: %v", err)
	}
	log.Printf("####### get server response: %s", r.ErrStr)
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
		log.Fatalf("destroy error: %v", err)
	}
	log.Printf("####### get server response: %s", "")
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
		log.Fatalf("discover error: %v", err)
	}
	log.Printf("####### get server response: %s", r.String())
}
