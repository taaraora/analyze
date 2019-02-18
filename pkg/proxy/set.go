package proxy

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/models"
)

type Set struct {
	logger  logrus.FieldLogger
	tr      http.RoundTripper
	m       sync.RWMutex
	proxies map[string]*httputil.ReverseProxy
}

func NewProxySet(transport http.RoundTripper, logger logrus.FieldLogger) *Set {
	return &Set{
		logger:  logger,
		tr:      transport,
		proxies: make(map[string]*httputil.ReverseProxy),
	}
}

func (s *Set) SetProxy(plugin *models.Plugin) error {
	s.m.Lock()
	defer s.m.Unlock()
	url, err := url.Parse("http://" + plugin.ServiceEndpoint)
	if err != nil {
		return errors.New("can't parse host")
	}

	s.logger.Debugf("create proxy for url: %+v", *url)
	reverseProxy := httputil.NewSingleHostReverseProxy(url)
	reverseProxy.Transport = s.tr
	reverseProxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		s.logger.Errorf("reverse proxy error params: %+v", *req.URL)
		s.logger.Errorf("reverse proxy error: %v", err)
		rw.WriteHeader(http.StatusBadGateway)
	}
	s.proxies[plugin.ID] = reverseProxy

	return nil
}

func (s *Set) GetProxies() map[string]*httputil.ReverseProxy {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.proxies
}

func (s *Set) RemoveProxy(pluginID string) {
	delete(s.proxies, pluginID)
}
