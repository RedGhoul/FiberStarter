# Adding Features to FiberStarter

This guide walks you through adding new features to FiberStarter with a complete, hands-on example.

## Table of Contents

- [Feature Development Workflow](#feature-development-workflow)
- [Example: Building a Blog Feature](#example-building-a-blog-feature)
- [Step 1: Plan the Feature](#step-1-plan-the-feature)
- [Step 2: Create the Model](#step-2-create-the-model)
- [Step 3: Create the Migration](#step-3-create-the-migration)
- [Step 4: Create the Repository](#step-4-create-the-repository)
- [Step 5: Create the Controller](#step-5-create-the-controller)
- [Step 6: Create the Views](#step-6-create-the-views)
- [Step 7: Add the Routes](#step-7-add-the-routes)
- [Step 8: Test the Feature](#step-8-test-the-feature)
- [Advanced Patterns](#advanced-patterns)

## Feature Development Workflow

The standard workflow for adding features:

```
1. Plan → 2. Model → 3. Migration → 4. Repository → 5. Controller → 6. Views → 7. Routes → 8. Test
```

### Quick Checklist

- [ ] Define data model (struct)
- [ ] Create database migration
- [ ] Write repository functions (CRUD)
- [ ] Build controller handlers
- [ ] Design view templates
- [ ] Register routes
- [ ] Test functionality

## Example: Building a Blog Feature

Let's build a complete blog feature with:
- Posts (title, content, author)
- List all posts
- View single post
- Create new post (authenticated users only)
- Edit post (owner only)
- Delete post (owner only)

## Step 1: Plan the Feature

### Database Schema

```
posts table:
- id (primary key)
- user_id (foreign key to users)
- title (string, required)
- slug (string, unique, URL-friendly)
- content (text, required)
- published (boolean, default false)
- created_at (timestamp)
- updated_at (timestamp)
- deleted_at (timestamp, for soft delete)
```

### Routes Planning

```
GET  /posts              → List all published posts
GET  /posts/:slug        → View single post
GET  /posts/create       → Show create form (auth required)
POST /posts/create       → Create post (auth required)
GET  /posts/:id/edit     → Show edit form (owner only)
POST /posts/:id/edit     → Update post (owner only)
POST /posts/:id/delete   → Delete post (owner only)
```

## Step 2: Create the Model

Create `models/post.go`:

```go
package models

import (
    "gorm.io/gorm"
)

type Post struct {
    gorm.Model
    UserID    uint   `gorm:"not null;index"`
    Title     string `gorm:"size:255;not null"`
    Slug      string `gorm:"size:255;not null;uniqueIndex"`
    Content   string `gorm:"type:text;not null"`
    Published bool   `gorm:"default:false;index"`

    // Relationship
    User User `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name
func (Post) TableName() string {
    return "posts"
}
```

**Key Points:**
- Embeds `gorm.Model` for ID, timestamps, soft delete
- `UserID` links to the User model
- `Slug` for URL-friendly identifiers
- `Published` flag for draft/published state
- Defines relationship to User

## Step 3: Create the Migration

Create `migrations/m_00002_Up.sql`:

```sql
CREATE TABLE IF NOT EXISTS `posts` (
  `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `user_id` bigint unsigned NOT NULL,
  `title` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `published` tinyint(1) DEFAULT 0,
  UNIQUE KEY `idx_posts_slug` (`slug`),
  KEY `idx_posts_deleted_at` (`deleted_at`),
  KEY `idx_posts_user_id` (`user_id`),
  KEY `idx_posts_published` (`published`),
  CONSTRAINT `fk_posts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

Create `migrations/m_00002_Down.sql`:

```sql
DROP TABLE IF EXISTS `posts`;
```

**Run the migration:**

```bash
mysql -u root -p fiberstarter < migrations/m_00002_Up.sql
```

**Verify:**

```bash
mysql -u root -p fiberstarter -e "DESCRIBE posts;"
```

## Step 4: Create the Repository

Create `repos/post.go`:

```go
package repos

import (
    "fiberstarter/database"
    "fiberstarter/models"
    "errors"
)

// Get all published posts
func GetAllPublishedPosts() []models.Post {
    var posts []models.Post
    database.DBConn.
        Preload("User").
        Where("published = ?", true).
        Order("created_at DESC").
        Find(&posts)
    return posts
}

// Get all posts by user
func GetPostsByUserID(userID uint) []models.Post {
    var posts []models.Post
    database.DBConn.
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Find(&posts)
    return posts
}

// Get post by slug
func GetPostBySlug(slug string) (models.Post, error) {
    var post models.Post
    result := database.DBConn.
        Preload("User").
        Where("slug = ? AND published = ?", slug, true).
        First(&post)

    if result.Error != nil {
        return post, result.Error
    }

    return post, nil
}

// Get post by ID (for editing)
func GetPostByID(id uint) (models.Post, error) {
    var post models.Post
    result := database.DBConn.
        Preload("User").
        First(&post, id)

    if result.Error != nil {
        return post, result.Error
    }

    return post, nil
}

// Create new post
func CreatePost(post *models.Post) error {
    result := database.DBConn.Create(post)
    return result.Error
}

// Update post
func UpdatePost(post *models.Post) error {
    result := database.DBConn.Save(post)
    return result.Error
}

// Delete post (soft delete)
func DeletePost(id uint, userID uint) error {
    // Verify ownership
    var post models.Post
    result := database.DBConn.First(&post, id)

    if result.Error != nil {
        return result.Error
    }

    if post.UserID != userID {
        return errors.New("unauthorized")
    }

    // Soft delete
    return database.DBConn.Delete(&post).Error
}

// Check if slug exists
func SlugExists(slug string) bool {
    var post models.Post
    result := database.DBConn.Where("slug = ?", slug).First(&post)
    return result.RowsAffected > 0
}
```

**Key Functions:**
- `GetAllPublishedPosts()` - Public post list
- `GetPostBySlug()` - View single post
- `GetPostsByUserID()` - User's posts
- `CreatePost()` - Create new post
- `UpdatePost()` - Update existing post
- `DeletePost()` - Delete with ownership check

## Step 5: Create the Controller

Create `controllers/Post.go`:

```go
package controllers

import (
    "fiberstarter/models"
    "fiberstarter/repos"
    "fiberstarter/utils"
    "github.com/gofiber/fiber/v2"
    "github.com/gosimple/slug"
    "strconv"
)

// List all published posts
func ListPosts(c *fiber.Ctx) error {
    posts := repos.GetAllPublishedPosts()

    return c.Render("Posts/index", fiber.Map{
        "posts": posts,
    })
}

// View single post
func ViewPost(c *fiber.Ctx) error {
    postSlug := c.Params("slug")

    post, err := repos.GetPostBySlug(postSlug)
    if err != nil {
        return c.Status(404).Render("Common/soft_error", fiber.Map{
            "msg": "Post not found",
        })
    }

    return c.Render("Posts/view", fiber.Map{
        "post": post,
    })
}

// Show create form
func ShowCreatePost(c *fiber.Ctx) error {
    return c.Render("Posts/create", fiber.Map{})
}

// Handle create post
func CreatePost(c *fiber.Ctx) error {
    // Get current user
    user, err := utils.GetCurrentUser(c)
    if err != nil {
        return c.Redirect("/Login")
    }

    // Get form data
    title := c.FormValue("title")
    content := c.FormValue("content")
    published := c.FormValue("published") == "on"

    // Validate
    if title == "" || content == "" {
        return c.Render("Posts/create", fiber.Map{
            "error": "Title and content are required",
        })
    }

    // Generate slug
    postSlug := slug.Make(title)

    // Ensure slug is unique
    originalSlug := postSlug
    counter := 1
    for repos.SlugExists(postSlug) {
        postSlug = originalSlug + "-" + strconv.Itoa(counter)
        counter++
    }

    // Create post
    post := models.Post{
        UserID:    user.ID,
        Title:     title,
        Slug:      postSlug,
        Content:   content,
        Published: published,
    }

    err = repos.CreatePost(&post)
    if err != nil {
        return c.Render("Posts/create", fiber.Map{
            "error": "Failed to create post",
        })
    }

    return c.Redirect("/posts")
}

// Show edit form
func ShowEditPost(c *fiber.Ctx) error {
    user, err := utils.GetCurrentUser(c)
    if err != nil {
        return c.Redirect("/Login")
    }

    postID, _ := strconv.Atoi(c.Params("id"))
    post, err := repos.GetPostByID(uint(postID))

    if err != nil {
        return c.Status(404).Render("Common/soft_error", fiber.Map{
            "msg": "Post not found",
        })
    }

    // Check ownership
    if post.UserID != user.ID {
        return c.Status(403).Render("Common/soft_error", fiber.Map{
            "msg": "Unauthorized",
        })
    }

    return c.Render("Posts/edit", fiber.Map{
        "post": post,
    })
}

// Handle edit post
func EditPost(c *fiber.Ctx) error {
    user, err := utils.GetCurrentUser(c)
    if err != nil {
        return c.Redirect("/Login")
    }

    postID, _ := strconv.Atoi(c.Params("id"))
    post, err := repos.GetPostByID(uint(postID))

    if err != nil {
        return c.Redirect("/posts")
    }

    // Check ownership
    if post.UserID != user.ID {
        return c.Status(403).SendString("Unauthorized")
    }

    // Get form data
    title := c.FormValue("title")
    content := c.FormValue("content")
    published := c.FormValue("published") == "on"

    // Update post
    post.Title = title
    post.Content = content
    post.Published = published

    // Regenerate slug if title changed
    newSlug := slug.Make(title)
    if newSlug != post.Slug {
        originalSlug := newSlug
        counter := 1
        for repos.SlugExists(newSlug) {
            newSlug = originalSlug + "-" + strconv.Itoa(counter)
            counter++
        }
        post.Slug = newSlug
    }

    err = repos.UpdatePost(&post)
    if err != nil {
        return c.Render("Posts/edit", fiber.Map{
            "post":  post,
            "error": "Failed to update post",
        })
    }

    return c.Redirect("/posts/" + post.Slug)
}

// Delete post
func DeletePost(c *fiber.Ctx) error {
    user, err := utils.GetCurrentUser(c)
    if err != nil {
        return c.Redirect("/Login")
    }

    postID, _ := strconv.Atoi(c.Params("id"))

    err = repos.DeletePost(uint(postID), user.ID)
    if err != nil {
        return c.Status(403).SendString("Unauthorized")
    }

    return c.Redirect("/posts")
}
```

**Install slug package:**

```bash
go get github.com/gosimple/slug
```

## Step 6: Create the Views

### Create directory:

```bash
mkdir -p views/Posts
```

### 6.1 List Posts (`views/Posts/index.django`)

```django
{% extends "layouts/main.django" %}

{% block content %}
<div class="container mt-4">
    <div class="d-flex justify-content-between align-items-center mb-4">
        <h1>Blog Posts</h1>
        {% if LoggedIn %}
        <a href="/posts/create" class="btn btn-primary">Create Post</a>
        {% endif %}
    </div>

    {% if posts %}
    <div class="row">
        {% for post in posts %}
        <div class="col-md-6 mb-4">
            <div class="card">
                <div class="card-body">
                    <h5 class="card-title">{{ post.Title }}</h5>
                    <p class="card-text">{{ post.Content|slice:":150" }}...</p>
                    <p class="text-muted small">By {{ post.User.Username }} on {{ post.CreatedAt }}</p>
                    <a href="/posts/{{ post.Slug }}" class="btn btn-sm btn-primary">Read More</a>
                </div>
            </div>
        </div>
        {% endfor %}
    </div>
    {% else %}
    <p class="text-muted">No posts yet. Be the first to create one!</p>
    {% endif %}
</div>
{% endblock %}
```

### 6.2 View Post (`views/Posts/view.django`)

```django
{% extends "layouts/main.django" %}

{% block content %}
<div class="container mt-4">
    <article>
        <h1>{{ post.Title }}</h1>
        <p class="text-muted">
            By {{ post.User.Username }} | {{ post.CreatedAt }}
        </p>
        <hr>
        <div class="post-content">
            {{ post.Content }}
        </div>
    </article>

    {% if LoggedIn %}
    <div class="mt-4">
        <a href="/posts/{{ post.ID }}/edit" class="btn btn-warning">Edit</a>
        <form method="POST" action="/posts/{{ post.ID }}/delete" style="display:inline">
            <button type="submit" class="btn btn-danger" onclick="return confirm('Delete this post?')">Delete</button>
        </form>
    </div>
    {% endif %}

    <div class="mt-4">
        <a href="/posts" class="btn btn-secondary">Back to Posts</a>
    </div>
</div>
{% endblock %}
```

### 6.3 Create Post (`views/Posts/create.django`)

```django
{% extends "layouts/main.django" %}

{% block content %}
<div class="container mt-4">
    <h1>Create New Post</h1>

    {% if error %}
    <div class="alert alert-danger">{{ error }}</div>
    {% endif %}

    <form method="POST" action="/posts/create">
        <div class="mb-3">
            <label for="title" class="form-label">Title</label>
            <input type="text" class="form-control" id="title" name="title" required>
        </div>

        <div class="mb-3">
            <label for="content" class="form-label">Content</label>
            <textarea class="form-control" id="content" name="content" rows="10" required></textarea>
        </div>

        <div class="mb-3 form-check">
            <input type="checkbox" class="form-check-input" id="published" name="published">
            <label class="form-check-label" for="published">Publish immediately</label>
        </div>

        <button type="submit" class="btn btn-primary">Create Post</button>
        <a href="/posts" class="btn btn-secondary">Cancel</a>
    </form>
</div>
{% endblock %}
```

### 6.4 Edit Post (`views/Posts/edit.django`)

```django
{% extends "layouts/main.django" %}

{% block content %}
<div class="container mt-4">
    <h1>Edit Post</h1>

    {% if error %}
    <div class="alert alert-danger">{{ error }}</div>
    {% endif %}

    <form method="POST" action="/posts/{{ post.ID }}/edit">
        <div class="mb-3">
            <label for="title" class="form-label">Title</label>
            <input type="text" class="form-control" id="title" name="title" value="{{ post.Title }}" required>
        </div>

        <div class="mb-3">
            <label for="content" class="form-label">Content</label>
            <textarea class="form-control" id="content" name="content" rows="10" required>{{ post.Content }}</textarea>
        </div>

        <div class="mb-3 form-check">
            <input type="checkbox" class="form-check-input" id="published" name="published" {% if post.Published %}checked{% endif %}>
            <label class="form-check-label" for="published">Published</label>
        </div>

        <button type="submit" class="btn btn-primary">Update Post</button>
        <a href="/posts/{{ post.Slug }}" class="btn btn-secondary">Cancel</a>
    </form>
</div>
{% endblock %}
```

## Step 7: Add the Routes

Edit `routes/routes.go`:

```go
package routes

import (
    "fiberstarter/controllers"
    "fiberstarter/utils"
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    // Existing routes
    setupAuthRoutes(app)
    setupBasicRoutes(app)

    // Add locals middleware
    app.Use(utils.AddLocals())

    // Public blog routes
    app.Get("/posts", controllers.ListPosts)
    app.Get("/posts/:slug", controllers.ViewPost)

    // Protected routes
    app.Use(utils.CheckAuth())

    setupHiddenRoutes(app)

    // Protected blog routes
    app.Get("/posts/create", controllers.ShowCreatePost)
    app.Post("/posts/create", controllers.CreatePost)
    app.Get("/posts/:id/edit", controllers.ShowEditPost)
    app.Post("/posts/:id/edit", controllers.EditPost)
    app.Post("/posts/:id/delete", controllers.DeletePost)
}
```

### Update Navigation

Edit `views/partials/header.django`:

```django
<nav class="navbar navbar-expand-lg navbar-light bg-light">
    <div class="container-fluid">
        <a class="navbar-brand" href="/">{{ APPNAME }}</a>
        <button class="navbar-toggler" ...>...</button>
        <div class="collapse navbar-collapse">
            <ul class="navbar-nav me-auto mb-2 mb-lg-0">
                <li class="nav-item">
                    <a class="nav-link" href="/">Home</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="/posts">Blog</a>
                </li>
                {% if LoggedIn %}
                <li class="nav-item">
                    <a class="nav-link" href="/Secret">Secret</a>
                </li>
                {% endif %}
            </ul>
            <ul class="navbar-nav">
                {% if LoggedIn %}
                <li class="nav-item">
                    <a class="nav-link" href="/Logout">Logout</a>
                </li>
                {% else %}
                <li class="nav-item">
                    <a class="nav-link" href="/Login">Login</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="/Register">Register</a>
                </li>
                {% endif %}
            </ul>
        </div>
    </div>
</nav>
```

## Step 8: Test the Feature

### 1. Run the application

```bash
go run main.go
```

### 2. Test the workflow

**Create a post:**
1. Login at `/Login`
2. Visit `/posts/create`
3. Fill out the form
4. Click "Create Post"

**View posts:**
1. Visit `/posts`
2. See your post listed
3. Click "Read More"

**Edit post:**
1. Click "Edit" on post detail page
2. Update the content
3. Click "Update Post"

**Delete post:**
1. Click "Delete" on post detail page
2. Confirm deletion

### 3. Test edge cases

- [ ] Create post without login → Redirects to login
- [ ] Edit someone else's post → Shows unauthorized
- [ ] Delete someone else's post → Shows unauthorized
- [ ] View draft post → Not visible in list
- [ ] Create post with duplicate title → Unique slug generated

## Advanced Patterns

### Adding Validation

Create `utils/validation.go`:

```go
package utils

import "errors"

func ValidatePost(title, content string) error {
    if len(title) < 3 {
        return errors.New("title must be at least 3 characters")
    }

    if len(title) > 255 {
        return errors.New("title must be less than 255 characters")
    }

    if len(content) < 10 {
        return errors.New("content must be at least 10 characters")
    }

    return nil
}
```

**Use in controller:**

```go
err := utils.ValidatePost(title, content)
if err != nil {
    return c.Render("Posts/create", fiber.Map{
        "error": err.Error(),
    })
}
```

### Adding File Uploads

```go
func CreatePost(c *fiber.Ctx) error {
    // ... existing code ...

    // Handle file upload
    file, err := c.FormFile("image")
    if err == nil {
        filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
        c.SaveFile(file, "./public/uploads/"+filename)
        post.ImagePath = "/uploads/" + filename
    }

    // ... create post ...
}
```

### Adding Rich Text Editor

In `views/Posts/create.django`:

```html
<link href="https://cdn.quilljs.com/1.3.6/quill.snow.css" rel="stylesheet">
<script src="https://cdn.quilljs.com/1.3.6/quill.js"></script>

<div id="editor"></div>
<textarea name="content" id="content" style="display:none"></textarea>

<script>
var quill = new Quill('#editor', { theme: 'snow' });
quill.on('text-change', function() {
    document.getElementById('content').value = quill.root.innerHTML;
});
</script>
```

---

**Congratulations!** You've built a complete blog feature for FiberStarter. Use this pattern to build any feature you need.

For more information, see:
- [Database Guide](DATABASE.md)
- [Authentication Guide](AUTHENTICATION.md)
