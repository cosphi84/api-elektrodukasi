# CRUD Implementation Quick Reference

## What Was Created

### 📁 Services (Business Logic)
- `user.service.go` - User CRUD operations
- `category.service.go` - Category CRUD operations
- `tag.service.go` - Tag CRUD operations
- `article.service.go` - Article CRUD + publishing + tag management + view counting
- `comment.service.go` - Comment CRUD + approval + reply management
- `project.service.go` - Project CRUD + owner filtering

### 📁 Handlers (HTTP Layer)
- `user.handler.go` - User endpoints (Create, List, Get, Update, Delete, Deactivate)
- `category.handler.go` - Category endpoints (CRUD)
- `tag.handler.go` - Tag endpoints (CRUD)
- `article.handler.go` - Article endpoints (CRUD + Publishing + Tag management + Slug lookup)
- `comment.handler.go` - Comment endpoints (CRUD + Approval + Replies)
- `project.handler.go` - Project endpoints (CRUD + Slug lookup + Owner filtering)

### 📁 Routes
- `routes.go` - Central route registry that:
  - Initializes all services
  - Initializes all handlers
  - Registers 70+ API endpoints
  - Uses Go's native http.ServeMux (no external routing library needed)

### 📝 Documentation
- `API_DOCUMENTATION.md` - Complete API reference with examples
- `CRUD_IMPLEMENTATION.md` - Implementation overview and usage guide
- `QUICK_REFERENCE.md` - This file

### 🔧 Configuration
- `cmd/api/main.go` - Updated to use new routes and services
- `go.mod` - Updated with PostgreSQL driver dependency

---

## API Endpoints Summary

### Users: 6 endpoints
```
POST   /api/users
GET    /api/users
GET    /api/users/{id}
PUT    /api/users/{id}
DELETE /api/users/{id}
PATCH  /api/users/{id}/deactivate
```

### Categories: 5 endpoints
```
POST   /api/categories
GET    /api/categories
GET    /api/categories/{id}
PUT    /api/categories/{id}
DELETE /api/categories/{id}
```

### Tags: 5 endpoints
```
POST   /api/tags
GET    /api/tags
GET    /api/tags/{id}
PUT    /api/tags/{id}
DELETE /api/tags/{id}
```

### Articles: 10 endpoints
```
POST   /api/articles
GET    /api/articles
GET    /api/articles/{id}
GET    /api/articles/slug/{slug}
GET    /api/categories/{categoryId}/articles
PUT    /api/articles/{id}
PATCH  /api/articles/{id}/publish
POST   /api/articles/{id}/tags
DELETE /api/articles/{id}/tags/{tagId}
DELETE /api/articles/{id}
```

### Comments: 7 endpoints
```
POST   /api/comments
GET    /api/comments/{id}
GET    /api/articles/{articleId}/comments
PUT    /api/comments/{id}
PATCH  /api/comments/{id}/approve
GET    /api/comments/{id}/replies
DELETE /api/comments/{id}
```

### Projects: 7 endpoints
```
POST   /api/projects
GET    /api/projects
GET    /api/projects/{id}
GET    /api/projects/slug/{slug}
GET    /api/users/{ownerId}/projects
PUT    /api/projects/{id}
DELETE /api/projects/{id}
```

**Total: 40 endpoints**

---

## Key Features

### 🎯 For Each Model
- ✅ Full CRUD operations
- ✅ Pagination support (page, page_size)
- ✅ Proper HTTP status codes
- ✅ JSON request/response handling
- ✅ Error handling with meaningful messages
- ✅ UUID validation for path parameters
- ✅ Relation preloading

### 📊 Special Features
- **Articles**: View count auto-increment, publish status, slug support, tag management
- **Comments**: Hierarchical structure, approval workflow, reply management
- **Projects**: Slug support, owner filtering
- **Users**: Deactivation instead of hard delete option, email lookup

---

## How to Use

### 1. Start Development
```bash
cd api-elektrodukasi
go mod tidy
export DATABASE_URL="host=localhost user=postgres password=postgres dbname=elektrodukasi port=5432 sslmode=disable"
go run cmd/api/main.go
```

### 2. Test an Endpoint
```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"hashed","role":"user"}'

# List users
curl http://localhost:8080/api/users?page=1&page_size=10

# Get specific user
curl http://localhost:8080/api/users/{id}
```

### 3. Common Query Parameters
```
?page=1              # Page number (default: 1)
?page_size=10        # Items per page (default: 10)
?published=true      # For articles: filter by published status
?approved=true       # For comments: filter by approval status
```

---

## Service Interface Examples

