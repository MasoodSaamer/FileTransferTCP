// This file will contain functions to handle encryption and decryption

package pkg

import (
	"crypto/aes"    // Provides the implementation of AES encryption.
	"crypto/cipher" // Provides interfaces for stream ciphers, including CFB.
	"crypto/rand"   // Provides the ability to read cryptographically secure random numbers.
	"fmt"           // Provides formatted I/O functions.
	"io"            // Provides basic interfaces for I/O.
)

// Encrypt encrypts the given data using AES encryption with the provided key.
func Encrypt(data []byte, key []byte) ([]byte, error) {
	// Create a new AES cipher block with the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Generate a random initialization vector (IV)
	ciphertext := make([]byte, aes.BlockSize+len(data))     // AES block size is 16 bytes
	iv := ciphertext[:aes.BlockSize]                        // Extracts the first 16 bytes of ciphertext slice to store in IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil { //Reads cryptographically secure random bytes from rand.Reader into the iv.
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Encrypt the data using CFB mode (Cipher Feedback mode)
	stream := cipher.NewCFBEncrypter(block, iv)           //Creates a new CFB (Cipher Feedback) encrypter stream using the AES cipher block and the IV.
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data) //Encrypts the plaintext 'data' using CFB and writes

	return ciphertext, nil // Return ciphertext slice which contains both the IV and the encrypted data
}

// Decrypt decrypts the given data using AES decryption with the provided key.
func Decrypt(data []byte, key []byte) ([]byte, error) { // data is the encrypted data(ciphetext). key is the decryption which must be the same as used for encryption
	// Create a new AES cipher block with the given key
	block, err := aes.NewCipher(key) // Same as encryption
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Extract the initialization vector (IV)
	if len(data) < aes.BlockSize { //if length of encrypted data is less than AES block size of 16 bytes
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := data[:aes.BlockSize]         // extract the first 16 bytes of encrypted data
	ciphertext := data[aes.BlockSize:] // extract the rest of the encrypted data as ciphertext to be decrypted

	// Decrypt the data using CFB mode
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext) //overwrite the ciphertext with decrypted data

	return ciphertext, nil // return the decrypted(now plaintext) ciphertext
}
