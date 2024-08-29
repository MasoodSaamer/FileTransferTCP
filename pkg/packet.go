// This file will define the packet structure and provide functions for creating, serializing (converting to bytes), and deserializing (parsing from bytes) the packets.
package pkg

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

// representing the packet format
// Packet represents the structure of a data packet
type Packet struct {
	PacketType     byte   //1 byte: Type of packet (eg Data = 0x01, ACK = 0x02)
	SequenceNumber uint32 //4 bytes: packet sequence number
	Checksum       uint32 // 4 bytes: checksume for error detection
	PayloadLength  uint16 // 2 bytes: Lenght of payload
	Payload        []byte // Upto 1024 bytes: Actual data being transferred
}

// now to convert the packet into bytes(serialization)
// ToBytes serializes the Packet struct into a byte slice. It is a method type function with a reciever being p *Packet
func (p *Packet) ToBytes() ([]byte, error) { //takes packet itself as a parameter (p) and returns the sliced byte array and an error if their is one
	buf := new(bytes.Buffer) //buffer to accumalte the serialized byte

	//Writing the packet type (1 byte)
	err := binary.Write(buf, binary.BigEndian, p.PacketType) //binary.Write writes a value to buffer buf in a specific byte order. binary.BigEndian specifies the byte order where the most significant byte is stored first. p.PacketType is the value to be written to the buffer
	if err != nil {
		return nil, err //error catching
	}

	// Writing the SequenceNumber(4 bytes). used to keep track of the order of packets to ensure they are reassembled correctly by the receiver.
	err = binary.Write(buf, binary.BigEndian, p.SequenceNumber)
	if err != nil {
		return nil, err
	}

	// Writing Checksum (4 bytes). value used for error detection to ensure data integrity during transmission.
	err = binary.Write(buf, binary.BigEndian, p.Checksum)
	if err != nil {
		return nil, err
	}

	// Writing PayloadLenght (2 byte). PayloadLength indicates the length of the actual data (payload) being transmitted in the packet.
	err = binary.Write(buf, binary.BigEndian, p.PayloadLength)
	if err != nil {
		return nil, err
	}

	// Writing the Payload (depends on the file sent). The payload is a byte slice that contains the actual data being sent (e.g., a chunk of a file).
	err = binary.Write(buf, binary.BigEndian, p.Payload)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil //called to convert the buffer into a byte slice and return it
}

// Now to convert the byte back into a packet (Deserialization)
func FromBytes(data []byte) (*Packet, error) { // Regular function with no method or reciever
	buf := bytes.NewReader(data)
	packet := &Packet{} // Pointer to the Packet Structure. & specifies the address and Packet{} is a the instance

	// Read PacketType (1 byte)
	err := binary.Read(buf, binary.BigEndian, &packet.PacketType) // Using & since read function requires pointer to the PacketType
	if err != nil {
		return nil, err
	}

	// Read the SequenceNumber (4 bytes)
	err = binary.Read(buf, binary.BigEndian, &packet.SequenceNumber)
	if err != nil {
		return nil, err
	}

	// Read Checksum (4 bytes)
	err = binary.Read(buf, binary.BigEndian, &packet.Checksum)
	if err != nil {
		return nil, err
	}

	// Read PayloadLength (2 bytes)
	err = binary.Read(buf, binary.BigEndian, &packet.PayloadLength)
	if err != nil {
		return nil, err
	}

	// Read payload (variable size)
	packet.Payload = make([]byte, packet.PayloadLength) // Allocating memory for the Payload slice since in byte slices arent automatically allocated memories. Length is Payloadlength. This is required for variable size data
	err = binary.Read(buf, binary.BigEndian, &packet.Payload)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

// Now we add a Checksume calculation and verification method for error detection

// Calculating the Checksum using CRC32 checksum.
func (p *Packet) CalculateChecksum() {
	p.Checksum = crc32.ChecksumIEEE(p.Payload) // uses Go's hash/crc32 package to calculate CRC32 Checksum
}

// Verifying the checksume matches the sender's checksum
func (p *Packet) VerifyChecksum() bool {
	return p.Checksum == crc32.ChecksumIEEE(p.Payload) // p.Checksum is the stored checksum. crc32.. is the newly calculated checksum
}
