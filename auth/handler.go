package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pquerna/otp/totp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/oauth2/callback",
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "random"
	userSecret       = "your-mfa-secret" // This should be securely stored and unique per user
)

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Render MFA form
	fmt.Fprintf(w, `<html><body>
        <form action="/mfa" method="POST">
            <label for="otp">Enter OTP:</label>
            <input type="text" id="otp" name="otp">
            <input type="hidden" name="access_token" value="%s">
            <input type="submit" value="Submit">
        </form>
        </body></html>`, token.AccessToken)
}

func handleMFA(w http.ResponseWriter, r *http.Request) {
	otp := r.FormValue("otp")
	accessToken := r.FormValue("access_token")

	if totp.Validate(otp, userSecret) {
		fmt.Fprintf(w, "MFA successful! Access Token: %s", accessToken)
	} else {
		fmt.Fprintf(w, "MFA failed! Invalid OTP.")
	}
}

func main() {
	http.HandleFunc("/oauth2/login", handleGoogleLogin)
	http.HandleFunc("/oauth2/callback", handleGoogleCallback)
	http.HandleFunc("/mfa", handleMFA)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
