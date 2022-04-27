package main

import (
	"log"
	"net"
)

func main() {
	tcpAddress, err := net.ResolveTCPAddr("localhost", "8080")
	if err != nil {
		log.Println(err.Error())
		return
	}
	connection, err := net.DialTCP("localhost", nil, tcpAddress)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = connection.Write([]byte("armin"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	var buffer []byte
	_, err = connection.Read(buffer[0:])
	if err != nil {
		log.Println(err.Error())
		return
	}
}
