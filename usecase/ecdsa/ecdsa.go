package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
)

var (
	DataSentitive string = "this is data"
)

func GenECDSA() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	var (
		privateKey *ecdsa.PrivateKey
		publicKey  *ecdsa.PublicKey
		err        error
	)
	// generate private key,
	privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey = &privateKey.PublicKey

	return privateKey, publicKey, err
}

// EncodePrivate, this is function convert ecdsa Private Key to string
func EncodePrivate(key *ecdsa.PrivateKey) (string, error) {
	var (
		keyByte []byte
		keyStr  string
		err     error
	)
	keyByte, err = x509.MarshalECPrivateKey(key)
	if err != nil {
		return keyStr, err
	}
	pemBlocl := &pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: keyByte,
	}
	pemByte := pem.EncodeToMemory(pemBlocl)
	return string(pemByte), nil
}

// EncodePublic, this is function convert ecdsa Public Key to string
func EncodePublic(key *ecdsa.PublicKey) (string, error) {
	var (
		keyByte []byte
		keyStr  string
		err     error
	)
	keyByte, err = x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return keyStr, err
	}
	pemBlocl := &pem.Block{
		Type:  "ECDSA PUBLIC KEY",
		Bytes: keyByte,
	}
	pemByte := pem.EncodeToMemory(pemBlocl)
	return string(pemByte), nil
}

// DecodePrivate, this is function convert string to *ecdsa Private Key
func DecodePrivate(keyStr string) (*ecdsa.PrivateKey, error) {
	blockPriv, _ := pem.Decode([]byte(keyStr))

	x509EncodedPriv := blockPriv.Bytes

	privateKey, err := x509.ParseECPrivateKey(x509EncodedPriv)

	return privateKey, err
}

// DecodePublic, this is function convert string to *ecdsa Public Key
func DecodePublic(keyStr string) (*ecdsa.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(keyStr))

	x509EncodedPub := blockPub.Bytes

	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey, err
}

func SignHash() []byte {
	h := md5.New()
	io.WriteString(h, DataSentitive)
	return h.Sum(nil)
}

// Signature, generate signature, propusal communication each service
func Signature(priv *ecdsa.PrivateKey) (string, error) {
	// sign data sign-hash
	signByte, err := ecdsa.SignASN1(rand.Reader, priv, SignHash())
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(signByte)
	return s, nil
}

// VerifySignature
func VerifySignature(publ *ecdsa.PublicKey, hash string) (bool, error) {
	hashByte, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}
	verify := ecdsa.VerifyASN1(publ, SignHash(), hashByte)
	return verify, nil
}

func MainECDSA() {
	priv, publi, err := GenECDSA()

	fmt.Println(priv, publi, err)

	signature, err := Signature(priv)

	fmt.Println(signature, err)

	verify, _ := VerifySignature(publi, signature)
	fmt.Println(verify)

}
