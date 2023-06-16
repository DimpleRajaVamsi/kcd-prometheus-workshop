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

var beerApi *beer

type beer struct {
	ctx context.Context
}

func GetBeer(ctx context.Context) Suggestions {
	if beerApi == nil {
		beerCtx := utils.CreateContextWithLogger(ctx, "beer")
		logger := utils.GetLogger(beerCtx)
		logger.Info("Created new Beer object")
		return &beer{beerCtx}
	}
	return beerApi
}

func (b *beer) Suggest(w http.ResponseWriter, req *http.Request) {
	beerName := fake.BeerName()
	logger := utils.GetLogger(b.ctx)
	httpCode := fake.HTTPStatusCodeSimple()

	if httpCode == 200 {
		logger.Info("Suggesting beer", zap.String("name", beerName))
		fmt.Fprintf(w, "%s", beerName)
		metrics.GetApiInvokeCountCounter(b.ctx).With(prometheus.Labels{"api_name": "beer", "success": "true"}).Add(1)
	} else {
		logger.Error("Failed in processing", zap.Int("http_code", httpCode))
		w.WriteHeader(httpCode)
		w.Write([]byte(utils.HttpErrorMessage))
		metrics.GetApiInvokeCountCounter(b.ctx).With(prometheus.Labels{"api_name": "beer", "success": "false"}).Add(1)
	}
}
