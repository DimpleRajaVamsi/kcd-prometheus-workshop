package suggestions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DimpleRajaVamsi/kcd-prometheus-workshop/server/metrics"
	"github.com/DimpleRajaVamsi/kcd-prometheus-workshop/server/utils"
	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var carApi *car

type car struct {
	ctx context.Context
}

func GetCar(ctx context.Context) Suggestions {
	if carApi == nil {
		carCtx := utils.CreateContextWithLogger(ctx, "car")
		logger := utils.GetLogger(carCtx)
		logger.Info("Created new Car object")
		return &car{carCtx}
	}
	return beerApi
}

func (c *car) Suggest(w http.ResponseWriter, req *http.Request) {
	carName := fake.Car()
	logger := utils.GetLogger(c.ctx)
	logger.Info("Suggesting car", zap.Any("name", carName))
	metrics.GetApiInvokeCountCounter(c.ctx).With(prometheus.Labels{"api_name": "car", "success": "true"}).Add(1)
	fmt.Fprintf(w, "%#v", carName)
}
