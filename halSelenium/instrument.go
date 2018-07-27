package halSelenium

import (
	"github.com/go-kit/kit/metrics"
	"github.com/tebeka/selenium"
	"golang.org/x/net/context"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	errorCount     metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentService(counter metrics.Counter, errorCount metrics.Counter,
	latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		errorCount:     errorCount,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) NewClient(seleniumServer string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "NewClient").Add(1)
		s.requestLatency.With("method", "NewClient").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "NewClient").Add(1)
		}
	}(time.Now())
	return s.Service.NewClient(seleniumServer)

}
func (s *instrumentingService) HandleSeleniumError(ctx context.Context, chatId uint32, internal bool, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "HandleSeleniumError").Add(1)
		s.requestLatency.With("method", "HandleSeleniumError").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.HandleSeleniumError(ctx, chatId, internal, err)
}
func (s *instrumentingService) Driver() selenium.WebDriver {
	defer func(begin time.Time) {
		s.requestCount.With("method", "Driver").Add(1)
		s.requestLatency.With("method", "Driver").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Driver()
}

func (s *instrumentingService) ClickByClassName(cn string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ClickByClassName").Add(1)
		s.requestLatency.With("method", "ClickByClassName").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.ClickByClassName(cn)
}
func (s *instrumentingService) ClickByXPath(xp string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ClickByXPath").Add(1)
		s.requestLatency.With("method", "ClickByXPath").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.ClickByXPath(xp)
}
func (s *instrumentingService) ClickByCSSSelector(cs string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ClickByCSSSelector").Add(1)
		s.requestLatency.With("method", "ClickByCSSSelector").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.ClickByCSSSelector(cs)
}

func (s *instrumentingService) WaitFor(findBy, selector string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "WaitFor").Add(1)
		s.requestLatency.With("method", "WaitFor").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.WaitFor(findBy, selector)
}
