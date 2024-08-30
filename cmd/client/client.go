// Handling client side connection
package main

import (
	"bufio"
	"file_transfer_protocol/pkg"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// Defining the constants for the client config
const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

func main() {
	// Connect to the server
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()
	fmt.Printf("Connected to server %s:%s\n", SERVER_HOST, SERVER_PORT)

	// Get the user input for the file path to be sent
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the path of the file to send: ")

	// Read the input until the newline character is encountered
	filePath, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	// Trim any newline and extra whitespace characters from the input
	filePath = strings.TrimSpace(filePath)

	// Verify if the path is correct (for debugging purposes)
	fmt.Printf("File path entered: %s\n", filePath)

	// Send the file to the server
	err = sendFile(filePath, conn)
	if err != nil {
		log.Fatalf("Error sending file: %v", err)
	}

	fmt.Println("File transfer complete.")
}

// Now to create the sendFile function which splits the file into chunks, creates data packets and sends them to the server
func sendFile(filePath string, conn net.Conn) error { // if an error occurs, it will return an error or else it will return nil
	key := []byte("mysecretbyteskey") // Ensure the key is 16 bytes long

	// Open the file
	file, err := os.Open(filePath) // Opens file in read-only mode
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 1024) // Create a byte slice of 1 KB chunk size. This buffer will be used to read chunks of the file to be sent to the server
	sequenceNumber := uint32(1)  // Initialize the sequence number to 1. keep track of order of packets sent to the server

	for { // Loop that continously reads all the chunks of file until they are read and sent
		n, err := file.Read(buffer) // n is the number of bytes actually read into the buffer
		if err != nil && err.Error() != "EOF" {
			return fmt.Errorf("failed to read file: %v", err)
		}

		// Check if the end of file has been reached
		if n == 0 { // If the number of bytes read is 0, it measn the EOF has been reached
			break
		}

		// Encrypt the data before creating a packet
		encryptedData, err := pkg.Encrypt(buffer[:n], key)
		if err != nil {
			return fmt.Errorf("failed to encrypt data: %v", err)
		}

		// Create a data packet
		packet := &pkg.Packet{ // Creates a new packet structure to repesent the data packet to be sent to the server
			PacketType:     0x01, // Data Packet
			SequenceNumber: sequenceNumber,
			Payload:        encryptedData, // Actual data read from the file. aka actual chunk of file daata
		}
		packet.PayloadLength = uint16(len(packet.Payload)) // Compute the length of the payload
		packet.CalculateChecksum()

		// Serialize the packet
		data, err := packet.ToBytes()
		if err != nil {
			return fmt.Errorf("failed to serialize packet: %v", err)
		}

		// Send the serialized packet to server
		_, err = conn.Write(data) //The return value (_) is ignored because we are not interested in the number of bytes sent; only that the operation succeeds.
		if err != nil {
			return fmt.Errorf("failed to send packet: %v", err)
		}

		// Confirmation
		fmt.Printf("Sent packet %d to server\n", sequenceNumber)

		// Wait for ACK from the server
		if !waitForAck(conn, sequenceNumber) {
			return fmt.Errorf("failed to receive ACK for packet %d", sequenceNumber)
		}

		// Increment sequence number for the next packet
		sequenceNumber++

	}

	return nil

}

// waitForAck handles Acknoledgements and Retransmissions if necessary
func waitForAck(conn net.Conn, expectedSeqNum uint32) bool {
	buffer := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Read timeout

	for {
		// Read data from the server
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from server: %v", err)
			return false
		}

		// Deserialize the recieved data packet
		ackPacket, err := pkg.FromBytes(buffer[:n])
		if err != nil {
			log.Printf("Error in deserializing ACK packet: %v", err)
			return false
		}

		// Check if its an ACK packet with the expected sequence number
		if ackPacket.PacketType == 0x02 && ackPacket.SequenceNumber == expectedSeqNum {
			fmt.Printf("Received ACK for packet %d\n", expectedSeqNum)
			return true
		}

		// Continue reading if the packet is not the expected packet
	}
}
