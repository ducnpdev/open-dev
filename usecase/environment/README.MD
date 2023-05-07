# environment
## load env yaml
- dựa vào biến `MODE` để xác định môi trường: `dev` `uat` `prod`. Thông qua func `getConfigName()`
- `AddConfigPath` là cấu hình đường dẫn của file env, ví dụ như bạn muốn đọc file yaml từ `/abc/config/dev.yaml`, trong khi source-code thì nằm ở `abc/source/`. Lúc này `AddConfigPath` thì sẽ truyền vào là `/abc/config/`
- `SetConfigName` filename của env.
- Vòng lặp `for` thì sẽ in ra tất cả env, hỗ trợ trong quá trình debug
## source-code
- get mode:
```go
func getConfigName() string {
	mode := "dev"
	switch os.Getenv("MODE") {
	case "uat":
		mode = "uat"
	case "prod":
		mode = "prod"
	}
	return mode
}
```
- LoadConfig từ file yaml:
```go
func LoadConfig() (*Config, error) {
	mode := getConfigName()
	fmt.Println("mode:", mode)
	cfg := &Config{}
	path := "***"
	vn := viper.New()
	vn.AddConfigPath(path)
	vn.SetConfigName(mode)
	vn.SetConfigType("yaml")
	vn.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	vn.AutomaticEnv()

	err := vn.ReadInConfig()
	if err != nil {
		log.Panic(err)
		return cfg, err
	}

	for _, key := range vn.AllKeys() {
		str := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		log.Println(key, str, vn.Get(key))
		vn.BindEnv(key, str)
	}

	err = vn.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}
```
## Demo
- link youtube: 
