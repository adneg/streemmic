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
	"time"

	"github.com/gordonklaus/portaudio"
)

func chk(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {

	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	chk(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	chk(err)
	defer ServerConn.Close()

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int32, 256)
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), &out)
	chk(err)
	defer stream.Close()
	buf := make([]byte, 44100)
	chk(stream.Start())
	for {
		bufReader := new(bytes.Buffer)
		ServerConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, _, err := ServerConn.ReadFromUDP(buf)
		data, err := gUnzipData(buf[0:n])
		chk(err)
		bufReader.Write(data)
		err = binary.Read(bufReader, binary.BigEndian, out)

		chk(err)
		stream.Write()

	}

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
