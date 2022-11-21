package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

var cvv = "434" // cvv number

// tạo ra cặp key rsa
func GenerateRSA(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	var (
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
	)
	if bits == 0 {
		bits = 2048
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	publicKey = &privateKey.PublicKey
	return privateKey, publicKey, err
}

// chuyển đổi public-key sang string,
// vì khi giao tiếp 2 service đa phần sẽ là json-string
func ConvertPublicString(publicKey *rsa.PublicKey) string {
	derPublicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	strPublicKey := base64.StdEncoding.EncodeToString(derPublicKeyBytes)
	return strPublicKey
}

// chuyển đổi từ string sang public-key
// sau khi nhận request có public-key, cần chuyển pulbic-key dạng string thành *rsa.PublicKey
func ConvertStringToPublic(publicKey string) (*rsa.PublicKey, error) {
	var (
		err        error
		bytePublic []byte
		rsaPublic  *rsa.PublicKey
	)
	bytePublic, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return rsaPublic, err
	}
	rsaPublic, err = x509.ParsePKCS1PublicKey(bytePublic)
	if err != nil {
		return rsaPublic, err
	}
	return rsaPublic, nil
}

// mã hoá data thông qua puplic-key
func EncryptData(public *rsa.PublicKey, dataSentitive string) (string, error) {
	var (
		dataEncrypt string
		err         error
	)
	encryptKey, err := rsa.EncryptPKCS1v15(rand.Reader, public, []byte(dataSentitive))
	if err != nil {
		return dataEncrypt, err
	}
	dataEncrypt = base64.StdEncoding.EncodeToString(encryptKey)
	return dataEncrypt, nil
}

// giải mã data thông qua private-key
func DecryptData(private *rsa.PrivateKey, data string) (string, error) {
	var (
		dataDecrypt []byte
		err         error
	)

	bytePublic, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	dataDecrypt, err = rsa.DecryptPKCS1v15(rand.Reader, private, bytePublic)
	if err != nil {
		return "", err
	}
	return string(dataDecrypt), nil
}

// main function
// code test example

func MainRsa() {
	// prepare key
	privKeyOri, publKeyOri, err := GenerateRSA(0)
	if err != nil {
		return
	}
	strPublKey := ConvertPublicString(publKeyOri)

	// encrypt proces
	bytePublKey, err := ConvertStringToPublic(strPublKey)
	if err != nil {
		fmt.Println("error decode rsa public key")
		return
	}

	publKey, err := EncryptData(bytePublKey, cvv)
	if err != nil {
		fmt.Println("error parser rsa public key")
		return
	}

	// decrypt process
	data, err := DecryptData(privKeyOri, publKey)
	if err != nil {
		fmt.Println("error decode rsa public key")
		return
	}
	if data == cvv {
		fmt.Println("ok")
	}
}
