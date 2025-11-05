# Database & Migrations Guide

This guide covers database operations, GORM usage, migrations, and data modeling in FiberStarter.

## Table of Contents

- [Overview](#overview)
- [Database Setup](#database-setup)
- [GORM Basics](#gorm-basics)
- [Models](#models)
- [Repository Pattern](#repository-pattern)
- [Migrations](#migrations)
- [Advanced Queries](#advanced-queries)
- [Best Practices](#best-practices)

## Overview

FiberStarter uses:
- **MySQL 8.0+** - Relational database
- **GORM** - Go ORM (Object-Relational Mapping)
- **SQL Migrations** - Version-controlled schema changes

### Key Files

- `database/database.go` - Database connection and configuration
- `models/` - GORM model definitions
- `repos/` - Repository pattern for data access
- `migrations/` - SQL migration scripts

## Database Setup

### Connection Configuration

Database connection is configured in `database/database.go`:

```go
package database

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func InitDb(url string) {
    var err error

    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,
            LogLevel:                  logger.Info,
            IgnoreRecordNotFoundError: true,
            Colorful:                  true,
        },
    )

    DBConn, err = gorm.Open(mysql.Open(url), &gorm.Config{
        Logger: newLogger,
    })

    if err != nil {
        panic(err)
    }

    fmt.Println("Database connected successfully")
}
```

### Connection String Format

```go
// providers/config.go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
    AppConfig.Db_User,
    AppConfig.Db_Password,
    AppConfig.Db_Host,
    AppConfig.Db_Port,
    AppConfig.Db_Name,
)
```

Example: `root:password@tcp(localhost:3306)/fiberstarter?charset=utf8mb4&parseTime=True&loc=Local`

### Connection Pool Configuration

```go
sqlDB, err := DBConn.DB()

// SetMaxIdleConns sets the maximum number of connections in the idle pool
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns sets the maximum number of open connections
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime sets the maximum lifetime of a connection
sqlDB.SetConnMaxLifetime(time.Hour)
```

## GORM Basics

### What is GORM?

GORM is a powerful ORM library that:
- Maps Go structs to database tables
- Provides type-safe query building
- Handles migrations automatically
- Supports relationships and associations
- Prevents SQL injection

### Basic Operations (CRUD)

**Create:**
```go
user := models.User{Username: "john", Email: "john@example.com"}
result := database.DBConn.Create(&user)
// user.ID is now populated
```

**Read:**
```go
var user models.User

// Find by primary key
database.DBConn.First(&user, 1) // WHERE id = 1

// Find by condition
database.DBConn.Where("email = ?", "john@example.com").First(&user)

// Find all
var users []models.User
database.DBConn.Find(&users)
```

**Update:**
```go
// Update single field
database.DBConn.Model(&user).Update("Username", "jane")

// Update multiple fields
database.DBConn.Model(&user).Updates(models.User{
    Username: "jane",
    Email:    "jane@example.com",
})
```

**Delete:**
```go
// Soft delete (default with gorm.Model)
database.DBConn.Delete(&user, 1)

// Permanent delete
database.DBConn.Unscoped().Delete(&user, 1)
```

## Models

### Defining Models

Models are Go structs with GORM tags in the `models/` directory.

**User Model (`models/user.go`):**

```go
package models

import "gorm.io/gorm"

type User struct {
    gorm.Model           // Adds ID, CreatedAt, UpdatedAt, DeletedAt
    Username string      `gorm:"size:255;not null"`
    Email    string      `gorm:"size:255;not null;uniqueIndex"`
    Password string      `gorm:"size:600;not null"`
}

// TableName overrides the default table name
func (User) TableName() string {
    return "users"
}
```

### The gorm.Model

```go
type gorm.Model struct {
    ID        uint           `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

This provides:
- **ID**: Auto-incrementing primary key
- **CreatedAt**: Timestamp on creation
- **UpdatedAt**: Timestamp on updates
- **DeletedAt**: Soft delete support

### GORM Tags Reference

| Tag | Description | Example |
|-----|-------------|---------|
| `primaryKey` | Set as primary key | `gorm:"primaryKey"` |
| `size` | Column size | `gorm:"size:255"` |
| `not null` | NOT NULL constraint | `gorm:"not null"` |
| `unique` | UNIQUE constraint | `gorm:"unique"` |
| `uniqueIndex` | Unique index | `gorm:"uniqueIndex"` |
| `index` | Regular index | `gorm:"index"` |
| `default` | Default value | `gorm:"default:'active'"` |
| `autoIncrement` | Auto increment | `gorm:"autoIncrement"` |
| `column` | Custom column name | `gorm:"column:user_name"` |
| `type` | Custom column type | `gorm:"type:varchar(100)"` |
| `comment` | Column comment | `gorm:"comment:'User email'"` |

### Example: Book Model

```go
package models

import "gorm.io/gorm"

type Book struct {
    gorm.Model
    Title  string `gorm:"size:255;not null"`
    Author string `gorm:"size:255;not null"`
    Rating int    `gorm:"check:rating >= 1 AND rating <= 5"`
}
```

### Relationships

**One-to-Many Example:**

```go
// User has many posts
type User struct {
    gorm.Model
    Username string
    Posts    []Post
}

type Post struct {
    gorm.Model
    Title  string
    Body   string
    UserID uint
    User   User
}
```

**Usage:**

```go
// Fetch user with posts
var user models.User
database.DBConn.Preload("Posts").First(&user, 1)

// Create post for user
post := models.Post{
    Title:  "My First Post",
    Body:   "Hello World",
    UserID: user.ID,
}
database.DBConn.Create(&post)
```

**Many-to-Many Example:**

```go
type User struct {
    gorm.Model
    Username string
    Roles    []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    gorm.Model
    Name  string
    Users []User `gorm:"many2many:user_roles;"`
}
```

## Repository Pattern

The repository pattern abstracts database operations:

```
Controller → Repository → Database
```

### Why Use Repositories?

- **Separation of Concerns**: Business logic separate from data access
- **Testability**: Easy to mock database operations
- **Reusability**: Common queries in one place
- **Maintainability**: Changes in one location

### Example: User Repository

**File: `repos/user.go`**

```go
package repos

import (
    "fiberstarter/database"
    "fiberstarter/models"
)

// Get all users
func GetAllUsers() []models.User {
    var users []models.User
    database.DBConn.Find(&users)
    return users
}

// Get user by ID
func GetUserByID(id uint) models.User {
    var user models.User
    database.DBConn.First(&user, id)
    return user
}

// Get user by username
func GetUserByUsername(username string) models.User {
    var user models.User
    database.DBConn.Where("username = ?", username).First(&user)
    return user
}

// Check if user exists by email
func CheckIfUserExists(email string) bool {
    var user models.User
    result := database.DBConn.Where("email = ?", email).First(&user)
    return result.RowsAffected > 0
}

// Create new user
func CreateUser(user *models.User) error {
    result := database.DBConn.Create(user)
    return result.Error
}

// Update user
func UpdateUser(user *models.User) error {
    result := database.DBConn.Save(user)
    return result.Error
}

// Delete user (soft delete)
func DeleteUser(id uint) error {
    result := database.DBConn.Delete(&models.User{}, id)
    return result.Error
}
```

### Using Repositories in Controllers

```go
// controllers/Index.go
func ShowHomePage(c *fiber.Ctx) error {
    users := repos.GetAllUsers()

    return c.Render("Home/index", fiber.Map{
        "users": users,
    })
}
```

### Creating a New Repository

**1. Create the repository file:**

```go
// repos/post.go
package repos

import (
    "fiberstarter/database"
    "fiberstarter/models"
)

func GetAllPosts() []models.Post {
    var posts []models.Post
    database.DBConn.Preload("User").Find(&posts)
    return posts
}

func GetPostByID(id uint) (models.Post, error) {
    var post models.Post
    result := database.DBConn.Preload("User").First(&post, id)
    return post, result.Error
}

func CreatePost(post *models.Post) error {
    result := database.DBConn.Create(post)
    return result.Error
}
```

**2. Use in controller:**

```go
func ShowPosts(c *fiber.Ctx) error {
    posts := repos.GetAllPosts()
    return c.Render("posts/index", fiber.Map{"posts": posts})
}
```

## Migrations

### What are Migrations?

Migrations are version-controlled database schema changes:
- Track database changes over time
- Enable rollback of changes
- Ensure consistency across environments
- Documented database structure

### Migration File Naming

```
migrations/
├── m_00001_Up.sql    # First migration (up)
├── m_00001_Down.sql  # First migration (down/rollback)
├── m_00002_Up.sql    # Second migration (up)
└── m_00002_Down.sql  # Second migration (down/rollback)
```

**Naming Convention:**
- `m_[number]_Up.sql` - Apply changes
- `m_[number]_Down.sql` - Rollback changes
- Number with leading zeros (00001, 00002, etc.)

### Initial Migration

**File: `migrations/m_00001_Up.sql`**

```sql
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL UNIQUE,
  `password` varchar(600) NOT NULL,
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `books` (
  `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `title` varchar(255) NOT NULL,
  `author` varchar(255) NOT NULL,
  `rating` smallint NULL,
  UNIQUE KEY `rating` (`rating`),
  KEY `idx_books_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**File: `migrations/m_00001_Down.sql`**

```sql
DROP TABLE IF EXISTS `books`;
DROP TABLE IF EXISTS `users`;
```

### Running Migrations

**Apply migrations:**
```bash
mysql -u root -p fiberstarter < migrations/m_00001_Up.sql
```

**Rollback migrations:**
```bash
mysql -u root -p fiberstarter < migrations/m_00001_Down.sql
```

### Creating a New Migration

**Example: Add `profile_picture` to users**

**File: `migrations/m_00002_Up.sql`**

```sql
ALTER TABLE `users`
ADD COLUMN `profile_picture` varchar(255) NULL;
```

**File: `migrations/m_00002_Down.sql`**

```sql
ALTER TABLE `users`
DROP COLUMN `profile_picture`;
```

**Apply the migration:**
```bash
mysql -u root -p fiberstarter < migrations/m_00002_Up.sql
```

### Migration Best Practices

1. **Always include both Up and Down**
   - Up: Apply changes
   - Down: Rollback changes

2. **One logical change per migration**
   ```
   ✅ m_00002 - Add profile_picture column
   ❌ m_00002 - Add 10 different columns
   ```

3. **Never modify existing migrations**
   - Create new migrations instead
   - Existing migrations may have already run

4. **Test rollbacks**
   ```bash
   # Apply
   mysql -u root -p db < m_00002_Up.sql
   # Rollback
   mysql -u root -p db < m_00002_Down.sql
   ```

5. **Use IF EXISTS / IF NOT EXISTS**
   ```sql
   CREATE TABLE IF NOT EXISTS `users` ...
   ALTER TABLE `users` ADD COLUMN IF NOT EXISTS ...
   ```

## Advanced Queries

### Query Chaining

```go
// Build complex queries
var users []models.User
database.DBConn.
    Where("created_at > ?", time.Now().AddDate(0, -1, 0)).
    Where("email LIKE ?", "%@example.com").
    Order("username ASC").
    Limit(10).
    Find(&users)
```

### Select Specific Fields

```go
var users []models.User
database.DBConn.Select("username", "email").Find(&users)
```

### Counting Records

```go
var count int64
database.DBConn.Model(&models.User{}).
    Where("created_at > ?", time.Now().AddDate(0, -1, 0)).
    Count(&count)
```

### Pagination

```go
func GetUsersPaginated(page, pageSize int) []models.User {
    var users []models.User
    offset := (page - 1) * pageSize

    database.DBConn.
        Offset(offset).
        Limit(pageSize).
        Find(&users)

    return users
}
```

### Raw SQL Queries

```go
type Result struct {
    Username string
    PostCount int
}

var results []Result
database.DBConn.Raw(`
    SELECT u.username, COUNT(p.id) as post_count
    FROM users u
    LEFT JOIN posts p ON u.id = p.user_id
    GROUP BY u.id
`).Scan(&results)
```

### Transactions

```go
func TransferCredits(fromUserID, toUserID uint, amount int) error {
    return database.DBConn.Transaction(func(tx *gorm.DB) error {
        var fromUser, toUser models.User

        // Lock rows for update
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&fromUser, fromUserID).Error; err != nil {
            return err
        }

        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&toUser, toUserID).Error; err != nil {
            return err
        }

        // Check balance
        if fromUser.Credits < amount {
            return errors.New("insufficient credits")
        }

        // Update balances
        tx.Model(&fromUser).Update("credits", fromUser.Credits - amount)
        tx.Model(&toUser).Update("credits", toUser.Credits + amount)

        return nil // Commit transaction
    })
}
```

### Soft Deletes

```go
// Soft delete (sets DeletedAt)
database.DBConn.Delete(&user, 1)

// Include soft-deleted records
var users []models.User
database.DBConn.Unscoped().Find(&users)

// Restore soft-deleted record
database.DBConn.Model(&user).Update("deleted_at", nil)

// Permanent delete
database.DBConn.Unscoped().Delete(&user, 1)
```

## Best Practices

### 1. Always Check Errors

```go
// ❌ Bad
database.DBConn.Create(&user)

// ✅ Good
if err := database.DBConn.Create(&user).Error; err != nil {
    return err
}
```

### 2. Use Preload for Relationships

```go
// ❌ Bad (N+1 queries)
var users []models.User
database.DBConn.Find(&users)
for _, user := range users {
    database.DBConn.Where("user_id = ?", user.ID).Find(&user.Posts)
}

// ✅ Good (2 queries)
var users []models.User
database.DBConn.Preload("Posts").Find(&users)
```

### 3. Use Indexes

```sql
-- Add indexes for frequently queried columns
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);
```

```go
// In model
type User struct {
    gorm.Model
    Email string `gorm:"index"`
}
```

### 4. Validate Before Saving

```go
func CreateUser(user *models.User) error {
    // Validate
    if len(user.Username) < 3 {
        return errors.New("username too short")
    }

    // Save
    return database.DBConn.Create(user).Error
}
```

### 5. Use Scopes for Reusable Queries

```go
// Define scope
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("deleted_at IS NULL")
}

func RecentUsers(db *gorm.DB) *gorm.DB {
    return db.Where("created_at > ?", time.Now().AddDate(0, -1, 0))
}

// Use scopes
var users []models.User
database.DBConn.Scopes(ActiveUsers, RecentUsers).Find(&users)
```

## Troubleshooting

### Connection Refused

```
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```

**Solutions:**
1. Start MySQL: `sudo systemctl start mysql`
2. Check port: `netstat -an | grep 3306`
3. Verify credentials in `.env`

### Table Doesn't Exist

```
Error: Table 'fiberstarter.users' doesn't exist
```

**Solution:** Run migrations:
```bash
mysql -u root -p fiberstarter < migrations/m_00001_Up.sql
```

### Incorrect String Value

```
Error: Incorrect string value: '\xF0\x9F...'
```

**Solution:** Use `utf8mb4` charset:
```sql
ALTER DATABASE fiberstarter CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE users CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### Deadlock Found

```
Error: Deadlock found when trying to get lock
```

**Solution:** Use transactions and proper locking:
```go
database.DBConn.Transaction(func(tx *gorm.DB) error {
    tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user)
    // ... operations ...
    return nil
})
```

---

For more information, see:
- [GORM Documentation](https://gorm.io/docs/)
- [Getting Started Guide](GETTING_STARTED.md)
- [Adding Features](ADDING_FEATURES.md)
