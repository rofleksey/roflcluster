package step

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/melbahja/goph"
	"log"
	"net/http"
	"sneakerdocs/config"
	"time"
)

type HealthCheckStep struct {
	Url string `yaml:"url"`
}

func (s *HealthCheckStep) Execute(client *goph.Client, cfg *config.Config) error {
	healthCheck := func() error {
		resp, err := http.Get(s.Url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("invalid status code: %d", resp.StatusCode)
		}

		return nil
	}

	onRetryFunc := func(n uint, err error) {
		log.Println(err.Error())
	}

	err := retry.Do(healthCheck, retry.OnRetry(onRetryFunc), retry.Attempts(15), retry.Delay(time.Second*20),
		retry.DelayType(retry.FixedDelay))
	if err != nil {
		return fmt.Errorf("timed out waiting for health check: %w", err)
	}

	return nil
}

func (s *HealthCheckStep) String() string {
	return "Health check url " + s.Url
}
