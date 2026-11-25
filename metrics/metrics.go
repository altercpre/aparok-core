package metrics

import (
	"time"

	"github.com/google/uuid"

	"github.com/anyproto/anytype-heart/core/anytype/config/loadenv"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
)

var log = logging.Logger("anytype-telemetry")

var (
	DefaultInHouseKey string
	Enabled           bool
)

func GenerateAnalyticsId() string {
	return uuid.New().String()
}

type noopCounter struct{}

func (n noopCounter) Inc()         {}
func (n noopCounter) Add(float64)  {}

type noopCounterVec struct{}

func (n noopCounterVec) WithLabelValues(lvs ...string) noopCounter {
	return noopCounter{}
}

var (
	ObjectFTUpdatedCounter              = noopCounter{}
	ObjectDetailsUpdatedCounter         = noopCounter{}
	ObjectDetailsHeadsNotChangedCounter = noopCounter{}
	LinkPreviewStatusCounter            = noopCounterVec{}
)

func MetricTimeBuckets(scale []time.Duration) []float64 {
	buckets := make([]float64, len(scale))
	for i, b := range scale {
		buckets[i] = b.Seconds()
	}
	return buckets
}

func init() {
	if DefaultInHouseKey == "" {
		DefaultInHouseKey = loadenv.Get("INHOUSE_KEY")
	}
}
