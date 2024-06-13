package middleware

import (
	"encoding/json"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type AuthProvider struct {
	Issuer string
	AppKey string
	Secret string
	Keys   jose.JSONWebKeySet
}

func NewAuthProvider(issuer, appKey, secret string) *AuthProvider {
	v := []byte(`{
		"keys": [{"kty": "RSA", "alg": "RS256", "use": "sig", "n": "vADQlQI1Fbjq_aJAK7epi5G86y12B44OI5oEjDMvI1jgJJ0VdDmbZ2lu7xb8Wr_RYtgtpoPO7kDIb-2-Vz0sx9Fnhy6ZYgy2PdP5WWcTGyKn2p8H2JcMGymOvDWfrzqYMPFTx1xGt-0eT9Yt10cbOrCCOnsFZfkqBmGyEOXYhr25oV_T5t5d49TKZ-G78w0NoJbTzgWtc89HX1OOj9yn49tonUdfVCCd0kU5Jng9YpTEVBW0sLlrUFIwHbS32fuAXhZpLrZpX6aTaSUDhyFFMJwAxYOdbxgPJjz0ZZsGNIVdUHRGqkys4_0cUKrZQTpFI63g0D0YEPiXcSbcM8vkdQ", "e": "AQAB", "kid": "vvg"}]
	}`)
	keys := jose.JSONWebKeySet{}
	json.Unmarshal(v, &keys)
	return &AuthProvider{
		Issuer: issuer,
		AppKey: appKey,
		Secret: secret,
		Keys:   keys,
	}
}

func (ap *AuthProvider) GetUserInfo(tokenString1 string) error {
	tokenString := "eyJhbGciOiJSUzI1NiIsImtpZCI6InZ2ZyIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QiLCJvdSI6ImNkIiwic2NvcGVzIjpbIjEiLCIyIl19.MuINdNH_-t6ryKC7z6Xz-D7e2eOlPzxZHECIO4vqB8DQ2nXZPOrTKhO3adMPQ46T6ZbV-Cw9pRHXXRjfqYoQWRZd4s9cmq_HRY4mCDHKPBXcqbkMvqNWGPy7fw7D3Ba9mK0dJsoU-mxnbF7--xehUKo1mdHVjjxFoA9T9McDpdYlKnOHyS1o3qNwmx8N8Ah43HeIJ-d7YG-F09sY9bGTHJ-SUvaJxUcRXoaSNow1dMrlSVSTXO37TB_CiXCEmndZCXDik8srskL2mEbHsAme-K1mKHdgi5MyBt-fnEOZ0WJvBLcSUkxnlQN0l618GozH-98wN_xMq76f5dDqqKKOeQ"
	token, err := jwt.ParseSigned(tokenString, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return err
	}
	keys := ap.Keys.Key(token.Headers[0].KeyID)
	for _, k := range keys {
		tmp := map[string]interface{}{}
		if err := token.Claims(k, &tmp); err == nil {
			v, _ := json.Marshal(tmp)
			println("Claims ", string(v))
		}
	}
	return nil
}
