package impersonation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
		newCtx, err := i.enrichContextWithSubject(ctx, req.Header())
		if err != nil {
			return nil, err
		}
		return next(newCtx, req)
	}
}

func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		newCtx, err := i.enrichContextWithSubject(ctx, conn.RequestHeader())
		if err != nil {
			return err
		}
		return next(newCtx, conn)
	}
}

func (i *Interceptor) enrichContextWithSubject(ctx context.Context, h http.Header) (context.Context, error) {
	token, err := extractToken(h)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	idToken, err := i.verifier.Verify(ctx, token)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
	}

	return context.WithValue(ctx, subjectKey, idToken.Subject), nil
}

func extractToken(h http.Header) (string, error) {
	auth := h.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.Fields(auth)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

func GetSubject(ctx context.Context) (string, bool) {
	sub, ok := ctx.Value(subjectKey).(string)
	return sub, ok
}
