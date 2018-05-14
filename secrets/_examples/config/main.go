package main

type myConfig struct {
	Environment string `json:"environment" yaml:"environment" env:"SERVICE_ENV"`
	Version     string `json:"version" yaml:"version" env:"SERVICE_VERSION"`

	Secret []byte `secret:"secret,base64"`
}

func main() {

}
