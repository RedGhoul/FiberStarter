# Authentication Guide

This guide explains how authentication works in FiberStarter, including user registration, login, session management, and protecting routes.

## Table of Contents

- [Overview](#overview)
- [Authentication Flow](#authentication-flow)
- [User Registration](#user-registration)
- [User Login](#user-login)
- [Session Management](#session-management)
- [Password Security](#password-security)
- [Protecting Routes](#protecting-routes)
- [Working with Sessions](#working-with-sessions)
- [Customizing Authentication](#customizing-authentication)

## Overview

FiberStarter uses a **session-based authentication system** with the following components:

- **MySQL Session Storage** - Sessions persisted in database
- **Argon2id Password Hashing** - Industry-standard password security
- **CSRF Protection** - Token-based request validation
- **Middleware-based Access Control** - Route protection
- **Secure Cookies** - HttpOnly, SameSite cookies

### Key Files

- `controllers/Auth.go` - Authentication handlers (login, register, logout)
- `providers/auth.go` - Authentication provider (`IsAuthenticated()`)
- `providers/hash.go` - Password hashing provider
- `providers/session.go` - Session store provider
- `utils/utils.go` - Auth utilities (cookies, session helpers)
- `repos/user.go` - User database operations
- `models/user.go` - User model definition

## Authentication Flow

### Registration Flow

```
1. User visits /Register
   ↓
2. ShowRegisterForm() renders register.django
   ↓
3. User submits form → POST /Register
   ↓
4. PostRegisterForm() validates:
   - Passwords match
   - Email not already registered
   ↓
5. Password hashed with Argon2id
   ↓
6. User saved to database
   ↓
7. Redirect to /Login
```

### Login Flow

```
1. User visits /Login
   ↓
2. ShowLoginForm() renders login.django
   ↓
3. User submits credentials → POST /Login
   ↓
4. PostLoginForm() validates:
   - User exists (repos.GetUserByUsername)
   - Password correct (MatchPasswords)
   ↓
5. Session created with userid
   ↓
6. Session cookie set
   ↓
7. Redirect to /
```

### Protected Route Flow

```
1. User visits /Secret
   ↓
2. CheckAuth() middleware executes
   ↓
3. providers.IsAuthenticated() checks session
   ↓
4. If authenticated: Continue to controller
5. If not authenticated: Redirect to /Login
```

## User Registration

### Registration Form (`views/Auth/register.django`)

```html
<form method="POST" action="/Register">
    <input type="text" name="username" required>
    <input type="email" name="email" required>
    <input type="password" name="password1" required>
    <input type="password" name="password2" required>
    <button type="submit">Register</button>
</form>
```

### Registration Handler (`controllers/Auth.go`)

```go
func PostRegisterForm(c *fiber.Ctx) error {
    username := c.FormValue("username")
    email := c.FormValue("email")
    password1 := c.FormValue("password1")
    password2 := c.FormValue("password2")

    // Validate passwords match
    if password1 != password2 {
        return c.Redirect("/Register")
    }

    // Check if user exists
    if repos.CheckIfUserExists(email) {
        return c.Redirect("/Register")
    }

    // Hash password
    hashedPassword, _ := providers.HashProvider().Hash(password1)

    // Create user
    newUser := models.User{
        Username: username,
        Email:    email,
        Password: hashedPassword,
    }
    repos.CreateUser(&newUser)

    return c.Redirect("/Login")
}
```

### Key Validations

1. **Password Matching** - `password1 == password2`
2. **Email Uniqueness** - Check database before creation
3. **Password Hashing** - Never store plain text passwords

### Customizing Registration

**Add Username Validation:**

```go
// Check username length
if len(username) < 3 {
    return c.Render("Auth/register", fiber.Map{
        "error": "Username must be at least 3 characters",
    })
}

// Check username uniqueness
if repos.CheckIfUsernameExists(username) {
    return c.Render("Auth/register", fiber.Map{
        "error": "Username already taken",
    })
}
```

**Add Email Validation:**

```go
import "net/mail"

// Validate email format
_, err := mail.ParseAddress(email)
if err != nil {
    return c.Render("Auth/register", fiber.Map{
        "error": "Invalid email address",
    })
}
```

**Add Password Requirements:**

```go
// Require minimum password length
if len(password1) < 8 {
    return c.Render("Auth/register", fiber.Map{
        "error": "Password must be at least 8 characters",
    })
}

// Require complexity (example)
import "regexp"

hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password1)
hasLower := regexp.MustCompile(`[a-z]`).MatchString(password1)
hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password1)

if !hasUpper || !hasLower || !hasNumber {
    return c.Render("Auth/register", fiber.Map{
        "error": "Password must contain uppercase, lowercase, and numbers",
    })
}
```

## User Login

### Login Form (`views/Auth/login.django`)

```html
<form method="POST" action="/Login">
    <input type="text" name="username" required>
    <input type="password" name="password" required>
    <button type="submit">Login</button>
</form>
```

### Login Handler (`controllers/Auth.go`)

```go
func PostLoginForm(c *fiber.Ctx) error {
    username := c.FormValue("username")
    password := c.FormValue("password")

    // Get user from database
    curUser := repos.GetUserByUsername(username)

    // Validate password
    if !utils.MatchPasswords(password, curUser.Password) {
        return c.Redirect("/Login")
    }

    // Set auth cookie (creates session)
    utils.SetAuthCookie(c, curUser.ID)

    return c.Redirect("/")
}
```

### Password Verification (`utils/utils.go`)

```go
func MatchPasswords(password string, hashedPassword string) bool {
    if len(hashedPassword) < 5 {
        return false
    }

    matched, err := providers.HashProvider().MatchHash(
        password,
        hashedPassword,
    )

    if err != nil {
        return false
    }

    return matched
}
```

### Setting the Session Cookie (`utils/utils.go`)

```go
func SetAuthCookie(c *fiber.Ctx, userId uint) error {
    store, err := providers.SessionProvider().Get(c)
    if err != nil {
        return err
    }

    // Store user ID in session
    store.Set("userid", userId)

    // Save to database
    return store.Save()
}
```

## Session Management

### Session Storage

Sessions are stored in MySQL using the `gofiber/storage/mysql` package:

```go
// providers/session.go
import "github.com/gofiber/storage/mysql/v2"

func SetUpSessionProvider(url string) {
    storage := mysql.New(mysql.Config{
        ConnectionURI: url,
        GCInterval:    time.Second * 10,
    })

    sessionP = session.New(session.Config{
        Storage: storage,
    })
}
```

### Session Table

The `session_storage` table is auto-created:

```sql
CREATE TABLE session_storage (
    k VARCHAR(64) PRIMARY KEY,    -- Session key (cookie value)
    v BLOB,                        -- Session data (serialized)
    e BIGINT                       -- Expiration timestamp
);
```

### Session Data Structure

Sessions store a map of key-value pairs:

```go
// What's stored
{
    "userid": uint(1)  // User ID
}
```

### Checking Authentication (`providers/auth.go`)

```go
func IsAuthenticated(c *fiber.Ctx) bool {
    store, err := SessionProvider().Get(c)
    if err != nil {
        return false
    }

    userid := store.Get("userid")
    return userid != nil
}
```

### Getting Current User

```go
// utils/utils.go
func GetCurrentUser(c *fiber.Ctx) (*models.User, error) {
    store, err := providers.SessionProvider().Get(c)
    if err != nil {
        return nil, err
    }

    userid := store.Get("userid")
    if userid == nil {
        return nil, errors.New("not authenticated")
    }

    user := repos.GetUserByID(userid.(uint))
    return &user, nil
}
```

**Usage in Controller:**

```go
func ShowProfile(c *fiber.Ctx) error {
    user, err := utils.GetCurrentUser(c)
    if err != nil {
        return c.Redirect("/Login")
    }

    return c.Render("Auth/profile", fiber.Map{
        "user": user,
    })
}
```

## Password Security

### Hashing Algorithm

FiberStarter uses **Argon2id** through the `fiber-hashing` library:

```go
// providers/hash.go
import "github.com/thomasvvugt/fiber-hashing"

var hashP *hashing.Argon2id

func SetUpHashProvider() {
    hashP = hashing.NewArgon2id()
}

func HashProvider() *hashing.Argon2id {
    return hashP
}
```

### Why Argon2id?

- **Industry Standard**: Winner of Password Hashing Competition (2015)
- **Memory-Hard**: Resistant to GPU/ASIC attacks
- **Configurable**: Tunable time/memory costs
- **Hybrid Mode**: Combines Argon2i and Argon2d benefits

### Hashing a Password

```go
hashedPassword, err := providers.HashProvider().Hash("plaintext")
// Result: $argon2id$v=19$m=65536,t=1,p=2$...
```

### Verifying a Password

```go
matched, err := providers.HashProvider().MatchHash(
    "plaintext",
    "$argon2id$v=19$m=65536,t=1,p=2$...",
)
```

### Password Storage

Passwords are stored in the `users` table:

```go
// models/user.go
type User struct {
    gorm.Model
    Username string
    Email    string
    Password string `gorm:"size:600"` // Large enough for hash
}
```

**Important**: The `Password` field must be at least 600 characters to store Argon2id hashes.

## Protecting Routes

### Middleware-Based Protection

Protected routes use the `CheckAuth()` middleware:

```go
// routes/routes.go
func SetupRoutes(app *fiber.App) {
    // Public routes (no authentication)
    setupAuthRoutes(app)
    setupBasicRoutes(app)

    // Apply authentication middleware
    app.Use(utils.CheckAuth())

    // Protected routes (require authentication)
    setupHiddenRoutes(app)
}
```

### CheckAuth Middleware (`utils/utils.go`)

```go
func CheckAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        if !providers.IsAuthenticated(c) {
            return c.Redirect("/Login")
        }
        return c.Next()
    }
}
```

### Route Groups

```go
// Public routes
func setupAuthRoutes(app *fiber.App) {
    app.Get("/Login", controllers.ShowLoginForm)
    app.Post("/Login", controllers.PostLoginForm)
    app.Get("/Register", controllers.ShowRegisterForm)
    app.Post("/Register", controllers.PostRegisterForm)
}

// Protected routes
func setupHiddenRoutes(app *fiber.App) {
    app.Get("/Secret", controllers.ShowSecret)
    app.Get("/Logout", controllers.PostLogoutForm)
}
```

### Inline Protection (Alternative)

You can also check authentication in individual handlers:

```go
func ShowProfile(c *fiber.Ctx) error {
    if !providers.IsAuthenticated(c) {
        return c.Redirect("/Login")
    }

    // Handler logic
    return c.Render("profile", fiber.Map{})
}
```

## Working with Sessions

### Storing Data in Sessions

```go
func SaveToSession(c *fiber.Ctx, key string, value interface{}) error {
    store, err := providers.SessionProvider().Get(c)
    if err != nil {
        return err
    }

    store.Set(key, value)
    return store.Save()
}
```

### Retrieving Data from Sessions

```go
func GetFromSession(c *fiber.Ctx, key string) (interface{}, error) {
    store, err := providers.SessionProvider().Get(c)
    if err != nil {
        return nil, err
    }

    return store.Get(key), nil
}
```

### Deleting Session Data

```go
func RemoveFromSession(c *fiber.Ctx, key string) error {
    store, err := providers.SessionProvider().Get(c)
    if err != nil {
        return err
    }

    store.Delete(key)
    return store.Save()
}
```

### Flash Messages (Session-based)

```go
// Set flash message
func SetFlash(c *fiber.Ctx, message string) error {
    store, _ := providers.SessionProvider().Get(c)
    store.Set("flash", message)
    return store.Save()
}

// Get and clear flash message
func GetFlash(c *fiber.Ctx) string {
    store, _ := providers.SessionProvider().Get(c)
    flash := store.Get("flash")

    if flash != nil {
        store.Delete("flash")
        store.Save()
        return flash.(string)
    }

    return ""
}
```

**Usage:**

```go
// In controller
func PostLogin(c *fiber.Ctx) error {
    // ... login logic ...
    utils.SetFlash(c, "Login successful!")
    return c.Redirect("/")
}

// In middleware (make available to views)
func AddFlashToLocals() fiber.Handler {
    return func(c *fiber.Ctx) error {
        c.Locals("Flash", utils.GetFlash(c))
        return c.Next()
    }
}
```

**In template:**

```html
{% if Flash %}
    <div class="alert alert-success">{{ Flash }}</div>
{% endif %}
```

## Customizing Authentication

### Adding Remember Me

```go
func PostLoginForm(c *fiber.Ctx) error {
    rememberMe := c.FormValue("remember_me") == "on"

    // ... authentication logic ...

    if rememberMe {
        store, _ := providers.SessionProvider().Get(c)
        store.SetExpiry(time.Hour * 24 * 30) // 30 days
        store.Save()
    }

    return c.Redirect("/")
}
```

### Adding Role-Based Access Control

**Update User Model:**

```go
type User struct {
    gorm.Model
    Username string
    Email    string
    Password string `gorm:"size:600"`
    Role     string `gorm:"default:'user'"` // user, admin, etc.
}
```

**Create Role Middleware:**

```go
// utils/utils.go
func RequireRole(role string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := GetCurrentUser(c)
        if err != nil {
            return c.Redirect("/Login")
        }

        if user.Role != role {
            return c.Status(403).SendString("Forbidden")
        }

        return c.Next()
    }
}
```

**Usage:**

```go
// routes/routes.go
admin := app.Group("/admin", utils.RequireRole("admin"))
admin.Get("/dashboard", controllers.AdminDashboard)
```

### Adding OAuth/Social Login

See our separate guide: **[OAuth Integration](OAUTH.md)** (coming soon)

## Logout

### Logout Handler (`controllers/Auth.go`)

```go
func PostLogoutForm(c *fiber.Ctx) error {
    utils.RemoveCookie(c)
    return c.Redirect("/")
}
```

### RemoveCookie Function (`utils/utils.go`)

```go
func RemoveCookie(c *fiber.Ctx) error {
    store, err := providers.SessionProvider().Get(c)
    if err != nil {
        return err
    }

    // Clear userid from session
    store.Delete("userid")
    store.Save()

    // Clear cookie from browser
    c.ClearCookie()

    return nil
}
```

## Security Best Practices

### 1. Always Use HTTPS in Production

```go
// app.go
if !providers.AppConfig.Enable_Debug {
    app.Use(func(c *fiber.Ctx) error {
        if c.Protocol() != "https" {
            return c.Redirect("https://" + c.Hostname() + c.OriginalURL())
        }
        return c.Next()
    })
}
```

### 2. Set Secure Cookie Attributes

```go
// providers/session.go
sessionP = session.New(session.Config{
    Storage: storage,
    CookieSecure: true,          // HTTPS only
    CookieHTTPOnly: true,         // No JavaScript access
    CookieSameSite: "Strict",     // CSRF protection
})
```

### 3. Implement Rate Limiting

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

// middleware/middleware.go
app.Use("/Login", limiter.New(limiter.Config{
    Max:        5,
    Expiration: 1 * time.Minute,
}))
```

### 4. Log Failed Login Attempts

```go
// controllers/Auth.go
func PostLoginForm(c *fiber.Ctx) error {
    // ... authentication logic ...

    if !utils.MatchPasswords(password, curUser.Password) {
        log.Printf("Failed login attempt for user: %s from IP: %s",
            username, c.IP())
        return c.Redirect("/Login")
    }

    // ...
}
```

## Testing Authentication

### Test User Registration

```go
// tests/auth_test.go
func TestUserRegistration(t *testing.T) {
    app := StartUp()

    req := httptest.NewRequest("POST", "/Register", strings.NewReader(
        "username=testuser&email=test@example.com&password1=pass123&password2=pass123",
    ))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, _ := app.Test(req)

    assert.Equal(t, 302, resp.StatusCode)
    assert.Equal(t, "/Login", resp.Header.Get("Location"))
}
```

## Troubleshooting

### Session Not Persisting

**Problem**: User logs in but immediately logged out.

**Solutions**:
1. Check MySQL session_storage table exists
2. Verify session GC interval: `GCInterval: time.Second * 10`
3. Ensure `store.Save()` is called after `store.Set()`

### Password Validation Always Fails

**Problem**: Correct password rejected.

**Solutions**:
1. Verify password field is at least 600 chars: `gorm:"size:600"`
2. Check hash is stored correctly (starts with `$argon2id$`)
3. Ensure no whitespace in password input

### Redirect Loop on Protected Routes

**Problem**: Accessing `/Secret` causes infinite redirect.

**Solutions**:
1. Verify `/Login` route is before `CheckAuth()` middleware
2. Check `IsAuthenticated()` logic
3. Ensure session cookie is being set

---

For more information, see:
- [Getting Started Guide](GETTING_STARTED.md)
- [Database Guide](DATABASE.md)
- [Adding Features](ADDING_FEATURES.md)
