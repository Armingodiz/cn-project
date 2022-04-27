package main

import (
	"log"
	"net"
)

func main() {
	tcpAddress, err := net.ResolveTCPAddr("tcp", ":80")
	if err != nil {
		log.Println(err.Error())
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		log.Println(err.Error())
		return
	}
	connection, err := listener.Accept()
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = connection.Write([]byte("hi armin"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	buffer := make([]byte, 512)
	_, err = connection.Read(buffer[0:])
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(buffer))
}
