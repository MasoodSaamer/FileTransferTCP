# File Transfer Protocol Design by Saamer Masood

## Introduction
This project was created by me as a way to dive into the world of networks and this document outlines the requirements and design of a custom File Transfer Protocol. Its basically about how data can be safely transferred from one system to another over a network.

## Features
- **File Chunking:** Breaks files into smaller parts for manageable transmission.
- **Error Detection:** Uses checksums to verify data integrity.
- **Sequence Numbers:** Ensures packets are received in the correct order.
- **Acknowledgments (ACKs):** Confirms receipt of each packet.
- **Timeout and Retransmission:** Handles lost packets by retransmitting if no acknowledgment is received.
- **Optional Encryption:** Secures data using AES encryption.

## Packet Structure
| Field          | Size (bytes) | Description                                        |
|----------------|--------------|----------------------------------------------------|
| Packet Type    | 1            | Specifies packet type (Data = 0x01, ACK = 0x02)    |
| Sequence Number| 4            | Tracks the order of packets                        |
| Checksum       | 4            | Used for error detection                           |
| Payload Length | 2            | Length of the data payload                         |
| Reserved       | 9            | Reserved for future use                            |
| Payload        | Up to 1024   | Actual data chunk being transferred                |

## Communication Flow
1. **Client Initiates Connection:** Connects to the server via TCP/UDP.
2. **File Transfer Begins:** Client sends file metadata (name, size).
3. **Data Transmission:** 
   - Client sends data packets.
   - Server sends ACK for each packet.
   - Retransmit if no ACK is received.
4. **Error Handling:** Server requests retransmission if checksum fails.
5. **Completion Confirmation:** Server sends a final acknowledgment when all packets are received.

## Example Packets
- **Data Packet:** `[0x01 | 0x00000001 | 0xABCD1234 | 0x0400 | ...data...]`
- **Acknowledgment Packet:** `[0x02 | 0x00000001 | 0x00000000 | 0x0000 | ...empty...]`
