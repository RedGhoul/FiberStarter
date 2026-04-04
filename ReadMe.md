# FiberStarter

A modern, production-ready Go web application framework built on Fiber v2 with server-side rendering, authentication, and MySQL database integration.

## Features

- **🚀 Fast & Modern**: Built on Go Fiber v2 for high performance
- **🔐 Authentication System**: Complete user registration, login, and session management
- **🗄️ Database Ready**: MySQL integration with GORM and migration support
- **🎨 Server-Side Rendering**: Django template engine with Bootstrap 5.2
- **🛡️ Security**: CSRF protection, Argon2id password hashing, secure sessions
- **📦 Clean Architecture**: Repository pattern with clear separation of concerns
- **🔥 Hot Reload**: Development-ready with Air for instant feedback
- **⚡ Middleware**: Logging, compression, recovery, and more out of the box

## Tech Stack

### Core Framework
- **Go 1.24+** - Programming language
- **Fiber v2** - Fast HTTP web framework for Go
- **GORM** - Feature-rich ORM for database operations

### Frontend
- **Django Template Language** - Server-side templating
- **Bootstrap 5.2** - Responsive CSS framework
- **Alpine.js** - Lightweight JavaScript framework

### Database & Storage
- **MySQL 8.0+** - Primary database
- **MySQL Session Storage** - Secure session management

### Security
- **fiber-hashing** - Argon2id password hashing
- **CSRF Protection** - Built-in token validation
- **Secure Sessions** - MySQL-backed session storage

## Quick Start

### Option A: Docker (Recommended)

The fastest way to get started -- no need to install MySQL locally:

```bash
git clone <repository-url>
cd FiberStarter
docker-compose up --build
```

Visit `http://localhost:8080` once both services are healthy.

### Option B: Local Setup

#### Prerequisites

- Go 1.24 or higher
- MySQL 8.0 or higher
- Git

#### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd FiberStarter
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Set up database**
   ```bash
   # Create the database
   mysql -u root -p -e "CREATE DATABASE fiberstarter;"

   # Run migrations
   mysql -u root -p fiberstarter < migrations/m_00001_Up.sql
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

6. **Visit the application**
   Open your browser to `http://localhost:8080`

## Project Structure

```
FiberStarter/
├── app.go                  # Application initialization and setup
├── main.go                 # Entry point
├── controllers/            # HTTP request handlers
│   ├── Auth.go            # Authentication logic
│   └── Index.go           # Home and protected routes
├── routes/                 # Route definitions
│   └── routes.go          # All route mappings
├── middleware/             # HTTP middleware
│   └── middleware.go      # CSRF, logging, sessions, etc.
├── models/                 # GORM database models
│   ├── user.go            # User model
│   └── book.go            # Example resource model
├── repos/                  # Data access layer
│   ├── user.go            # User repository
│   └── book.go            # Book repository
├── database/               # Database connection
│   └── database.go        # GORM initialization
├── providers/              # Service providers
│   ├── config.go          # Configuration management
│   ├── session.go         # Session provider
│   ├── auth.go            # Authentication provider
│   └── hash.go            # Password hashing provider
├── utils/                  # Utility functions
│   └── utils.go           # Helpers for auth, cookies, etc.
├── migrations/             # SQL migrations
│   ├── m_00001_Up.sql     # Initial schema
│   └── m_00001_Down.sql   # Rollback scripts
├── views/                  # Django template files
│   ├── layouts/           # Base templates
│   ├── partials/          # Reusable components
│   ├── Auth/              # Authentication views
│   ├── Home/              # Home page views
│   └── Common/            # Common views (errors, etc.)
└── public/                 # Static assets
    └── site.js            # Frontend JavaScript
```

## Configuration

Configure your application using environment variables in `.env`:

```env
# Application Settings
APP_NAME=FiberStarter
SERVER_PORT=8080
ENABLE_DEBUG=true
ENABLE_PREFORK=false
ENABLE_COMPRESSION=true
ENABLE_RECOVERY=true
CONCURRENCY_LEVEL=200

# Database Configuration
DB_USER=root
DB_PASSWORD=password
DB_NAME=fiberstarter
DB_HOST=localhost
DB_PORT=3306
```

## Core Features

### Authentication
- User registration with email validation
- Secure login with Argon2id password hashing
- Session-based authentication
- Protected routes with middleware
- Logout functionality

### Database
- GORM ORM with MySQL
- Repository pattern for data access
- SQL-based migrations
- Soft delete support
- Connection pooling

### Middleware Stack
1. Request ID tracking
2. Structured logging
3. Static file serving
4. CSRF protection
5. Session management
6. Panic recovery
7. Gzip compression
8. Custom 404 handling

### Views & Templates
- Django template engine
- Layout inheritance
- Partial templates
- Bootstrap 5.2 styling
- Alpine.js for interactivity
- CSRF token injection

