package suggestions

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/DimpleRajaVamsi/kcd-prometheus-workshop/server/metrics"
	"github.com/DimpleRajaVamsi/kcd-prometheus-workshop/server/utils"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var delayApi *delay

const (
	minDelayInSec int = 0
	maxDelayInSec int = 10
)

type delay struct {
	ctx context.Context
}

func GetDelay(ctx context.Context) Suggestions {
	if delayApi == nil {
		delayCtx := utils.CreateContextWithLogger(ctx, "delay")
		logger := utils.GetLogger(delayCtx)
		logger.Info("Created new Delay object")
		return &delay{delayCtx}
	}
	return delayApi
}

func (d *delay) Suggest(w http.ResponseWriter, req *http.Request) {
	rand.NewSource(time.Now().UnixNano())
	delay := rand.Intn(maxDelayInSec-minDelayInSec+1) + minDelayInSec
	logger := utils.GetLogger(d.ctx)
	logger.Info("Suggesting delay", zap.Any("value", time.Duration(delay)*time.Second))
	time.Sleep(time.Duration(delay) * time.Second)
	metrics.GetApiInvokeCountCounter(d.ctx).With(prometheus.Labels{"api_name": "delay", "success": "true"}).Add(1)
	metrics.GetDelayHistogram(d.ctx).Observe(float64(delay))
	metrics.GetDelaySummary(d.ctx).Observe(float64(delay))
	fmt.Fprintf(w, "%d", delay)
}
