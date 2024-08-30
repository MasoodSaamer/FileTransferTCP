// Server implementation
// Starting by setting up a TCP Server that listens on a specific port

package main //entry point of the Go program

import (
	"file_transfer_protocol/pkg"
	"fmt"
	"io"
	"log"
	"net"
)

const ( // define the constants for the server
	SERVER_HOST = "localhost" // local machine
	SERVER_PORT = "8080"      //port on which server will listen for incoming connections
	SERVER_TYPE = "tcp"       // TCP = transmission control protocol
)

func main() {
	// Start the TCP Server
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer server.Close()
	fmt.Printf("Listening on %s:%s...\n", SERVER_HOST, SERVER_PORT)

	// Handle incoming client connections
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("Error accepting connections: %v", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()                                                  // Ensure the connection is properly closed when the handleClient function exits out
	fmt.Printf("Connected to client: %s\n", conn.RemoteAddr().String()) // Show the client's remote address

	// Buffer to store incoming data
	buffer := make([]byte, 4096) // Create a bytle slice with a capacity of 4096 to store incoming data from the client

	for { // Read data from client until connection is close
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF { // If client disconnects
				log.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())
				return
			}
			log.Printf("Error reading from client: %v", err)
			return
		}

		// Process the incoming data (deserialize packet)
		packet, err := pkg.FromBytes(buffer[:n]) // Slice the buffer to actual number of bytes read
		if err != nil {
			log.Printf("Error deserializing packet: %v", err)
			continue
		}

		// Verify the packet checksum
		if !packet.VerifyChecksum() {
			log.Printf("Invalid checksum for packet from client %s", conn.RemoteAddr().String())
			continue
		}

		fmt.Printf("Received packet from %s: %+v\n", conn.RemoteAddr().String(), packet)

		// Handle received packet (eg, ACK, retransmit)
		processPacket(packet, conn)
	}
}

// Function to process incoming packets
func processPacket(packet *pkg.Packet, conn net.Conn) { //net.Conn is the client connections
	switch packet.PacketType { // Switch case to handle the different packet types
	case 0x01: // Data packet
		fmt.Printf("Processing data packet with sequence number %d\n", packet.SequenceNumber)

		// Send an ACK(Acknowledgement) for received data packet
		ackPacket := &pkg.Packet{
			PacketType:     0x02, // ACK Packet
			SequenceNumber: packet.SequenceNumber,
			Payload:        nil, // ACK packet does not need to carry a payload
		}
		ackPacket.CalculateChecksum() // needed to ensure integrity
		sendPacket(ackPacket, conn)   // Calls sendPacket func to send ACK packet back to client over the existing connection of conn

	case 0x02: // ACK Packet
		fmt.Printf("Received ACK for packet with sequence number: %d\n", packet.SequenceNumber)

	default: // For handling unexpected or unrecognized packet types
		log.Printf("Unknown packet type: %d\n", packet.PacketType)
	}
}

// function to send ACK packet back to client
func sendPacket(packet *pkg.Packet, conn net.Conn) {
	//Serialize the packet
	data, err := packet.ToBytes()
	if err != nil {
		log.Printf("Error serializing packet: %v", err)
		return
	}

	_, err = conn.Write(data) // Send the serialized byte slice(data) to the client over the exisiting connection of conn
	if err != nil {
		log.Printf("Error sending packet: %v", err)
	}

}
