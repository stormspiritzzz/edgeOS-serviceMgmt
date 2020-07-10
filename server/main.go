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

const (
	port = ":50051"
)

// the map of AppUUID to ServUUID
var appID2ServIDMap map[string]string /*创建集合 */

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
	creatDir()

	s := grpc.NewServer()
	eservicemgmt.RegisterEdgeServiceMgmtServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Deploy(ctx context.Context, req *eservicemgmt.DeployReq) (*eservicemgmt.DeployResp, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method Deploy not implemented")
	var f *os.File
	defer f.Close()
	var err1 error
	filename := "./yamls/" + req.GetServUUID() + ".yaml"
	exist, err := PathExists(filename)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return &eservicemgmt.DeployResp{T: eservicemgmt.DeployResp_ERR,
			ErrStr: "检测文件是否存在出错:" + fmt.Sprintf("%s", err)}, nil
	}

	if exist { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_CREATE, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		fmt.Println("文件不存在,创建文件")
		f, err1 = os.Create(filename) //创建文件
	}

	if err != nil {
		return &eservicemgmt.DeployResp{T: eservicemgmt.DeployResp_ERR,
			ErrStr: "打开文件出错:" + fmt.Sprintf("%s", err)}, nil
	}

	_, err1 = f.Write([]byte(req.GetMajorManifest()))
	if err1 != nil {
		log.Println(err1.Error())
		return &eservicemgmt.DeployResp{T: eservicemgmt.DeployResp_ERR,
			ErrStr: "写入文件出错:" + fmt.Sprintf("%s", err)}, nil
	}

	cmd := fmt.Sprintf("docker-compose -f ./yamls/%s.yml up", req.GetServUUID())
	out := string(Cmd(cmd, true))
	fmt.Println(out)

	appID2ServIDMap[req.GetAppUUID()] = req.GetServUUID()

	return &eservicemgmt.DeployResp{T: eservicemgmt.DeployResp_OK}, nil
}
func (s *server) Destroy(ctx context.Context, req *eservicemgmt.DestroyReq) (*eservicemgmt.Empty, error) {
	cmd := fmt.Sprintf("docker-compose -f ./yamls/%s.yml down", req.GetServUUID())
	out := string(Cmd(cmd, true))
	fmt.Println(out)

	var f *os.File
	defer f.Close()

	filename := "./yamls/" + req.GetServUUID() + ".yaml"
	exist, err := PathExists(filename)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return &eservicemgmt.Empty{}, status.Errorf(codes.OK, fmt.Sprintf("error:%v", err))
	}

	if exist { //如果文件存在
		err = os.Remove(filename)
		if err != nil {
			return &eservicemgmt.Empty{}, status.Errorf(codes.OK, fmt.Sprintf("error:%v", err))
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

	filename := "./yamls/" + appID2ServIDMap[req.GetAppUUID()] + ".yaml"

	content, err := readAllIntoMemory(filename)
	if err != nil {
		log.Fatal(err)
		return &eservicemgmt.DiscoverResp{}, status.Errorf(codes.OK, fmt.Sprintf("error:%v", err))
	}
	fmt.Printf("%s\n", content)
	return &eservicemgmt.DiscoverResp{Manifest: string(content)}, nil
}

// SayHello implements helloworld.GreeterServer
// func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
// 	fmt.Println("######### get client request name :" + in.Name)
// 	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
// }

// test whether the dir or file exists
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

func creatDir() {
	_dir := "./yamls"
	exist, err := PathExists(_dir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}

	if exist {
		fmt.Printf("has dir![%v]\n", _dir)
	} else {
		fmt.Printf("no dir![%v]\n", _dir)
		// 创建文件夹
		err := os.Mkdir(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
}

func Cmd(cmd string, shell bool) []byte {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	}
}

// * 整个文件读到内存，适用于文件较小的情况
func readAllIntoMemory(filename string) (content []byte, err error) {
	fp, err := os.Open(filename) // 获取文件指针
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	fileInfo, err := fp.Stat()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = fp.Read(buffer) // 文件内容读取到buffer中
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
