package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	serverip          string        = "kcd-server"
	port              int           = 8089
	delay             time.Duration = 120 * time.Second
	cocurrentRoutines int           = 3
	iterations        int           = 10
)

func main() {
	serverlUrl := fmt.Sprintf("http://%s:%d", serverip, port)
	logger, _ := zap.NewProduction()
	logger.Info("Starting the client with", zap.String("server", serverlUrl))

	apis := []string{fmt.Sprintf("%s/car", serverlUrl),
		fmt.Sprintf("%s/beer", serverlUrl),
		fmt.Sprintf("%s/delay", serverlUrl)}
	logger.Info("APIs", zap.Any("list", apis))

	for i := 0; i < iterations; i++ {
		var wg sync.WaitGroup
		for _, api := range apis {
			temp := strings.Split(api, "/")
			name := temp[len(temp)-1]
			for j := 0; j < cocurrentRoutines; j++ {
				wg.Add(1)
				go invokeApi(logger.Named(fmt.Sprintf("%s-%d-%d", name, i, j)), api, &wg)
			}
		}
		wg.Wait()
		time.Sleep(delay)
	}
}

func invokeApi(logger *zap.Logger, api string, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Get(api)
	if err != nil {
		logger.Error("Invoking", zap.String("api", api), zap.Error(err))
		return
	}
	body, _ := io.ReadAll(res.Body)
	logger.Info("Invoked ", zap.String("API", api), zap.ByteString("response", body))
}
