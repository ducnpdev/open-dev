# jwt

## parse token
- input token data: eyJraWQ******o10SMU_Zw1un3Q
- output:
  - header:
```json
{
  "kid": "YuyXoY",
  "alg": "RS256"
}
```
  - payload:
```json
{
  "iss": "https://appleid.apple.com",
  "aud": "com.tn.vpn.test",
  "exp": 16****606,
  "iat": 16****06,
  "sub": "000830.b5269f****6b552f.0130",
  "nonce": "fb5753a43f4****e14ed875d",
  "c_hash": "84eg****EiQ",
  "email": "q****v@pr****ay.appleid.com",
  "email_verified": "true",
  "is_private_email": "true",
  "auth_time": 16****06,
  "nonce_supported": true
}
```
- code:
```go
func ParserToken(tokenString string) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})
	// ... error handling
	fmt.Println(token, err)
	// do something with decoded claims
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
	}
}
```