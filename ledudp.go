package lcv

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

// UDP Client
type udpC struct {
	// the host and port the client connects to
	host string
	port string
	// the udp endpoint address, contains IP and port information
	udpaddr *net.UDPAddr
	// the udp connection
	udpconn *net.UDPConn
}

// Generates a new udp client at localhost with the specified port number
func newUdpC(port string) *udpC {
	var client = &udpC{
		host: "127.0.0.1:",
		port: port,
	}
	return client
}

// Starts the udp client on the specified host and port
func (c *udpC) start() {
	var err error
	c.udpaddr, err = net.ResolveUDPAddr("udp4", c.host+c.port)
	c.udpconn, err = net.DialUDP("udp4", nil, c.udpaddr)
	chk(err)

	fmt.Printf("The UDP server is %s\n", c.udpconn.RemoteAddr().String())
}

// Sends a string message on the udp client
func (c *udpC) sendMsg(m string) {
	data := []byte(m + "\n")
	_, err := c.udpconn.Write(data) // write to buffer and then exit the client after notifying server of shutdown#
	// Find a way to close the connection
	if strings.TrimSpace(string(data)) == "STOP" {
		fmt.Println("Exiting UDP client!")
		return
	}
	chk(err)
}

// Closes the connection of the client to the specified host and port
func (c *udpC) closeConnection() {
	c.udpconn.Close()
}

// UDP Server used for testing
type udpS struct {
	// the port to run the server on
	port string
	// the udp endpoint address, contains IP and port information
	udpaddr *net.UDPAddr
	// the udp connection
	udpconn *net.UDPConn
}

// Generates a new server object
func NewUdpS() *udpS {
	var client = &udpS{
		port: "127.0.0.1:6969",
	}
	return client
}

// Runs the server to print all messages it receives from the client
func (s *udpS) Run() {
	var err error
	s.udpaddr, err = net.ResolveUDPAddr("udp4", s.port) // Gets the udp endpoint address
	chk(err)
	s.udpconn, err = net.ListenUDP("udp4", s.udpaddr) // Listen for a connection from a udp client
	chk(err)

	defer s.udpconn.Close()
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, _, err := s.udpconn.ReadFromUDP(buffer)
		// n is the number of bytes read from the buffer, we minus 1 to compensate for the newline which the client writes
		fmt.Print("-> ", string(buffer[0:n-1]), "\n")

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		chk(err)
	}
}
