package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("")

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// สำหรับ User สร้างรายการใหม่
type CreateRequestInput struct {
	Title string `json:"title" binding:"required"`
}

// สำหรับ Admin อัปเดตสถานะและใส่เหตุผล
type UpdateMultipleRequestInput struct {
	IDs         []int  `json:"ids" binding:"required"` // รับเป็น [1, 2, 3]
	Status      string `json:"status" binding:"required"`
	AdminReason string `json:"admin_reason" binding:"required"`
}

type ApprovalResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	AdminReason string    `json:"admin_reason"`
	Status      string    `json:"status"`
	Username    string    `json:"username"` // โชว์ชื่อคนขอ
	CreatedAt   time.Time `json:"created_at"`
}

var db *sql.DB

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "กรุณาเข้าสู่ระบบ"})
			c.Abort()
			return
		}
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ไม่ถูกต้องหรือหมดอายุ"})
			c.Abort()
			return
		}

		// เก็บข้อมูลไว้ใช้ใน Handler
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleCheckMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "คุณไม่มีสิทธิ์เข้าถึงส่วนนี้"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func LoginHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 1. ตรวจสอบข้อมูลที่ส่งมาจาก Body (JSON)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอก Username และ Password"})
		return
	}

	// 2. ค้นหา User ใน Database
	var user User
	query := "SELECT id, username, password, role FROM users WHERE username = $1"
	err := db.QueryRow(query, input.Username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		// ถ้าหาไม่เจอ หรือ DB มีปัญหา
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบผู้ใช้งานนี้ในระบบ"})
		return
	}

	// 3. ตรวจสอบรหัสผ่านด้วย Bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัสผ่านไม่ถูกต้อง"})
		return
	}

	// 4. สร้าง JWT Token (ถ้าผ่านทุกขั้นตอนข้างบน)
	expirationTime := time.Now().Add(24 * time.Hour) // Token หมดอายุใน 24 ชั่วโมง
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// สร้าง Token โดยใช้ Secret Key จาก .env
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้าง Token ได้"})
		return
	}

	// 5. ส่ง Token และสิทธิ์กลับไปให้ Frontend
	c.JSON(http.StatusOK, gin.H{
		"message":  "เข้าสู่ระบบสำเร็จ",
		"token":    tokenString,
		"role":     user.Role,
		"username": user.Username,
	})
}

func CreateUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"` // 'admin' หรือ 'user'
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 1. Hash Password ก่อนบันทึก
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 2. บันทึกลง Database
	query := "INSERT INTO users (username, password, role) VALUES ($1, $2, $3)"
	_, err = db.Exec(query, input.Username, string(hashedPassword), input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User already exists or DB error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func CreateApproval(c *gin.Context) {
	var input CreateRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุชื่อรายการ"})
		return
	}

	username, _ := c.Get("username")
	var userID int
	db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)

	query := "INSERT INTO approval_requests (title, user_id) VALUES ($1, $2)"
	_, err := db.Exec(query, input.Title, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถเพิ่มรายการได้"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ส่งรายการเรียบร้อยแล้ว"})
}

func UpdateMultipleApprovals(c *gin.Context) {
	var input UpdateMultipleRequestInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง หรือลืมเลือกรายการ"})
		return
	}

	query := `UPDATE approval_requests 
              SET status = $1, admin_reason = $2 
              WHERE id = ANY($3)`

	result, err := db.Exec(query, input.Status, input.AdminReason, pq.Array(input.IDs))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตรายการได้"})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("ดำเนินการสำเร็จ %d รายการ", rowsAffected),
	})
}

func GetUserRequests(c *gin.Context) {
	username, _ := c.Get("username")

	query := `SELECT a.id, a.title, 
       COALESCE(a.admin_reason, '') as admin_reason,
       a.status, u.username, a.created_at 
FROM approval_requests a 
JOIN users u ON a.user_id = u.id 
WHERE u.username = $1 
ORDER BY a.created_at DESC`

	rows, err := db.Query(query, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ดึงข้อมูลล้มเหลว"})
		return
	}
	defer rows.Close()

	var requests []ApprovalResponse
	for rows.Next() {
		var r ApprovalResponse
		err := rows.Scan(&r.ID, &r.Title, &r.AdminReason, &r.Status, &r.Username, &r.CreatedAt)
		if err != nil {
			fmt.Println("SCAN ERROR:", err)
			continue
		}
		requests = append(requests, r)
	}

	c.JSON(http.StatusOK, requests)
}

func GetAllRequests(c *gin.Context) {
	query := `SELECT a.id, a.title, COALESCE(a.admin_reason, ''), a.status, u.username, a.created_at 
			  FROM approval_requests a 
			  JOIN users u ON a.user_id = u.id 
			  ORDER BY a.created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ดึงข้อมูลล้มเหลว"})
		return
	}
	defer rows.Close()

	var requests []ApprovalResponse
	for rows.Next() {
		var r ApprovalResponse
		rows.Scan(&r.ID, &r.Title, &r.AdminReason, &r.Status, &r.Username, &r.CreatedAt)
		requests = append(requests, r)
	}

	c.JSON(http.StatusOK, requests)
}

func main() {
	err := godotenv.Load()
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // URL ของ Angular
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/login", LoginHandler)

	userGroup := r.Group("/user").Use(AuthMiddleware())
	{
		userGroup.POST("/request", CreateApproval)
		userGroup.GET("/my-requests", GetUserRequests) // ดูคำขอของตัวเอง
	}

	adminGroup := r.Group("/admin").Use(AuthMiddleware(), RoleCheckMiddleware("admin"))
	{
		adminGroup.GET("/all-requests", GetAllRequests)
		adminGroup.POST("/create-user", CreateUser)
		adminGroup.PUT("/approve-multiple", UpdateMultipleApprovals)
	}

	fmt.Println("Server started at :8080")
	r.Run(":8080")
}
