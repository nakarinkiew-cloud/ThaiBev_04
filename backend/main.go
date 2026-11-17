package main

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Database Model
type UserProfile struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	First_Name     string    `gorm:"not null" json:"first_name"`
	Last_Name      string    `gorm:"not null" json:"last_name"`
	Email          string    `gorm:"not null" json:"email"`
	Phone          string    `gorm:"not null" json:"phone"`
	Profile_Base64 string    `gorm:"type:text;not null" json:"profile_base64"`
	Birth_Day      time.Time `gorm:"not null" json:"birth_day"`
	Occupation     string    `gorm:"not null" json:"occupation"`
	Sex            string    `gorm:"not null" json:"sex"`
	Create_data    time.Time `gorm:"not null"`
}

// Request Payload DTO
type CreateUserRequest struct {
	First_Name     string `json:"first_name" binding:"required"`
	Last_Name      string `json:"last_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone" binding:"required"`
	Profile_Base64 string `json:"profile_base64" binding:"required"`
	Birth_Day      string `json:"birth_day" binding:"required"`
	Occupation     string `json:"occupation" binding:"required"`
	Sex            string `json:"sex" binding:"required,oneof=Male Female"`
}

func main() {
	// เชื่อมต่อฐานข้อมูล SQLite
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}
	db.AutoMigrate(&UserProfile{})

	r := gin.Default()
	r.Use(cors.Default()) // เปิดใช้งาน CORS สำหรับ Frontend

	r.POST("/api/profiles", func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกข้อมูลให้ครบถ้วนและถูกต้อง: " + err.Error()})
			return
		}

		// 1. ตรวจสอบ Format เบอร์โทรศัพท์ (ตัวเลข 9-10 หลัก ขึ้นต้นด้วย 0)
		phoneRegex := regexp.MustCompile(`^0[0-9]{8,9}$`)
		if !phoneRegex.MatchString(req.Phone) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบเบอร์โทรศัพท์ไม่ถูกต้อง"})
			return
		}

		// 2. ตรวจสอบและแปลง Format วันที่ วัน/เดือน/ปี (DD/MM/YYYY)
		parsedDate, err := time.Parse("02/01/2006", req.Birth_Day)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบวันเกิดต้องเป็น วัน/เดือน/ปี (DD/MM/YYYY)"})
			return
		}

		// บันทึกลงฐานข้อมูล
		profile := UserProfile{
			First_Name:     req.First_Name,
			Last_Name:      req.Last_Name,
			Email:          req.Email,
			Phone:          req.Phone,
			Profile_Base64: req.Profile_Base64,
			Birth_Day:      parsedDate,
			Occupation:     req.Occupation,
			Sex:            req.Sex,
			Create_data:    time.Now(),
		}

		if err := db.Create(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถบันทึกข้อมูลลงฐานข้อมูลได้"})
			return
		}

		// ส่งคืน ID ที่ถูกสร้างจากฐานข้อมูล
		c.JSON(http.StatusCreated, gin.H{
			"message": "save data success",
			"id":      profile.ID,
		})
	})

	r.Run(":8080")
}
