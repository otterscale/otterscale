package mux

import (
	"encoding/json"
	"net/http"

	"connectrpc.com/connect"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

type JWKSProxy struct {
	*http.ServeMux

	cache *jwk.Cache
	url   string
}

func NewJWKSProxy() *JWKSProxy {
	return &JWKSProxy{
		ServeMux: &http.ServeMux{},
	}
}

func (p *JWKSProxy) RegisterHandlers(_ []connect.HandlerOption) error {
	p.Handle("/.well-known/jwks.json", http.HandlerFunc(p.serve))
	return nil
}

func (p *JWKSProxy) SetCache(cache *jwk.Cache) {
	p.cache = cache
}

func (p *JWKSProxy) SetURL(url string) {
	p.url = url
}

func (p *JWKSProxy) serve(w http.ResponseWriter, r *http.Request) {
	set, err := p.cache.CachedSet(p.url)
	if err != nil {
		http.Error(w, "Failed to get JWKS", http.StatusInternalServerError)
		return
	}

	buf, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal JWKS", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(buf); err != nil {
		http.Error(w, "Failed to write JWKS", http.StatusInternalServerError)
		return
	}
}
