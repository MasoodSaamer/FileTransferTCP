//A simple test function to ensure serialization, deserilization and checksum functions work

package pkg

import (
	"reflect"
	"testing"
)

func TestPacketSerialization(t *testing.T) {
	// first we create an original packet for our testing purpose
	originalPacket := &Packet{
		PacketType:     0x01,
		SequenceNumber: 1,
		Payload:        []byte("Test Payload"),
	}

	// Calculating the payload length and checksum of the packet
	originalPacket.PayloadLength = uint16(len(originalPacket.Payload)) // length of the payload bytes
	originalPacket.CalculateChecksum()

	// Serializing the packet into a byte slice (serialized)
	serialized, err := originalPacket.ToBytes()
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	//Deserializing the packet
	deserializedPacket, err := FromBytes(serialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	// now to compare the original and deserialized packets
	if !reflect.DeepEqual(originalPacket, deserializedPacket) { // Ensuring the packets are deeply equal (all their fields are the same value) using reflect
		t.Fatalf("Deserialized packet does not match original: got %+v, want %+v", deserializedPacket, originalPacket)
	}

	//Verify the checksum
	if !deserializedPacket.VerifyChecksum() {
		t.Fatalf("Checksum verification failed")
	}
}
