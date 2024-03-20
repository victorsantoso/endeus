package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/victorsantoso/endeus/domain"
	"github.com/victorsantoso/endeus/helper"
)

// AuthMiddleware to validate user access control
func AuthMiddleware(userRepository domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// validate authorization
		authorization := c.GetHeader("Authorization")
		if len(authorization) == 0 {
			handleForbiddenAccess(c)
			return
		}
		// validate bearer token existence
		token := authorization[len("Bearer "):]
		jwtObj, err := helper.VerifyJwt(token)
		var registeredClaims *jwt.RegisteredClaims
		if jwtObj == nil {
			handleForbiddenAccess(c)
			return
		}
		// validate jwtObj with type assertions
		switch t := jwtObj.Claims.(type) {
		case *jwt.RegisteredClaims:
			registeredClaims = t
			if !jwtObj.Valid {
				handleForbiddenAccess(c)
				return
			}
		// set default status bad request
		default:
			handleForbiddenAccess(c)
			return
		}
		// validate err and signature
		if err != nil || err == jwt.ErrSignatureInvalid {
			handleForbiddenAccess(c)
			return
		}
		// validate id and role
		id, err := strconv.Atoi(registeredClaims.ID)
		if err != nil {
			handleForbiddenAccess(c)
			return
		}
		// validate user's existence
		validateUser, err := userRepository.FindById(context.Background(), int64(id))
		if validateUser == nil || err != nil {
			handleForbiddenAccess(c)
			return
		}
		// validate user's role
		if validateUser.Role != registeredClaims.Subject {
			handleForbiddenAccess(c)
			return
		}
		// set context with validated data
		c.Set("user", validateUser)
	}
}

// handle forbidden access for invalid access control
func handleForbiddenAccess(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, &AuthMiddlewareResponse{
		Message: "forbidden access",
		Code:    http.StatusForbidden,
	})
}

type AuthMiddlewareResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
