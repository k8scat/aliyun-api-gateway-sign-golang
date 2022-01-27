# aliyun-api-gateway-sign-golang

Unofficial golang sdk for [Aliyun API Gateway](https://www.aliyun.com/product/apigateway?spm=5176.19720258.J_8058803260.53.c89b2c4akF0F92).

## Install

```shell
go get github.com/k8scat/aliyun-api-gateway-sign-golang
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/k8scat/aliyun-api-gateway-sign-golang"
	"net/http"
)

func main() {
	appKey := ""
	appSecret := ""
	api := sign.NewAPIGateway(appKey, appSecret)

	rawurl := ""
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		panic(err)
	}

	err = api.Sign(req)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)
}
```

## Used by

[k8scat/Articli](https://github.com/k8scat/Articli)

## License

[MIT](./LICENSE)
