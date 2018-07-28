package main

import (
	"fmt"
	"net"
	//"strconv"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	//"io"
	//"time"

	"github.com/gordonklaus/portaudio"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int32, 256)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	for {
		chk(stream.Read())
		buf := new(bytes.Buffer)
		chk(binary.Write(buf, binary.BigEndian, in))
		//fmt.Println(len(buf.Bytes()))
		//fmt.Println(".")
		data, err := gZipData(buf.Bytes())
		chk(err)

		_, err = Conn.Write(data)
		CheckError(err)
	}
	chk(stream.Stop())

}

//func gUnzipData(data []byte) (resData []byte, err error) {
//	b := bytes.NewBuffer(data)

//	var r io.Reader
//	r, err = gzip.NewReader(b)
//	if err != nil {
//		return
//	}

//	var resB bytes.Buffer
//	_, err = resB.ReadFrom(r)
//	if err != nil {
//		return
//	}

//	resData = resB.Bytes()

//	return
//}

func gZipData(data []byte) (compressedData []byte, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err = gz.Write(data)
	if err != nil {
		return
	}

	if err = gz.Flush(); err != nil {
		return
	}

	if err = gz.Close(); err != nil {
		return
	}

	compressedData = b.Bytes()

	return
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
