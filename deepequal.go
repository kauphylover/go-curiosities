package main

import "reflect"

type JwtAuthnConfig struct {
	IdpName          string
	Issuer           string
	JwksUri          string
	CallbackEndpoint string
}

func main() {
	jwt := &JwtAuthnConfig{
		"asdf", "asdf", "jwks", "call",
	}

	jwt2 := &JwtAuthnConfig{
		"asdf", "asdf", "jwks", "call",
	}

	println(reflect.DeepEqual(jwt, jwt2))
}
