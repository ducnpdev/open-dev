# rsa

## Đề Bài: 
- bài toán đặt ra là cần mã hoá thông tin giao tiếp giữa hệ thống server của HDBANK và trường đaị học HUTECH.
- serverA: server được đặt tại trường Đại Học Hutech.
- serverB: server được đặt tại ngân hàng HDBANK.
- thông tin cần mã hoá là mã số cvv trên thẻ ATM hoặc thẻ tín dụng.
## Step gọi
- HUTECH tạo ra cặp key rsa và gửi qua hệ thống hdbank public-key.
- Sau khi nhận public, HDBANK sẽ thực hiện query số cvv trên hệ thống, và mà hoá băng public được nhận, gửi lại về phía HUTECH.
- Sau khi nhận data được mã hoá từ HDBANK, HUTECH dựa vào private-key đã giữ để giải mã thông tin.

## Code:

### serverA:
- hàm tạo ra cặp key rsa:
```go
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
```
- chuyển public sang string:
```go
func ConvertPublicString(publicKey *rsa.PublicKey) string {
	derPublicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	strPublicKey := base64.StdEncoding.EncodeToString(derPublicKeyBytes)
	return strPublicKey
}
```

### serverB
- sau khi nhận public dạng string từ serverA, cần chuyển về *rsq.PublicKey
```go
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
```
- Mã hoá data:
```go
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
```

### serverA
- sau khi nhận data từ phía serverB, sẽ thực hiện giải mã:
```go
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
```

## Full Source Code:
- https://github.com/ducnpdev/open-dev/blob/master/usecase/rsa/gen.go