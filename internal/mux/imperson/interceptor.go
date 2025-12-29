import (
	"context"
	"strings"
	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5" // 建議使用 jwt-go 
)

func NewImpersonationInterceptor(baseConfig *rest.Config) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// 1. 從 Header 提取 Authorization: Bearer <JWT>
			authHeader := req.Header().Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing token"))
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			// 2. 驗證 JWT 並提取資訊 (此處簡化驗證邏輯，實務需搭配 Keycloak 公鑰)
			claims := &jwt.MapClaims{}
			_, _, err := new(jwt.Parser).ParseUnverified(tokenStr, claims)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			// 假設 JWT 裡有 email 和 groups 欄位
			username := (*claims)["email"].(string)
			groups := []string{}
			for _, g := range (*claims)["groups"].([]interface{}) {
				groups = append(groups, g.(string))
			}

			// 3. 克隆基礎配置並設定 Impersonation
			impersonatedConfig := rest.CopyConfig(baseConfig)
			impersonatedConfig.Impersonate = rest.ImpersonationConfig{
				UserName: username,
				Groups:   groups,
			}

			// 4. 建立該身份專用的 Clientset (或你的 CRD Clientset)
			client, err := kubernetes.NewForConfig(impersonatedConfig)
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}

			// 5. 將 Client 塞入 Context 傳給下一層
			newCtx := context.WithValue(ctx, k8sClientKey{}, client)
			return next(newCtx, req)
		}
	}
}