package metrics

import (
	"fmt"
	"time"

	"github.com/tyriis/go-locking-service/internal/domain"
	"github.com/tyriis/go-locking-service/internal/repositories"
)

type MetricsUpdater struct {
	lockRepo       *repositories.LockRepository
	metricsService *PrometheusMetricsService
	logger         domain.Logger
	quit           chan struct{}
	updateInterval time.Duration
}

func NewMetricsUpdater(
	lockRepo *repositories.LockRepository,
	metricsService *PrometheusMetricsService,
	logger domain.Logger,
) *MetricsUpdater {
	return &MetricsUpdater{
		lockRepo:       lockRepo,
		metricsService: metricsService,
		logger:         logger,
		quit:           make(chan struct{}),
		updateInterval: 10 * time.Second,
	}
}

func (m *MetricsUpdater) Start() {
	go func() {
		ticker := time.NewTicker(m.updateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				count, err := m.lockRepo.Count()
				if err != nil {
					m.logger.Error("MetricsUpdater.Start > Failed to get locks count for metrics")
					continue
				}
				m.metricsService.SetLockCount(float64(count))
				m.logger.Debug("MetricsUpdater.Start > set lock count to " + fmt.Sprintf("%d", count))
			case <-m.quit:
				return
			}
		}
	}()
}

func (m *MetricsUpdater) Stop() {
	close(m.quit)
}
