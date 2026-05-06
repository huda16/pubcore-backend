package database

import (
	"fmt"
	"log"
	"os"
	"pubcore/models"

	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	err error
)

// StartDB initializes the database connection, runs migrations, and seeds demo data.
func StartDB() {
	host := os.Getenv("PGHOST")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbPort := os.Getenv("PGPORT")
	dbname := os.Getenv("PGDATABASE")
	sslMode := os.Getenv("PGSSLMODE")

	if sslMode == "" {
		sslMode = "disable"
	}

	config := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, dbPort, sslMode,
	)

	// Configure GORM logger to be less aggressive with slow query warnings
	// given the high latency of the free cloud service.
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             2000 * time.Millisecond, // Increase threshold to 500ms
			LogLevel:                  logger.Warn,             // Log level
			IgnoreRecordNotFoundError: true,                    // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                    // Enable color
		},
	)

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("error connecting to database:", err)
	}

	// Setup Connection Pool
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	fmt.Println("success connect to database")

	// Auto-migrate all models
	db.AutoMigrate(
		&models.Users{},
		&models.Authors{},
		&models.Publishers{},
		&models.Books{},
	)

	// Seed initial demo data (only if tables are empty)
	seedDemoData(db)
}

// GetDB returns the global database instance.
func GetDB() *gorm.DB {
	return db
}

// seedDemoData inserts initial demo data into empty tables.
func seedDemoData(db *gorm.DB) {
	// Seed Users
	var userCount int64
	db.Model(&models.Users{}).Count(&userCount)
	if userCount == 0 {
		users := []models.Users{
			{
				Username: "admin",
				Email:    "admin@pubcore.io",
				Password: "admin123",
				Bio:      "Platform administrator",
			},
			{
				Username: "editor",
				Email:    "editor@pubcore.io",
				Password: "editor123",
				Bio:      "Content editor",
			},
		}
		db.Create(&users)
		fmt.Println("Seeded users")
	}

	// Seed Publishers
	var pubCount int64
	db.Model(&models.Publishers{}).Count(&pubCount)
	if pubCount == 0 {
		publishers := []models.Publishers{
			{Name: "Bloomsbury", Address: "50 Bedford Square, London, WC1B 3DP", Website: "https://www.bloomsbury.com"},
			{Name: "Penguin Random House", Address: "1745 Broadway, New York, NY 10019", Website: "https://www.penguinrandomhouse.com"},
			{Name: "HarperCollins", Address: "195 Broadway, New York, NY 10007", Website: "https://www.harpercollins.com"},
		}
		db.Create(&publishers)
		fmt.Println("Seeded publishers")
	}

	// Seed Authors
	var authorCount int64
	db.Model(&models.Authors{}).Count(&authorCount)
	if authorCount == 0 {
		authors := []models.Authors{
			{Name: "J.K. Rowling", Bio: "British author best known for the Harry Potter series.", Nationality: "British"},
			{Name: "George Orwell", Bio: "English novelist and essayist, journalist and critic.", Nationality: "British"},
			{Name: "Haruki Murakami", Bio: "Japanese writer whose works of fiction and non-fiction have sold millions of copies worldwide.", Nationality: "Japanese"},
		}
		db.Create(&authors)
		fmt.Println("Seeded authors")
	}

	// Seed Books
	var bookCount int64
	db.Model(&models.Books{}).Count(&bookCount)
	if bookCount == 0 {
		// Fetch seeded data IDs
		var authors []models.Authors
		var publishers []models.Publishers
		db.Find(&authors)
		db.Find(&publishers)

		if len(authors) >= 3 && len(publishers) >= 3 {
			books := []models.Books{
				{
					Title:       "Harry Potter and the Philosopher's Stone",
					ISBN:        "978-0-7475-3269-9",
					Description: "A young boy discovers he is a wizard and enrolls in Hogwarts School of Witchcraft and Wizardry.",
					Genre:       "Fantasy",
					Year:        1997,
					AuthorID:    authors[0].ID,
					PublisherID: publishers[0].ID,
				},
				{
					Title:       "Harry Potter and the Chamber of Secrets",
					ISBN:        "978-0-7475-3849-3",
					Description: "Harry Potter's second year at Hogwarts is marked by a series of mishaps.",
					Genre:       "Fantasy",
					Year:        1998,
					AuthorID:    authors[0].ID,
					PublisherID: publishers[0].ID,
				},
				{
					Title:       "Nineteen Eighty-Four",
					ISBN:        "978-0-452-28423-4",
					Description: "A dystopian novel set in a totalitarian society ruled by Big Brother.",
					Genre:       "Dystopian Fiction",
					Year:        1949,
					AuthorID:    authors[1].ID,
					PublisherID: publishers[1].ID,
				},
				{
					Title:       "Norwegian Wood",
					ISBN:        "978-0-375-70402-7",
					Description: "A nostalgic story of loss and sexuality, set in late-1960s Tokyo.",
					Genre:       "Literary Fiction",
					Year:        1987,
					AuthorID:    authors[2].ID,
					PublisherID: publishers[2].ID,
				},
				{
					Title:       "Kafka on the Shore",
					ISBN:        "978-1-4000-7927-4",
					Description: "A metaphysical novel interwoven with a mysterious alternative reality.",
					Genre:       "Magical Realism",
					Year:        2002,
					AuthorID:    authors[2].ID,
					PublisherID: publishers[2].ID,
				},
			}
			db.Create(&books)
			fmt.Println("Seeded books")
		}
	}
}
