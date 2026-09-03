package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "/app/data/invoices.db"
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("Failed to create data dir: %v", err)
	}
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS invoices (
		id TEXT PRIMARY KEY,
		filename TEXT,
		vendor_name TEXT,
		vendor_address TEXT,
		invoice_number TEXT,
		invoice_date TEXT,
		due_date TEXT,
		subtotal TEXT,
		tax TEXT,
		total TEXT,
		raw_text TEXT,
		created_at TEXT
	)`)
	if err != nil {
		log.Fatalf("Failed to create invoices table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS line_items (
		id TEXT PRIMARY KEY,
		invoice_id TEXT,
		description TEXT,
		quantity TEXT,
		unit_price TEXT,
		total TEXT
	)`)
	if err != nil {
		log.Fatalf("Failed to create line_items table: %v", err)
	}
}

func main() {
	initDB()
	defer db.Close()

	// Ensure upload and output dirs
	os.MkdirAll("/app/data/uploads", 0755)
	os.MkdirAll("/app/output/exports", 0755)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": "Invoice Extract API", "version": "1.0.0"})
	})

	router.POST("/extract", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "No file uploaded"})
			return
		}
		iid := uuid.New().String()
		filepath := fmt.Sprintf("/app/data/uploads/%s_%s", iid, file.Filename)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save file"})
			return
		}
		_, err = db.Exec("INSERT INTO invoices (id,filename,vendor_name,invoice_number,invoice_date,total,raw_text,created_at) VALUES (?,?,?,?,?,?,?,?)",
			iid, file.Filename, "Unknown Vendor", "N/A", "N/A", "N/A", "File: "+file.Filename, time.Now().UTC().Format(time.RFC3339))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"invoice_id": iid, "message": "Invoice uploaded, extraction pending"})
	})

	router.GET("/invoices", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM invoices ORDER BY created_at DESC")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		invoices := []map[string]interface{}{}
		for rows.Next() {
			var id, filename, vendorName, vendorAddress, invoiceNumber, invoiceDate, dueDate, subtotal, tax, total, rawText, createdAt string
			if err := rows.Scan(&id, &filename, &vendorName, &vendorAddress, &invoiceNumber, &invoiceDate, &dueDate, &subtotal, &tax, &total, &rawText, &createdAt); err != nil {
				continue
			}
			invoices = append(invoices, map[string]interface{}{
				"id": id, "filename": filename, "vendor_name": vendorName, "vendor_address": vendorAddress,
				"invoice_number": invoiceNumber, "invoice_date": invoiceDate, "due_date": dueDate,
				"subtotal": subtotal, "tax": tax, "total": total, "raw_text": rawText, "created_at": createdAt,
			})
		}
		c.JSON(200, gin.H{"invoices": invoices})
	})

	router.GET("/invoices/:iid", func(c *gin.Context) {
		iid := c.Param("iid")
		var id, filename, vendorName, vendorAddress, invoiceNumber, invoiceDate, dueDate, subtotal, tax, total, rawText, createdAt string
		err := db.QueryRow("SELECT * FROM invoices WHERE id=?", iid).Scan(&id, &filename, &vendorName, &vendorAddress, &invoiceNumber, &invoiceDate, &dueDate, &subtotal, &tax, &total, &rawText, &createdAt)
		if err != nil {
			c.JSON(404, gin.H{"error": "Invoice not found"})
			return
		}
		rows, err := db.Query("SELECT * FROM line_items WHERE invoice_id=?", iid)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		items := []map[string]interface{}{}
		for rows.Next() {
			var itemID, invoiceID, desc, qty, unitPrice, itemTotal string
			if err := rows.Scan(&itemID, &invoiceID, &desc, &qty, &unitPrice, &itemTotal); err != nil {
				continue
			}
			items = append(items, map[string]interface{}{
				"id": itemID, "invoice_id": invoiceID, "description": desc,
				"quantity": qty, "unit_price": unitPrice, "total": itemTotal,
			})
		}
		result := map[string]interface{}{
			"id": id, "filename": filename, "vendor_name": vendorName, "vendor_address": vendorAddress,
			"invoice_number": invoiceNumber, "invoice_date": invoiceDate, "due_date": dueDate,
			"subtotal": subtotal, "tax": tax, "total": total, "raw_text": rawText,
			"created_at": createdAt, "line_items": items,
		}
		c.JSON(200, result)
	})

	router.GET("/export", func(c *gin.Context) {
		rows, err := db.Query("SELECT id,filename,vendor_name,invoice_number,invoice_date,due_date,total FROM invoices ORDER BY created_at DESC")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		filename := "invoices_export.csv"
		filepath := "/app/output/exports/" + filename
		f, err := os.Create(filepath)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer f.Close()

		w := csv.NewWriter(f)
		w.Write([]string{"id", "filename", "vendor_name", "invoice_number", "invoice_date", "due_date", "total"})
		for rows.Next() {
			var id, fn, vn, in, id2, dd, total string
			if err := rows.Scan(&id, &fn, &vn, &in, &id2, &dd, &total); err != nil {
				continue
			}
			w.Write([]string{id, fn, vn, in, id2, dd, total})
		}
		w.Flush()
		c.FileAttachment(filepath, filename)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Invoice Extract server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}