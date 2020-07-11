package main

import (
	eservicemgmt "edgeOS/edgeService"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

const (
	port = ":50051"
)

// the map of AppUUID to ServUUID
var appID2ServIDMap map[string]string

// type Empty struct{}

// server is used to implement eservicemgmt.EdgeServiceMgmtServer.
type server struct {
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	appID2ServIDMap = make(map[string]string)
	createYamlDir()

	s := grpc.NewServer()
	eservicemgmt.RegisterEdgeServiceMgmtServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Deploy(ctx context.Context, req *eservicemgmt.DeployReq) (*eservicemgmt.DeployResp, error) {
	var f *os.File
	defer f.Close()
	filename := "./yamls/" + req.GetServUUID() + ".yml"
	exist, err := PathExists(filename)
	if err != nil {
		log.Printf("get dir error![%v]\n", err)
		return &eservicemgmt.DeployResp{T: eservicemgmt.DeployResp_ERR,
			ErrStr: "check file whether exist error:" + fmt.Sprintf("%s", err)}, nil
	}

	if exist {
		f, err = os.OpenFile(filename, os.O_CREATE, 0666) //open file
		log.Println("file exist")
	} else {
		log.Println("file not exist,create it")
		f, err = os.Create(filename) // create file
	}

	if err != nil {
		return &eservicemgmt.DeployResp{T: eservicemgmt.DeployResp_ERR,
			ErrStr: "open file eeror:" + fmt.Sprintf("%s", err)}, nil
	}

	_, err = f.Write([]byte(req.GetMajorManifest()))
	if err != nil {
		log.Println(err.Error())
		return &eservicemgmt.DeployResp{T: eservicemgmt.DeployResp_ERR,
			ErrStr: "write file error:" + fmt.Sprintf("%s", err)}, nil
	}

	// cmd := fmt.Sprintf("docker-compose -f ./yamls/%s.yml up", req.GetServUUID())
	// out, err1 := Cmd(cmd, true)
	// if err1 != nil {
	// 	return &eservicemgmt.DeployResp{T: eservicemgmt.DeployResp_ERR,
	// 		ErrStr: "cmd error:" + fmt.Sprintf("%s", err)}, nil
	// }
	// outStr := string(out)
	// log.Println(outStr)

	appID2ServIDMap[req.GetAppUUID()] = req.GetServUUID()

	return &eservicemgmt.DeployResp{T: eservicemgmt.DeployResp_OK}, nil
}
func (s *server) Destroy(ctx context.Context, req *eservicemgmt.DestroyReq) (*eservicemgmt.Empty, error) {
	// cmd := fmt.Sprintf("docker-compose -f ./yamls/%s.yml down", req.GetServUUID())
	// out, err1 := Cmd(cmd, true)
	// if err1 != nil {
	// 	return &eservicemgmt.Empty{}, status.Errorf(codes.OK, fmt.Sprintf("error:%v", err1))
	// }
	// outStr := string(out)
	// log.Println(outStr)

	var f *os.File
	defer f.Close()

	filename := "./yamls/" + req.GetServUUID() + ".yml"
	exist, err := PathExists(filename)
	if err != nil {
		log.Printf("get dir error![%v]\n", err)
		return &eservicemgmt.Empty{}, status.Errorf(codes.OK, fmt.Sprintf("error:%v", err))
	}

	if exist {
		err = os.Remove(filename)
		if err != nil {
			return &eservicemgmt.Empty{}, err)
		}
	}

	for app := range appID2ServIDMap {
		if appID2ServIDMap[app] == req.GetServUUID() {
			delete(appID2ServIDMap, app)
			break
		}
	}

	return &eservicemgmt.Empty{}, nil
}

func (s *server) Discover(ctx context.Context, req *eservicemgmt.DiscoverReq) (*eservicemgmt.DiscoverResp, error) {
	var f *os.File
	defer f.Close()

	filename := "./yamls/" + appID2ServIDMap[req.GetAppUUID()] + ".yml"

	content, err := readAllIntoMemory(filename)
	if err != nil {
		log.Fatal(err)
		return &eservicemgmt.DiscoverResp{}, status.Errorf(codes.OK, fmt.Sprintf("error:%v", err))
	}
	log.Printf("%s\n", content)
	return &eservicemgmt.DiscoverResp{Manifest: string(content)}, nil
}

// check whether the dir or file exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// create the folder to contain the yamls
func createYamlDir() {
	_dir := "./yamls"
	exist, err := PathExists(_dir)
	if err != nil {
		log.Printf("create yaml dir error![%v]\n", err)
		return
	}

	if exist {
		log.Printf("has dir![%v]\n", _dir)
	} else {
		// create dir
		err := os.Mkdir(_dir, os.ModePerm)
		if err != nil {
			log.Printf("mkdir failed![%v]\n", err)
		} else {
			log.Printf("mkdir success!\n")
		}
	}
}

func Cmd(cmd string, shell bool) ([]byte, error) {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			panic("some error found")
			return out, err
		}
		return out, nil
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			panic("some error found")
			return out, err
		}
		return out, nil
	}
}

// read the whole file into memory
func readAllIntoMemory(filename string) (content []byte, err error) {
	fp, err := os.Open(filename) // get the file pointer
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	fileInfo, err := fp.Stat()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = fp.Read(buffer) // put the content of the file into the buffer
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
