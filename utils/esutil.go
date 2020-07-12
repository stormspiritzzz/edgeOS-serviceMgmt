package esutil

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

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

func Cmd(cmd string, shell bool) ([]byte, error) {
	log.Println("cmd string:", cmd)
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
func ReadAllIntoMemory(filename string) (content []byte, err error) {
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

func GetHostIp() (string, error) {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer conn.Close()
	string_slice := strings.Split(conn.LocalAddr().String(), ":")

	return string_slice[0], nil
}
