package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
)

const xApiGwReqId = "X-Apigw-Req-Id" // make configurable

func main() {
	e := echo.New()

	InitOauthConfig()
	SetupRoutes(e)

	e.Logger.Fatal(e.Start(":3000"))
}

func SetupRoutes(e *echo.Echo) {
	url, err := url.Parse("https://httpbin.org")
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/*", doNothing, AddRequestId, OidcAuth, middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{Name: "httpbin", URL: url},
	})))

	e.POST("/*", doNothing, AddRequestId, OidcAuth, middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{Name: "httpbin", URL: url},
	})))

	//e.GET("/login", loginHandler)

	e.GET("/healthz", healthHandler)

	e.GET("/callback", callbackHandler)

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("user-sessions"))))
}

func callbackHandler(c echo.Context) error {
	fmt.Println("In callback handler...")
	fmt.Printf("code=%s state=%s\n", c.QueryParam("code"), c.QueryParam("state"))

	s, err := session.Get("session", c)
	if err != nil {
		c.String(http.StatusForbidden, fmt.Sprintf("%s", err))
		return err
	}

	// Check the state that was returned in the query string is the same as the above state
	if c.QueryParam("state") == "" || c.QueryParam("state") != s.Values["state"] {
		c.String(http.StatusForbidden, fmt.Sprintf("%s", err))
		return err
	}

	// Make sure the code was provided
	if c.QueryParam("error") != "" {
		c.String(http.StatusForbidden, fmt.Sprintf("authorization server returned an error: %s", c.QueryParam("error")))
		return fmt.Errorf("authorization server returned an error: %s", c.QueryParam("error"))
	}

	// Make sure the code was provided
	if c.QueryParam("code") == "" {
		c.String(http.StatusForbidden, fmt.Sprintf("the code was not returned or is not accessible"))
		return fmt.Errorf("the code was not returned or is not accessible")
	}

	token, err := oktaOauthConfig.Exchange(
		context.Background(),
		c.QueryParam("code"),
		oauth2.SetAuthURLParam("status", s.Values["state"].(string)),
	)
	if err != nil {
		c.String(http.StatusUnauthorized, fmt.Sprintf("%s", err))
		return err
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.String(http.StatusUnauthorized, fmt.Sprintf("id token missing from OAuth2 token"))
		return fmt.Errorf("id token missing from OAuth2 token")
	}
	_, err = verifyToken(rawIDToken)

	if err != nil {
		c.String(http.StatusForbidden, fmt.Sprintf("%s", err))
		return err
	} else {
		s.Values["access_token"] = token.AccessToken

		s.Save(c.Request(), c.Response())
	}

	c.Redirect(http.StatusFound, c.QueryParam("state"))
	return nil
}

func verifyToken(t string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["aud"] = os.Getenv("OAUTH2_CLIENT_ID")
	jv := verifier.JwtVerifier{
		Issuer:           os.Getenv("OAUTH2_ISSUER"),
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified")
}

func loginHandler(c echo.Context) error {
	//state := c.Request().RequestURI
	//
	//newSession, err := session.Get("session", c)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//newSession.Options = &sessions.Options{
	//	Path:     "/",
	//	MaxAge:   86400 * 7,
	//	HttpOnly: true,
	//}
	//newSession.Values["state"] = state
	//newSession.Save(c.Request(), c.Response())
	//
	//c.Redirect(http.StatusTemporaryRedirect, oktaOauthConfig.AuthCodeURL(state))
	return nil
}

func doNothing(c echo.Context) error {
	return nil
}

func healthHandler(c echo.Context) error {
	fmt.Printf("In healthHandler: %s\n", c.Request().Header.Get(xApiGwReqId))
	c.Response().Status = 200
	return nil
}

func AddRequestId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := uuid.New().String()
		c.Logger().Infof("RequestID: %s", requestID)
		c.Set(xApiGwReqId, requestID)
		c.Request().Header.Set(xApiGwReqId, requestID)
		c.Response().Header().Set(xApiGwReqId, requestID)
		return next(c)
	}
}

func OidcAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !isAuthenticated(c) {
			fmt.Printf("Unauthorized route: %s", c.Request().RequestURI)
			// c.Redirect(http.StatusFound, "/login")

			state := c.Request().RequestURI

			newSession, err := session.Get("session", c)
			if err != nil {
				fmt.Println(err) // TODO
			}
			newSession.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   86400 * 7,
				HttpOnly: false,
				Secure:   false,
			}
			newSession.Values["state"] = state
			newSession.Save(c.Request(), c.Response())

			c.Redirect(http.StatusTemporaryRedirect, oktaOauthConfig.AuthCodeURL(state))

			return nil
		}
		return next(c)
	}
}

func isAuthenticated(c echo.Context) bool {
	s, _ := session.Get("session", c)
	if s == nil || s.IsNew {
		return false
	} else {
		return true // TODO check for if session is valid (does the session get cleaned up automatically?)
	}
}

var oktaOauthConfig = &oauth2.Config{}

func InitOauthConfig() {
	godotenv.Load("./.okta.env")

	oktaOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("CALLBACK_URL"),
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   os.Getenv("OAUTH2_ISSUER") + "/v1/authorize",
			TokenURL:  os.Getenv("OAUTH2_ISSUER") + "/v1/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}
