# Getting Started with FiberStarter

This guide will walk you through setting up FiberStarter from scratch and building your first feature.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24 or higher** - [Download Go](https://golang.org/dl/)
- **MySQL 8.0 or higher** - [Download MySQL](https://dev.mysql.com/downloads/)
- **Git** - [Download Git](https://git-scm.com/downloads)
- **A code editor** - We recommend [VSCode](https://code.visualstudio.com/) with the Go extension

### Verifying Your Installation

```bash
# Check Go version
go version
# Expected: go version go1.24.0 or higher

# Check MySQL version
mysql --version
# Expected: mysql Ver 8.0.x or higher

# Check Git version
git --version
# Expected: git version 2.x.x or higher
```

## Installation Steps

### 1. Clone the Repository

```bash
git clone <your-repository-url>
cd FiberStarter
```

### 2. Install Go Dependencies

```bash
go mod download
```

This will download all required packages including:
- Fiber v2 (web framework)
- GORM (ORM)
- MySQL driver
- Session storage
- Hashing library
- And more...

### 3. Configure Your Environment

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

If `.env.example` doesn't exist, create `.env` with the following content:

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
DB_PASSWORD=your_mysql_password
DB_NAME=fiberstarter
DB_HOST=localhost
DB_PORT=3306
```

**Important**: Update `DB_PASSWORD` with your actual MySQL password.

### 4. Set Up the Database

#### Create the Database

```bash
# Login to MySQL
mysql -u root -p

# Create the database
CREATE DATABASE fiberstarter CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# Exit MySQL
exit;
```

Alternatively, use a single command:

```bash
mysql -u root -p -e "CREATE DATABASE fiberstarter CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

#### Run Migrations

```bash
# Run the initial migration
mysql -u root -p fiberstarter < migrations/m_00001_Up.sql
```

This will create the following tables:
- `users` - User accounts with authentication
- `books` - Example resource table
- `session_storage` - Session management (auto-created)

### 5. Run the Application

```bash
go run main.go
```

You should see output similar to:

```
 ┌───────────────────────────────────────────────────┐
 │                   Fiber v2.52.5                   │
 │               http://127.0.0.1:8080               │
 │       (bound on host 0.0.0.0 and port 8080)       │
 │                                                   │
 │ Handlers ............ 15  Processes ........... 1 │
 │ Prefork ....... Disabled  PID ............. 12345 │
 └───────────────────────────────────────────────────┘
```

### 6. Test the Application

Open your browser and visit: `http://localhost:8080`

You should see the home page with the FiberStarter navbar.

## Your First Steps

### 1. Register a New User

1. Click **Register** in the navbar
2. Fill out the form:
   - Username: `testuser`
   - Email: `test@example.com`
   - Password: `password123`
   - Confirm Password: `password123`
3. Click **Register**
4. You'll be redirected to the login page

### 2. Login

1. Enter your credentials:
   - Username: `testuser`
   - Password: `password123`
2. Click **Login**
3. You'll be redirected to the home page
4. Notice the navbar now shows **Logout** instead of **Login/Register**

### 3. Access Protected Routes

1. Click **Secret** in the navbar (only visible when logged in)
2. You'll see the protected page with the message: "I like dogs"
3. Try accessing `/Secret` in an incognito window - you'll be redirected to login

### 4. Logout

1. Click **Logout** in the navbar
2. You'll be redirected to the home page
3. Your session will be cleared

## Development Workflow

### Hot Reload with Air

For a better development experience, use Air for automatic reloading:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

Now any changes to `.go` or `.django` files will automatically reload the server.

### VSCode Debugging

1. Open the project in VSCode
2. Press `F5` or click **Run and Debug**
3. Set breakpoints in your code
4. Navigate to pages in your browser
5. Debugger will pause at breakpoints

### Project Structure Overview

```
FiberStarter/
├── main.go                 # Entry point - calls StartUp()
├── app.go                  # Application setup
│
├── controllers/            # Request handlers
│   ├── Auth.go            # Login, Register, Logout
│   └── Index.go           # Home, Secret pages
│
├── routes/                 # Route definitions
│   └── routes.go          # Maps URLs to controllers
│
├── middleware/             # Middleware components
│   └── middleware.go      # CSRF, sessions, logging, etc.
│
├── models/                 # Database models
│   ├── user.go            # User struct
│   └── book.go            # Book struct
│
├── repos/                  # Data access layer
│   ├── user.go            # User database operations
│   └── book.go            # Book database operations
│
├── database/               # Database connection
│   └── database.go        # GORM setup
│
├── providers/              # Service providers
│   ├── config.go          # Environment config
│   ├── session.go         # Session management
│   ├── auth.go            # Authentication checks
│   └── hash.go            # Password hashing
│
├── utils/                  # Helper functions
│   └── utils.go           # Auth, cookies, etc.
│
├── views/                  # Django templates
│   ├── layouts/           # Base templates
│   ├── partials/          # Reusable components
│   ├── Auth/              # Login, Register forms
│   ├── Home/              # Home, Secret pages
│   └── Common/            # Error pages
│
└── public/                 # Static files
    └── site.js            # Frontend JavaScript
```

## Understanding the Request Flow

Let's trace a request through the application:

### Example: Accessing a Protected Page

```
User visits: http://localhost:8080/Secret

1. Request arrives at Fiber router
   ↓
2. Middleware chain executes (in order):
   - RequestID: Assigns unique ID
   - Logger: Logs request
   - AddLocals: Sets LoggedIn, APPNAME variables
   - CSRF: Validates token (if POST)
   - CheckAuth: Verifies user session
   ↓
3. Controller handler executes:
   - Index.ShowSecret() in controllers/Index.go
   - Prepares data for view
   ↓
4. View rendering:
   - Renders views/Home/secret.django
   - Wrapped in views/layouts/main.django
   - Includes views/partials/header.django
   ↓
5. Response sent to browser
   - HTML with Bootstrap styling
   - CSRF token in meta tags
   - Session cookie set
```

## Common Tasks

### View Application Logs

Logs are printed to stdout with the format:

```
[PID] [RequestID] [Status] - [Method] [Path]
```

Example:
```
12345 abc123def456 200 - GET /
12345 xyz789ghi012 302 - POST /Login
```

### Access MySQL Database

```bash
# Connect to database
mysql -u root -p fiberstarter

# View users
SELECT * FROM users;

# View sessions
SELECT * FROM session_storage;

# Exit
exit;
```

### Clear All Sessions

```sql
DELETE FROM session_storage;
```

### Reset Database

```bash
# Drop and recreate database
mysql -u root -p -e "DROP DATABASE fiberstarter; CREATE DATABASE fiberstarter;"

# Run migrations again
mysql -u root -p fiberstarter < migrations/m_00001_Up.sql
```

## Environment Configuration Details

### Application Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_NAME` | FiberStarter | Displayed in navbar |
| `SERVER_PORT` | 8080 | Port to listen on |
| `ENABLE_DEBUG` | true | Enable debug logging |
| `ENABLE_PREFORK` | false | Enable Fiber prefork mode |
| `ENABLE_COMPRESSION` | true | Enable gzip compression |
| `ENABLE_RECOVERY` | true | Enable panic recovery |
| `CONCURRENCY_LEVEL` | 200 | Max concurrent connections |

### Database Settings

| Variable | Required | Description |
|----------|----------|-------------|
| `DB_USER` | Yes | MySQL username |
| `DB_PASSWORD` | Yes | MySQL password |
| `DB_NAME` | Yes | Database name |
| `DB_HOST` | Yes | Database host |
| `DB_PORT` | Yes | Database port |

## Troubleshooting

### Port Already in Use

```
Error: listen tcp :8080: bind: address already in use
```

**Solution**: Change `SERVER_PORT` in `.env` or kill the process using port 8080:

```bash
# Find process
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Database Connection Failed

```
Error: Error 1045: Access denied for user 'root'@'localhost'
```

**Solutions**:
1. Check `DB_PASSWORD` in `.env`
2. Verify MySQL is running: `sudo systemctl status mysql`
3. Test connection: `mysql -u root -p`

### Migration Failed

```
Error: Table 'users' already exists
```

**Solution**: Drop existing tables first:

```bash
mysql -u root -p fiberstarter < migrations/m_00001_Down.sql
mysql -u root -p fiberstarter < migrations/m_00001_Up.sql
```

### Template Not Found

```
Error: template: views/Home/index.django: file not found
```

**Solution**: Check that:
1. File exists at the path
2. Extension is `.django`
3. Path is relative to `views/` directory

## Next Steps

Now that you have FiberStarter running, explore these guides:

- **[Authentication Guide](AUTHENTICATION.md)** - Learn how authentication works
- **[Database & Migrations](DATABASE.md)** - Working with the database
- **[Adding Features](ADDING_FEATURES.md)** - Build your first feature
- **[Deployment](DEPLOYMENT.md)** - Deploy to production

## Quick Reference Commands

```bash
# Development
go run main.go              # Run application
air                         # Run with hot reload
go test ./...               # Run tests

# Database
mysql -u root -p fiberstarter                    # Connect to database
mysql -u root -p fiberstarter < migrations/*.sql # Run migration

# Dependencies
go mod download             # Install dependencies
go mod tidy                 # Clean up dependencies
go get <package>            # Add new dependency

# Build
go build -o app main.go     # Build binary
./app                       # Run binary
```

## Getting Help

If you encounter issues:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review the other documentation guides
3. Open an issue on the repository
4. Check Fiber documentation: https://docs.gofiber.io/

---

**Congratulations!** You now have a fully functional FiberStarter application running. Happy coding!
