package impersonation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/otterscale/otterscale/internal/config"
)

type contextKey struct{}

var subjectKey = contextKey{}

type Interceptor struct {
	verifier *oidc.IDTokenVerifier
}

func NewInterceptor(conf *config.Config) (*Interceptor, error) {
	const timeout = 30 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	provider, err := oidc.NewProvider(ctx, conf.KeycloakRealmURL())
	if err != nil {
		return nil, fmt.Errorf("failed to init oidc provider: %w", err)
	}

	return &Interceptor{
		verifier: provider.Verifier(&oidc.Config{
			ClientID: conf.KeycloakClientID(),
		}),
	}, nil
}

func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		token, err := extractToken(req)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}

		idToken, err := i.verifier.Verify(ctx, token)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid token: %w", err))
		}

		newCtx := context.WithValue(ctx, subjectKey, idToken.Subject)
		return next(newCtx, req)
	}
}

func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

func extractToken(req connect.AnyRequest) (string, error) {
	auth := req.Header().Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	parts := strings.Fields(auth)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid authorization format")
	}

	return parts[1], nil
}

func GetSubject(ctx context.Context) (string, bool) {
	sub, ok := ctx.Value(subjectKey).(string)
	return sub, ok
}
