package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	e := echo.New()

	//e.GET("/api", kubeproxy)
	//e.GET("/apis", kubeproxy)
	//e.GET("/nexus", nexusApiHandler)

	rproxy := echo.HandlerFunc(func(c echo.Context) error {
		if c.Request().Header.Get("proxy") != "" {
			return proxyHandler(c)
		} else {
			uri, err := url.ParseRequestURI(c.Request().RequestURI)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			c.Redirect(http.StatusFound, uri.Path)
		}
		return nil
	})

	e.GET("*", rproxy)
	//e.GET("/get", func(c echo.Context) error {
	//	c.String(http.StatusOK, "called /get")
	//	return nil
	//})

	e.Logger.Fatal(e.Start(":3000"))
}

func proxyHandler(c echo.Context) error {
	req := c.Request()
	res := c.Response()
	tgt := getProxyTarget(req)

	// Fix header
	// Basically it's not good practice to unconditionally pass incoming x-real-ip header to upstream.
	// However, for backward compatibility, legacy behavior is preserved unless you configure Echo#IPExtractor.
	if req.Header.Get(echo.HeaderXRealIP) == "" || c.Echo().IPExtractor != nil {
		req.Header.Set(echo.HeaderXRealIP, c.RealIP())
	}
	if req.Header.Get(echo.HeaderXForwardedProto) == "" {
		req.Header.Set(echo.HeaderXForwardedProto, c.Scheme())
	}
	if c.IsWebSocket() && req.Header.Get(echo.HeaderXForwardedFor) == "" { // For HTTP, it is automatically set by Go HTTP reverse proxy.
		req.Header.Set(echo.HeaderXForwardedFor, c.RealIP())
	}

	// Proxy
	switch {
	case req.Header.Get(echo.HeaderAccept) == "text/event-stream":
	default:
		proxyHTTP(tgt, c).ServeHTTP(res, req)
	}
	if e, ok := c.Get("_error").(error); ok {
		return e
	}
	return nil
}

func getProxyTarget(req *http.Request) *ProxyTarget {
	//upstream, _ := url.Parse("default.nexus-api-gw")
	upstream, _ := url.Parse("https://httpbin.org")
	return &ProxyTarget{
		Name: "httpbin",
		URL:  upstream,
	}
}

// ProxyTarget defines the upstream target.
type ProxyTarget struct {
	Name string
	URL  *url.URL
	Meta echo.Map
}

const StatusCodeContextCanceled = 499

func proxyHTTP(tgt *ProxyTarget, c echo.Context) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(tgt.URL)
	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		desc := tgt.URL.String()
		if tgt.Name != "" {
			desc = fmt.Sprintf("%s(%s)", tgt.Name, tgt.URL.String())
		}
		// If the client canceled the request (usually by closing the connection), we can report a
		// client error (4xx) instead of a server error (5xx) to correctly identify the situation.
		// The Go standard library (at of late 2020) wraps the exported, standard
		// context.Canceled error with unexported garbage value requiring a substring check, see
		// https://github.com/golang/go/blob/6965b01ea248cabb70c3749fd218b36089a21efb/src/net/net.go#L416-L430
		if err == context.Canceled || strings.Contains(err.Error(), "operation was canceled") {
			httpError := echo.NewHTTPError(StatusCodeContextCanceled, fmt.Sprintf("client closed connection: %v", err))
			httpError.Internal = err
			c.Set("_error", httpError)
		} else {
			httpError := echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("remote %s unreachable, could not forward: %v", desc, err))
			httpError.Internal = err
			c.Set("_error", httpError)
		}
	}
	//proxy.Transport = config.Transport
	//proxy.ModifyResponse = config.ModifyResponse
	return proxy
}

func ProxyMiddlewareFunc() echo.MiddlewareFunc {
	urlHttpbin, _ := url.Parse("https://httpbin.org")
	return middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{Name: "httpbin", URL: urlHttpbin},
	}))
}
