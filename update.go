package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ramantehlan/mateix/packages/e"
)

// ReadJSON a json file and return the byte data
func ReadJSON(file string) Config {
	fileData, err := os.Open(file)
	e.Check(err)
	defer fileData.Close()
	byteData, err := ioutil.ReadAll(fileData)
	e.Check(err)
	jsonData := ByteToJSON(byteData)
	return jsonData
}

// ByteToJSON Parse a json file
func ByteToJSON(byteData []byte) Config {
	var jsonData Config
	err := json.Unmarshal(byteData, &jsonData)
	e.Check(err)
	return jsonData
}

func connectToServer(addr string, port int) net.Conn {
	dest := addr + ":" + strconv.Itoa(port)
	fmt.Printf("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)

	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Some problem connecting.")
		} else {
			fmt.Println("Unknown error: " + err.Error())
		}
		os.Exit(1)
	}

	return conn
}

// Update is to update the path with the sync system
func Update(path string) {
	fmt.Println(path)
	conf := ReadJSON(path + "/.mateix/config.json")
	conn := connectToServer(conf.TargetIP, 1248)
	conn.Write([]byte(conf.TargetDir + "\n"))
	dataFile := path + "/data"
	dataFileHash := GetHash(dataFile)
	fmt.Println("LocalFile Hash: ", dataFileHash)
	incomingHash, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("incoming Hash: ", incomingHash)
	e.Check(err)
	if incomingHash != dataFileHash+"\n" {
		fmt.Println("Different Hash")
		dat, err := ioutil.ReadFile(dataFile)
		e.Check(err)
		text := string(dat)
		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		_, err = conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error writing to stream.")
		}
	} else {
		fmt.Println("Same Hash")
	}

}
