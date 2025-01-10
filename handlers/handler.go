package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
	"mobile-banking-v3/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

)



// Handlers
func RegisterHandler(c *fiber.Ctx, db *gorm.DB) error {
    user := new(models.User)
    if err := c.BodyParser(user); err != nil {
        return err
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)

    // Create user
    if err := db.Create(user).Error; err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Username already exists",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(user)
}

func LoginHandler(c *fiber.Ctx, db *gorm.DB) error {
    loginRequest := new(models.LoginRequest)
    if err := c.BodyParser(loginRequest); err != nil {
        return err
    }

    // Find user
    var user models.User
    if err := db.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }

    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }

    // Generate token
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = user.ID

    t, err := token.SignedString([]byte("your-secret-key"))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not generate token",
        })
    }

    return c.JSON(fiber.Map{
        "token": t,
    })
}

func GetBalanceHandler(c *fiber.Ctx, db *gorm.DB) error {
    userID := c.Locals("userID").(uint)
    var user models.User
    if err := db.First(&user, userID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    return c.JSON(fiber.Map{
        "balance": user.Balance,
    })
}

func TransferHandler(c *fiber.Ctx, db *gorm.DB) error {
    userID := c.Locals("userID").(uint)
    transferRequest := new(models.TransferRequest)
    if err := c.BodyParser(transferRequest); err != nil {
        return err
    }

    // Start transaction
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Get sender
    var sender models.User
    if err := tx.First(&sender, userID).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Sender not found",
        })
    }

    // Get recipient
    var recipient models.User
    if err := tx.Where("username = ?", transferRequest.ToUsername).First(&recipient).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Recipient not found",
        })
    }

    // Check balance
    if sender.Balance < transferRequest.Amount {
        tx.Rollback()
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Insufficient balance",
        })
    }

    // Update balances
    sender.Balance -= transferRequest.Amount
    recipient.Balance += transferRequest.Amount

    // Save changes
    if err := tx.Save(&sender).Error; err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Save(&recipient).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Create transaction record
    transaction := models.Transaction{
        FromUserID: sender.ID,
        ToUserID:   recipient.ID,
        Amount:     transferRequest.Amount,
        Type:       "transfer",
    }
    if err := tx.Create(&transaction).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Commit transaction
    tx.Commit()

    return c.JSON(fiber.Map{
        "message": "Transfer successful",
        "balance": sender.Balance,
    })
}

func DepositHandler(c *fiber.Ctx,  db *gorm.DB) error {
    // Get user ID from JWT token
    userID := c.Locals("userID").(uint)
    
    // Parse request body
    deposit := new(models.DepositRequest)
    if err := c.BodyParser(deposit); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Validate deposit amount
    if deposit.Amount <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Deposit amount must be greater than 0",
        })
    }

    // Start database transaction
    tx := db.Begin()
    
    // Get user
    var user models.User
    if err := tx.First(&user, userID).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    // Update balance
    user.Balance += deposit.Amount
    if err := tx.Save(&user).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update balance",
        })
    }

    // Create transaction record
    transaction := models.Transaction{
        FromUserID: userID,
        ToUserID:   userID, // same user for deposit
        Amount:     deposit.Amount,
        Type:       "deposit",
    }

    if err := tx.Create(&transaction).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to record transaction",
        })
    }

    // Commit transaction
    tx.Commit()

    // Return updated balance
    return c.JSON(fiber.Map{
        "message": "Deposit successful",
        "balance": user.Balance,
        "amount":  deposit.Amount,
    })
}
