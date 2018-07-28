package main

import (
	//"bytes"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/gordonklaus/portaudio"
)

/* A Simple function to verify error */
func chk(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	chk(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	chk(err)
	defer ServerConn.Close()

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int32, 256)
	//out2 := make([]int32, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), &out)
	chk(err)
	defer stream.Close()
	buf := make([]byte, 44100)
	chk(stream.Start())
	for {
		bufReader := new(bytes.Buffer)
		n, _, err := ServerConn.ReadFromUDP(buf)

		data, err := gUnzipData(buf[0:n])
		//fmt.Println(len(data))

		chk(err)
		bufReader.Write(data)
		err = binary.Read(bufReader, binary.BigEndian, out)
		//fmt.Println(len(out))
		//fmt.Println(data)
		chk(err)
		stream.Write()
		//fmt.Println(len(out2))
		//binary.Write(buf, binary.LittleEndian, in)
	}
	fmt.Println("dupa")
	//	for {
	//		_, _, err := ServerConn.ReadFromUDP(buf)
	//		//		fmt.Println("Received ", len(buf[0:n]), " from ", addr)
	//		if err != nil {
	//			fmt.Println("Error: ", err)
	//		}

	//		stream.Write()

	//	}
}
func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}
