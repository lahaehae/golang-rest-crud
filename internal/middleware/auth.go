package middleware

import (
	"github.com/gin-gonic/gin"
	kratos "github.com/ory/kratos-client-go"

	"net/http"
)

type KratosAuthMiddleware struct {
    client *kratos.APIClient
}

func NewKratosAuthMiddleware() *KratosAuthMiddleware {
    configuration := kratos.NewConfiguration()
    configuration.Servers = []kratos.ServerConfiguration{
        {
            URL: "http://127.0.0.1:4433",
        },
    }

    return &KratosAuthMiddleware{
        client: kratos.NewAPIClient(configuration),
    }
}

func (m *KratosAuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        cookie, err := c.Cookie("ory_kratos_session")
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        // Проверяем сессию через Kratos
        session, _, err := m.client.FrontendAPI.ToSession(c).
            Cookie("ory_kratos_session=" + cookie).
            Execute()
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        // Добавляем ID пользователя в контекст
        c.Set("user_id", session.Identity.Id)
        c.Next()
    }
}