## Documentation

For detailed guides, see:

- **[Getting Started Guide](docs/GETTING_STARTED.md)** - Detailed setup and first steps
- **[Authentication Guide](docs/AUTHENTICATION.md)** - Working with user authentication
- **[Database & Migrations](docs/DATABASE.md)** - Database operations and migrations
- **[Creating Features](docs/ADDING_FEATURES.md)** - Step-by-step guide to add new features
- **[Deployment](docs/DEPLOYMENT.md)** - Production deployment guide

## Development

### Hot Reload

For development with automatic reloading:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Running Tests

```bash
go test ./...
```

### VSCode Debugging

A launch configuration is included in `.vscode/launch.json`. Press F5 to start debugging.

## Routes

### Public Routes
- `GET /` - Home page
- `GET /health` - Health check (JSON)
- `GET /Login` - Login form
- `POST /Login` - Process login
- `GET /Register` - Registration form
- `POST /Register` - Process registration

### Protected Routes (Require Authentication)
- `GET /Secret` - Example protected page
- `GET /Logout` - Logout and clear session

## Usage Examples

### Adding a New Page

Create a controller, a view, and register the route.

**1. Add the controller** (`controllers/Products.go`):

```go
package controllers

import "github.com/gofiber/fiber/v2"

func ShowProducts(c *fiber.Ctx) error {
    return c.Render("Products/index", fiber.Map{
        "Title": "Our Products",
    })
}
```

**2. Add the view** (`views/Products/index.django`):

```django
{% include "partials/header.django" %}
<div class="container py-3">
    <h1>{{ Title }}</h1>
    <p>Product list goes here.</p>
</div>
{% include "partials/footer.django" %}
```

**3. Register the route** in `routes/routes.go`:

```go
// Inside setupBasicRoutes():
app.Get("/Products", controllers.ShowProducts)
```

### Adding a New Model and Database Table

**1. Define the model** (`models/product.go`):

```go
package models

import "time"

type Product struct {
    ID        uint      `json:"Id"`
    CreatedAt time.Time `json:"CreatedAt"`
    UpdatedAt time.Time `json:"UpdatedAt"`
    DeletedAt time.Time `json:"DeletedAt"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
}
```

**2. Create a migration** (`migrations/m_00002_Up.sql`):

```sql
CREATE TABLE IF NOT EXISTS `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**3. Add a repository** (`repos/product.go`):

```go
package repos

import (
    "FiberStarter/database"
    "FiberStarter/models"
)

func GetAllProducts() []models.Product {
    var products []models.Product
    database.DBConn.Find(&products)
    return products
}

func GetProductByID(id int) models.Product {
    var product models.Product
    database.DBConn.First(&product, id)
    return product
}

func CreateProduct(name string, price float64) bool {
    product := models.Product{Name: name, Price: price}
    database.DBConn.Create(&product)
    return product.ID > 0
}
```

### Adding a Protected Route

Any route registered **after** the `app.Use(utils.CheckAuth())` middleware is automatically protected. Unauthenticated users are redirected to `/Login`.

```go
// In routes/routes.go:
func SetupRoutes(app *fiber.App) {
    // ...public routes...

    app.Use(utils.CheckAuth()) // Auth barrier

    // Everything below requires login
    app.Get("/Dashboard", controllers.ShowDashboard)
    app.Get("/Profile", controllers.ShowProfile)
}
```

### Adding a JSON API Endpoint

You can mix server-rendered pages with JSON endpoints:

```go
// controllers/Api.go
package controllers

import (
    "FiberStarter/repos"
    "github.com/gofiber/fiber/v2"
)

func GetProductsJSON(c *fiber.Ctx) error {
    products := repos.GetAllProducts()
    return c.JSON(products)
}
```

```go
// In routes/routes.go, add before the auth middleware:
app.Get("/api/products", controllers.GetProductsJSON)
```

Test it with curl:

```bash
curl http://localhost:8080/api/products
```

### Using the Health Check

The built-in health endpoint reports application and database status:

```bash
curl http://localhost:8080/health
# {"database":"ok","status":"ok"}
```

### Using the Makefile

```bash
make help          # Show all available commands
make run           # Start the application
make build         # Build a binary to bin/fiberstarter
make test          # Run all tests
make dev           # Hot reload with Air
make docker-up     # Start with Docker Compose
make docker-down   # Stop Docker Compose
```

## Security Features

- **Password Security**: Argon2id hashing (industry standard)
- **Session Security**: MySQL-backed sessions with secure cookies
- **CSRF Protection**: Token validation on all POST requests
- **SQL Injection Prevention**: GORM parameterized queries
- **Access Control**: Middleware-based authentication checks

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the MIT License.

## Support

For issues, questions, or contributions, please open an issue on the repository.

---

Built with ❤️ using Go Fiber