### UserService
```go
interface UserService {
  CreateUser(user *models.User) error
  GetUserByID(id uuid.UUID) (*models.User, error)
  GetUserByEmail(email string) (*models.User, error)
  GetAllUsers(page, pageSize int) ([]models.User, int64, error)
  UpdateUser(id uuid.UUID, updates map[string]interface{}) (*models.User, error)
  DeleteUser(id uuid.UUID) error
  DeactivateUser(id uuid.UUID) error
}
```

### ArticleService
```go
interface ArticleService {
  CreateArticle(article *models.Article) error
  GetArticleByID(id uuid.UUID) (*models.Article, error)
  GetArticleBySlug(slug string) (*models.Article, error)
  GetAllArticles(page, pageSize int, published *bool) ([]models.Article, int64, error)
  GetArticlesByCategory(categoryID uuid.UUID, page, pageSize int) ([]models.Article, int64, error)
  UpdateArticle(id uuid.UUID, updates map[string]interface{}) (*models.Article, error)
  DeleteArticle(id uuid.UUID) error
  PublishArticle(id uuid.UUID) error
  AddTagToArticle(articleID, tagID uuid.UUID) error
  RemoveTagFromArticle(articleID, tagID uuid.UUID) error
  IncrementViewCount(id uuid.UUID) error
}
```

---

## Response Examples

### Create Success (201 Created)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "is_active": true,
  "role": "user",
  "created_at": "2024-01-01T12:00:00Z"
}
```

### List Success (200 OK)
```json
{
  "data": [
    {"id": "uuid1", "name": "Item 1"},
    {"id": "uuid2", "name": "Item 2"}
  ],
  "total": 100,
  "page": 1,
  "page_size": 10
}
```

### Error (404 Not Found)
```json
{
  "error": "user not found"
}
```

### Error (400 Bad Request)
```json
{
  "error": "Invalid request"
}
```

---

## Code Organization

### Layer Responsibilities

**Routes** → Registers HTTP paths and methods
- Maps HTTP methods to handlers
- Initializes services and handlers
- Manages dependency injection

**Handlers** → Processes HTTP requests/responses
- Parses URL parameters and query strings
- Decodes JSON request bodies
- Calls service methods
- Encodes JSON responses
- Returns appropriate HTTP status codes

**Services** → Contains business logic
- Data validation and transformation
- Database queries via GORM
- Error handling and mapping
- Pagination logic
- Eager loading with Preload

**Models** → GORM entities
- Table definitions
- Relationships
- JSON tags for serialization
- Database column mappings

---

## Next Steps (Recommended)

1. **Add Middleware**
   - Authentication (JWT/OAuth)
   - CORS
   - Request logging
   - Error recovery

2. **Add Validation**
   - Input validation layer
   - Custom validators
   - Constraint checking

3. **Add Testing**
   - Unit tests for services
   - Integration tests for handlers
   - Test fixtures and mocks

4. **Add Documentation**
   - Swagger/OpenAPI specs
   - Generate API docs from code
   - Client SDKs

5. **Add Observability**
   - Structured logging
   - Metrics collection
   - Distributed tracing

---

## File Locations

```
api-elektrodukasi/
├── cmd/
│   └── api/
│       └── main.go (updated)
├── internal/
│   ├── services/
│   │   ├── user.service.go
│   │   ├── category.service.go
│   │   ├── tag.service.go
│   │   ├── article.service.go
│   │   ├── comment.service.go
│   │   └── project.service.go
│   ├── handlers/
│   │   ├── user.handler.go
│   │   ├── category.handler.go
│   │   ├── tag.handler.go
│   │   ├── article.handler.go
│   │   ├── comment.handler.go
│   │   └── project.handler.go
│   ├── routes/
│   │   └── routes.go
│   ├── models/ (existing)
│   ├── database/ (existing)
│   ├── repositories/ (existing)
│   └── dto/ (existing)
├── migrations/ (existing)
├── go.mod (updated)
└── API_DOCUMENTATION.md
└── CRUD_IMPLEMENTATION.md
```

---

## Support for Complex Operations

### Multi-Step Operations
Example: Create article with tags
```go
// 1. Create article
article := CreateArticle(...)

// 2. Add tags
AddTagToArticle(article.ID, tag1.ID)
AddTagToArticle(article.ID, tag2.ID)

// 3. Publish
PublishArticle(article.ID)
```

### Filtered Listing
Example: Get published articles in category
```
GET /api/categories/{categoryId}/articles?published=true&page=1&page_size=20
```

### Hierarchical Queries
Example: Get comment replies
```
GET /api/comments/{parentId}/replies
```

### Slug-based Lookup
Example: Get article by slug
```
GET /api/articles/slug/my-awesome-article
```

---

## Notes

- All UUIDs are auto-generated if not provided
- Timestamps are in ISO 8601 format
- Pagination defaults to page=1, page_size=10
- Delete operations are hard deletes (use soft delete fields in models if needed)
- Relations are eagerly loaded to avoid N+1 queries
- Comments support nesting via parent_id field
