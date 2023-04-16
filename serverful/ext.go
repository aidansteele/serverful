package main

import (
	"context"
	"fmt"
	"github.com/aidansteele/serverful/extension"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()

	ln, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	go serveTunnel(ln)

	name, _ := os.Executable()
	name = path.Base(name)

	extclient := extension.NewClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	_, err = extclient.Register(ctx, name)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	go handleExtension(ctx, extclient)

	lambda.Start(func(ctx context.Context, input *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error) {
		return &events.LambdaFunctionURLResponse{
			StatusCode: http.StatusTemporaryRedirect,
			Headers: map[string]string{
				"Location": fmt.Sprintf("https://%s", ln.Addr()),
			},
		}, nil
	})
}

func serveTunnel(ln net.Listener) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	u, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	err := http.Serve(ln, httputil.NewSingleHostReverseProxy(u))
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

func handleExtension(ctx context.Context, extclient *extension.Client) {
	for {
		next, err := extclient.NextEvent(ctx)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		fmt.Printf("ext next(): %+v\n", next)

		deadline := time.UnixMilli(next.DeadlineMs)
		time.Sleep(time.Until(deadline) - time.Second)
	}

}
