package main

import (
	"encoding/json"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-co-op/gocron"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"io"
	"net/http"
	"strings"
	"time"
)

var idp map[string][]string

//const which = "okta"

const which = "csp"

//const which = "auth0"

//const which = "google"

func main() {
	idp = make(map[string][]string)
	idp["okta"] = []string{
		"https://dev-16151438.okta.com/oauth2/default", "eyJraWQiOiIwdHB6TnREQmEtZF9DeTVJOHJaUWMyYVVPX1pORGNXb29velNzZVFYb0VZIiwiYWxnIjoiUlMyNTYifQ.eyJ2ZXIiOjEsImp0aSI6IkFULlo4emJKZXppNUVyb0JTUzhzOFRtX2JOVkdsRDJNVHAzd0FCcXdwX3N2OVEiLCJpc3MiOiJodHRwczovL2Rldi0xNjE1MTQzOC5va3RhLmNvbS9vYXV0aDIvZGVmYXVsdCIsImF1ZCI6ImFwaTovL2RlZmF1bHQiLCJpYXQiOjE2NTMxMTE5MzcsImV4cCI6MTY1MzExNTUzNywiY2lkIjoiMG9hNTFzYnRua3VpY0tkSWc1ZDciLCJ1aWQiOiIwMHU1MXIxM2NmeDZKZjlSOTVkNyIsInNjcCI6WyJvcGVuaWQiLCJwcm9maWxlIl0sImF1dGhfdGltZSI6MTY1MzExMTkzNSwic3ViIjoiYW1hbGxlbGFAdm13YXJlLmNvbSJ9.UU9JH-mCPFruMFHuE3KhXBhUnPlmr_HgVQyYd6bqRnpxReeWN5ogmOuDMe2uulpI5MgW8bVVjveMConO_KN-pp7aCLtQNk3qX49ZBA46RVmNs6S3XYislDYqNAcWaqjom8zPytjixt-FX5eMMXEQvGmbLhBA6j2fdT3bivQTCTCjbo-xrtIXgJ-T-ABHqRz9QlOnaEl7lqSj-wayssdjWzc9cEDw_iexv_iOt_aBpVmdKZ2hBZe0zxZ9Cu0Ff8F1OfXjca3zITW-dkKB9o49ooZFk0We3QZ183PQoGSMuhNbT0jZNRFwHvKWxXvAXSN0FocmbxdfivjwfFuD6N67NA"}

	idp["csp"] = []string{
		"https://console-stg.cloud.vmware.com/csp/gateway/am/api", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InNpZ25pbmdfMiJ9.eyJzdWIiOiJ2bXdhcmUuY29tOjFjNjI5NGNmLWI3NWItNDRmYi04ZjAwLWZkNGU0MDViZDc5MCIsImlzcyI6Imh0dHBzOi8vZ2F6LXByZXZpZXcuY3NwLXZpZG0tcHJvZC5jb20iLCJjb250ZXh0X25hbWUiOiJhM2ZkOWY1MS0wY2JjLTQ2MDYtYWE2NS02MGZhOWQ4MzljZjkiLCJfbm9uY2UiOiI0MjljMTJmMC1kODJmLTExZWMtODcwNS0wMTk2YWRkYzkxZDQiLCJhenAiOiJTMmQ0ZnJ5eHdwQkpUUk1URjRGZFBDVlRqazhDQk92Z1M3MiIsImF1dGhvcml6YXRpb25fZGV0YWlscyI6W10sImRvbWFpbiI6InZtd2FyZS5jb20iLCJjb250ZXh0IjoiYjc3NGI3NjctYTNiOS00NzJjLTgwOGMtNjg0MGExODJlNTQyIiwicGVybXMiOltdLCJleHAiOjE2NTM0NTg3NzksImlhdCI6MTY1MzQ1Njk3OSwianRpIjoiMWEwNTUyYmEtYjQzZS00MWIzLTg1MjktODdiZTQ3ZmNlMDRkIiwiYWNjdCI6ImFtYWxsZWxhQHZtd2FyZS5jb20iLCJ1c2VybmFtZSI6ImFtYWxsZWxhIn0.AXACW_UuvPfITJcErGxsbuh2Ee1ICMo_KZC2VavNCXgj1NPU8tg5OoBSHky0rlxb5ecbL-DqoXhL861UVi_ujvjgXkByoGSINT6NIHGzR4Ozx-zMnpBu6FVjNzmeFNezLFSlnlF0ZWVD6f7Zpp40r_Y9WhF9kgi67G8FzdhMcJYc2gkjhYe3lqmGSfHApdeneBbCd0VEzzYoxxey0ziuyU9Ab4pufisCQtzE5_Kj2CsJ8BYS0qsvYyqS6CaKkTP6HLxuYGihMJ5TzzWvCFw0y8-SZ1nLRie4LYHVJUEJahK_xUbyl69FmrOpDWQISA3N7Xj2um35_fUf8lUNwyIGzQ"}

	idp["auth0"] = []string{
		"https://dev-hs46gmys.us.auth0.com",
		"eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIiwiaXNzIjoiaHR0cHM6Ly9kZXYtaHM0NmdteXMudXMuYXV0aDAuY29tLyJ9..CPA0Tfjmr3QZ7Ftf.9JfUL1TUpCQZQCpYb7q03rgmqwlvPVDOCrAe-LfN9jihni05QOlvIrvEgzi_FM0XRU-S14eJVelnwNL-5awm9mnQIk85mJ4frDIYCBEZl-v8Ce9YIT9TQpC4JviKuMlCf7v4WO1OgVDNmGZEuuCE56m1qyA3roORmUgNrJSA_NXERH4tqPMRlOrCzcc919Dx7JjQOeXorNdjJsgdXNNXS-RRNa-9lfly5RCJ0vjmPiJugULfQ4zaxM0ydehLwCPuhE7wFJaxa7kjiOd-PgF081IeDZcBiyulEN_9DYqTwB5Ve2nzY9xb8etrc3E.4hXQmWWwp6UShZUA6IviiA"}

	idp["google"] = []string{
		"https://accounts.google.com", "ya29.a0ARrdaM8CQYK3Zvj2W43fMyMvLTOib-JlvYVe1wsRil4eW7zWTcSbe0tANt1HxWxoIeXvN9BnDW_xGdeY1ltkoE28nG_Qg3HqXPyr__FlCYkVrpmXtrBjAJpTzIOZ7jeSh3HVnAVaJp2MrigyxOHdOkPaCGW7",
	}
	var issuerURL = idp[which][0]
	var accessToken = idp[which][1]

	// Create the JWKS from the resource at the given URL.
	keyFuncOptions := keyfunc.Options{
		RefreshInterval:  10 * time.Second,
		RefreshRateLimit: 10 * time.Second,
		RefreshErrorHandler: func(err error) {
			fmt.Println(err)
		},
		RefreshUnknownKID: true,
	}
	jwks, err := keyfunc.Get(getJwksUri(issuerURL), keyFuncOptions)
	if err != nil {
		fmt.Printf("Failed to get the JWKS from the given URL.\nError:%s", err.Error())
	}

	// Parse the JWT.
	_, err = jwt.Parse(accessToken, jwks.Keyfunc)
	if err != nil {
		fmt.Printf("failed to parse token: %s\n", err)
	}

	//fetchKeys(issuerURL)
	//_, err := jws.VerifySet([]byte(accessToken), keyset)
	//if err != nil {
	//	fmt.Printf("failed to parse payload: %s\n", err)
	//} else {
	//	fmt.Println("AUTH SUCCESS 1")
	//}
	//
	//remoteKeySet := oidc.NewRemoteKeySet(context.Background(), getJwksUri(issuerURL))
	//_, err = remoteKeySet.VerifySignature(context.Background(), accessToken)
	//if err != nil {
	//	fmt.Println("AUTH FAIL 2")
	//} else {
	//	fmt.Println("AUTH SUCCESS 2")
	//}
	//
	//// 3
	//jwks, err := keyfunc.Get("https://www.googleapis.com/oauth2/v3/certs", keyfunc.Options{
	//	RefreshErrorHandler: func(err error) {
	//		fmt.Printf("There was an error with the jwt.KeyFunc\nError:%s\n", err.Error())
	//	},
	//})

}

var scheduler *gocron.Scheduler
var keyset jwk.Set
var staticKeySet oidc.StaticKeySet

func fetchKeys(issuerURL string) {
	fmt.Printf("BEGIN Fetching authorization server keys from %s scheduler\n", issuerURL)

	jwksUri := getJwksUri(issuerURL)

	//if which == "csp" {
	//	jwksUri = issuerURL + "/am/api/auth/token-public-key"
	//}

	resp2, err := http.Get(jwksUri)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp2.Body.Close()

	bodyBytes, err := io.ReadAll(resp2.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	keyset, err = jwk.Parse(bodyBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// FOR CSP - set 'alg' if not set
	n := keyset.Len()
	for i := 0; i < n; i++ {
		key, _ := keyset.Get(i)
		if key.Algorithm() == "" {
			key.Set(jwk.AlgorithmKey, "RS256")
		}
	}

	fmt.Printf("DONE Fetching authorization server keys from %s scheduler\n", issuerURL)

}

func getJwksUri(issuerURL string) string {
	wellKnown := strings.TrimSuffix(issuerURL, "/") + "/.well-known/openid-configuration"
	resp, err := http.Get(wellKnown)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	var jsonObject map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	err = json.Unmarshal(bodyBytes, &jsonObject)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	jwksUri := jsonObject["jwks_uri"].(string)
	if jwksUri == "" {
		fmt.Println("jwks_uri empty. aborting")
		return ""
	}
	return jwksUri
}
