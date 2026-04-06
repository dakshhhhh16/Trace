package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dakshhhhh16/trace/pkg/api"
	"github.com/dakshhhhh16/trace/pkg/aws"
	"github.com/dakshhhhh16/trace/pkg/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Print CLI Banner
	printBanner()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("[FATAL] Unable to load SDK config: %v", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Get bucket name from environment variable or use default
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "trace-scans"
	}

	// Initialize handlers with dependencies
	scanHandler := api.NewScanHandler(aws.BucketBasics{
		S3Client: s3Client,
	}, bucketName)

	// Bundle all handlers
	handlers := &routes.Handlers{
		ScanHandler: scanHandler,
	}

	// Setup Gin router
	router := gin.Default()

	// Premium Middleware: Custom Headers
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Powered-By", "Trace Engine")
		c.Writer.Header().Set("X-Trace-Version", "1.0.0-alpha")
		c.Next()
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Trace",
			"status":  "operational",
			"version": "1.0.0-alpha",
		})
	})

	// Setup routes with handlers
	routes.SetupRoutes(router, handlers)

	// Start server
	log.Println("[INFO] Trace Engine starting on port 7789...")
	if err := router.Run(":7789"); err != nil {
		log.Fatalf("[FATAL] Failed to start server: %v", err)
	}
}

func printBanner() {
	banner := `
 ______   ______     ______     ______     ______    
/\__  _\ /\  == \   /\  __ \   /\  ___\   /\  ___\   
\/_/\ \/ \ \  __<   \ \  __ \  \ \ \____  \ \  __\   
   \ \_\  \ \_\ \_\  \ \_\ \_\  \ \_____\  \ \_____\ 
    \/_/   \/_/ /_/   \/_/\/_/   \/_____/   \/_____/ 
                                                     
   >> Supply Chain Security & Vulnerability Analyzer
   >> (c) Daksh Pathak
`
	log.Println(banner)
}
