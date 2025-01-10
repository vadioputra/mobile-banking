package pkg

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
)

// JWT middleware
func JwtMiddleware(c *fiber.Ctx) error {
    token := c.Get("Authorization")
    if token == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Missing authorization token",
        })
    }

    // Verify token
    claims := jwt.MapClaims{}
    _, err := jwt.ParseWithClaims(token[7:], claims, func(token *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key"), nil
    })

    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }

    c.Locals("userID", uint(claims["user_id"].(float64)))
    return c.Next()
}