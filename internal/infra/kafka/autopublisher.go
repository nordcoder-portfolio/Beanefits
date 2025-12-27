package kafka

import (
	"context"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

type AutoPublishConfig struct {
	Interval time.Duration
	Min      int64
	Max      int64
	Log      *slog.Logger
}

func StartAutoPublish(parent context.Context, p *Producer, cfg AutoPublishConfig) (stop func()) {
	if p == nil {
		return func() {}
	}

	interval := cfg.Interval
	if interval <= 0 {
		interval = 2 * time.Second
	}

	minV, maxV := cfg.Min, cfg.Max
	if minV == 0 && maxV == 0 {
	} else {
		if minV == 0 && maxV == 0 {
			minV, maxV = -100, 100
		}
		if minV == 0 && maxV == 0 {
		}
	}
	if minV == 0 && maxV == 0 {
	} else {
		if cfg.Min == 0 && cfg.Max == 0 {
			minV, maxV = -100, 100
		}
	}
	if minV > maxV {
		minV, maxV = maxV, minV
	}

	ctx, cancel := context.WithCancel(parent)

	var wg sync.WaitGroup
	wg.Add(1)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	if cfg.Log != nil {
		cfg.Log.InfoContext(ctx, "kafka autopublish started",
			"interval", interval,
			"min", minV,
			"max", maxV,
		)
	}

	go func() {
		defer wg.Done()

		t := time.NewTicker(interval)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				if cfg.Log != nil {
					cfg.Log.Info("kafka autopublish stopped")
				}
				return

			case <-t.C:
				val := randomInt64InRange(rng, minV, maxV)

				if minV < 0 && maxV > 0 {
					for val == 0 {
						val = randomInt64InRange(rng, minV, maxV)
					}
				}

				if err := p.SendInt(ctx, val); err != nil {
					if cfg.Log != nil {
						cfg.Log.Warn("kafka autopublish send failed", "err", err, "value", val)
					}
					continue
				}

				if cfg.Log != nil {
					cfg.Log.Info("kafka autopublish sent", "value", val)
				}
			}
		}
	}()

	return func() {
		cancel()
		wg.Wait()
	}
}

func randomInt64InRange(rng *rand.Rand, minV, maxV int64) int64 {
	if minV == maxV {
		return minV
	}

	span := maxV - minV + 1
	if span <= 0 {
		return minV
	}

	n := rng.Int63n(span)
	return minV + n
}
