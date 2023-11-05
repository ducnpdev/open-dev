# convert file key and pem to base64
- code example
```go
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func main() {
	// Replace these with the actual file paths
	certFilePath := "file.pem"
	keyFilePath := "file.key"

	// Read the PEM certificate file
	certPEM, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		fmt.Printf("Failed to read certificate file: %v\n", err)
		return
	}

	// Read the PEM private key file
	keyPEM, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		fmt.Printf("Failed to read private key file: %v\n", err)
		return
	}

	// Encode the certificate and private key to base64
	certBase64 := base64.StdEncoding.EncodeToString(certPEM)
	keyBase64 := base64.StdEncoding.EncodeToString(keyPEM)

	// Print the base64-encoded certificate and private key (for demonstration)
	fmt.Println("Base64-encoded certificate:")
	fmt.Println(certBase64)
	fmt.Println("\nBase64-encoded private key:")
	fmt.Println(keyBase64)

	// You can save the base64-encoded data to files or use it as needed
	// For example, to save the base64-encoded certificate to a file:
	err = ioutil.WriteFile("certificate_base64.txt", []byte(certBase64), 0644)
	if err != nil {
		fmt.Printf("Failed to write base64-encoded certificate to file: %v\n", err)
		return
	}

	// Similarly, you can save the base64-encoded private key to a file.
}
```