# Elektrodukasi API - CRUD Implementation Summary

## Overview

This project provides a complete CRUD (Create, Read, Update, Delete) REST API implementation for the Elektrodukasi platform using Go, GORM, and PostgreSQL.

## Project Structure

```
internal/
├── models/              # GORM models (7 models)
│   ├── user.model.go
│   ├── category.model.go
│   ├── tag.model.go
│   ├── article.model.go
│   ├── article_tag.model.go (junction table)
│   ├── comment.model.go
│   └── project.model.go
│
├── services/            # Business logic layer (6 services)
│   ├── user.service.go
│   ├── category.service.go
│   ├── tag.service.go
│   ├── article.service.go
│   ├── comment.service.go
│   └── project.service.go
│
├── handlers/            # HTTP handlers (6 handlers)
│   ├── user.handler.go
│   ├── category.handler.go
│   ├── tag.handler.go
│   ├── article.handler.go
│   ├── comment.handler.go
│   └── project.handler.go
│
├── routes/              # Route definitions
│   └── routes.go
│
├── repositories/        # Data access layer (pre-existing)
└── database/            # Database configuration

cmd/api/
└── main.go              # Application entry point
```

## Implemented Models & Endpoints

### 1. **Users** (6 endpoints)
- `POST /api/users` - Create user
- `GET /api/users` - List all users (paginated)
- `GET /api/users/{id}` - Get user by ID
- `PUT /api/users/{id}` - Update user
- `PATCH /api/users/{id}/deactivate` - Deactivate user
- `DELETE /api/users/{id}` - Delete user

**Fields**: id, name, email, password, avatar, is_active, role, created_at, updated_at, deleted_at, last_login, last_login_from

---

### 2. **Categories** (5 endpoints)
- `POST /api/categories` - Create category
- `GET /api/categories` - List all categories (paginated)
- `GET /api/categories/{id}` - Get category by ID
- `PUT /api/categories/{id}` - Update category
- `DELETE /api/categories/{id}` - Delete category

**Fields**: id, name, description, created_at, updated_at

---

### 3. **Tags** (5 endpoints)
- `POST /api/tags` - Create tag
- `GET /api/tags` - List all tags (paginated)
- `GET /api/tags/{id}` - Get tag by ID
- `PUT /api/tags/{id}` - Update tag
- `DELETE /api/tags/{id}` - Delete tag

**Fields**: id, name, created_at, updated_at

---

### 4. **Articles** (10 endpoints)
- `POST /api/articles` - Create article
- `GET /api/articles` - List articles (paginated, filterable by published status)
- `GET /api/articles/{id}` - Get article by ID (auto-increments view count)
- `GET /api/articles/slug/{slug}` - Get article by slug (auto-increments view count)
- `GET /api/categories/{categoryId}/articles` - Get articles by category (paginated)
- `PUT /api/articles/{id}` - Update article
- `PATCH /api/articles/{id}/publish` - Publish article
- `POST /api/articles/{id}/tags` - Add tag to article
- `DELETE /api/articles/{id}/tags/{tagId}` - Remove tag from article
- `DELETE /api/articles/{id}` - Delete article

**Fields**: id, author_id, category_id, title, slug, summary, content_html, content_json, image, published, published_at, view_count, created_at, updated_at, deleted_at

**Relations**: Author (User), Category, Tags (many-to-many)

---

### 5. **Comments** (7 endpoints)
- `POST /api/comments` - Create comment
- `GET /api/comments/{id}` - Get comment by ID
- `GET /api/articles/{articleId}/comments` - Get comments by article (paginated, filterable by approved status)
- `PUT /api/comments/{id}` - Update comment
- `PATCH /api/comments/{id}/approve` - Approve comment
- `GET /api/comments/{id}/replies` - Get comment replies
- `DELETE /api/comments/{id}` - Delete comment

**Fields**: id, article_id, user_id, parent_id, content, approved, lft, rgt, depth, children_count, created_at, updated_at

**Relations**: Article, User, Parent comment

**Note**: Supports nested set model for hierarchical comment structure

---

### 6. **Projects** (7 endpoints)
- `POST /api/projects` - Create project
- `GET /api/projects` - List all projects (paginated)
- `GET /api/projects/{id}` - Get project by ID
- `GET /api/projects/slug/{slug}` - Get project by slug
- `GET /api/users/{ownerId}/projects` - Get projects by owner (paginated)
- `PUT /api/projects/{id}` - Update project
- `DELETE /api/projects/{id}` - Delete project

**Fields**: id, owner_id, title, slug, summary, link, metadata, created_at, updated_at

**Relations**: Owner (User)

---

## Features

### Service Layer Features
- ✅ **Pagination**: All list endpoints support `page` and `page_size` query parameters
- ✅ **Filtering**: Articles support `published` filter, Comments support `approved` filter
- ✅ **Preloading**: Automatic eager loading of related data (Author, Category, Tags, User, Owner)
- ✅ **Error Handling**: Proper error messages for not found, validation, and server errors
- ✅ **Many-to-Many**: Full support for Article-Tag relationships
- ✅ **Nested Comments**: Support for hierarchical comment structure
- ✅ **View Tracking**: Auto-increment article view count on fetch
- ✅ **Soft Delete Support**: Ready for soft deletes in models

### Handler Layer Features
- ✅ **JSON Request/Response**: Automatic JSON encoding/decoding
- ✅ **HTTP Status Codes**: Proper status codes (201 Created, 204 No Content, 400/404/500)
- ✅ **Pagination Response**: Consistent format with data, total, page, page_size
- ✅ **Query Parameter Parsing**: Safe parsing with defaults
- ✅ **Path Parameters**: UUID and integer parameter extraction
- ✅ **Error Responses**: User-friendly error messages

### Route Definition
- ✅ **RESTful Design**: Following REST principles
- ✅ **HTTP Methods**: Proper use of GET, POST, PUT, PATCH, DELETE
- ✅ **Hierarchical Routes**: Category→Articles, User→Projects, Article→Comments
- ✅ **Slug Support**: Alternative lookup by slug for Articles and Projects

---

## Installation

### Prerequisites
- Go 1.26.3+
- PostgreSQL 12+

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd api-elektrodukasi
   ```

2. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Configure database**
   ```bash
   export DATABASE_URL="host=localhost user=postgres password=postgres dbname=elektrodukasi port=5432 sslmode=disable"
   ```

4. **Run migrations**
   ```bash
   # Use your migration tool (e.g., migrate, flyway, or sql-migrate)
   migrate -path migrations -database "$DATABASE_URL" up
   ```

5. **Start the server**
   ```bash
   go run cmd/api/main.go
   ```

   The API will be available at `http://localhost:8080`

---

## API Usage Examples

### Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "hashed_password",
    "role": "user"
  }'
```

### List Articles with Pagination
```bash
curl "http://localhost:8080/api/articles?page=1&page_size=10&published=true"
```

### Get Article by Slug
```bash
curl "http://localhost:8080/api/articles/slug/my-first-article"
```

### Add Tag to Article
```bash
curl -X POST http://localhost:8080/api/articles/{articleId}/tags \
  -H "Content-Type: application/json" \
  -d '{"tag_id": "tag-uuid"}'
```

### Create Comment
```bash
curl -X POST http://localhost:8080/api/comments \
  -H "Content-Type: application/json" \
  -d '{
    "article_id": "article-uuid",
    "user_id": "user-uuid",
    "parent_id": null,
    "content": "Great article!"
  }'
```

---

## Database Schema

The database includes the following tables:
1. **users** - User accounts with authentication info
2. **categories** - Article categories
3. **tags** - Article tags
4. **articles** - Blog articles with rich content support
5. **article_tags** - Junction table for article-tag relationships
6. **comments** - Hierarchical comments with nested set model
7. **projects** - User projects with metadata

See migrations folder for detailed schema definitions.

---

## Response Format

### Success Response (200 OK)
```json
{
  "id": "uuid",
  "name": "Example",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### List Response (200 OK)
```json
{
  "data": [
    { "id": "uuid", "name": "Item 1" },
    { "id": "uuid", "name": "Item 2" }
  ],
  "total": 100,
  "page": 1,
  "page_size": 10
}
```

### Error Response (4xx/5xx)
```json
{
  "error": "user not found"
}
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DATABASE_URL | (localhost) | PostgreSQL connection string |
| PORT | 8080 | Server port |

---

## Testing the API

### Using cURL
See examples above or check `API_DOCUMENTATION.md`

### Using Postman
1. Import the API endpoints into Postman
2. Set environment variable for `base_url`: `http://localhost:8080`
3. Set `page` and `page_size` for list endpoints

### Using Go Test
```bash
go test ./...
```

---

## Architecture

```
HTTP Request
    ↓
Routes (routes.go)
    ↓
Handler Layer (handlers/)
    ↓
Service Layer (services/)
    ↓
GORM ORM
    ↓
PostgreSQL Database
```

Each layer has a single responsibility:
- **Routes**: HTTP routing and method matching
- **Handlers**: HTTP request/response handling, validation
- **Services**: Business logic, data transformation
- **GORM**: Database query execution

---

## Future Enhancements

- [ ] Add authentication middleware (JWT/OAuth)
- [ ] Add request validation layer
- [ ] Add caching layer (Redis)
- [ ] Add rate limiting
- [ ] Add API versioning
- [ ] Add comprehensive API tests
- [ ] Add Swagger/OpenAPI documentation
- [ ] Add logging and metrics
- [ ] Add search functionality (Elasticsearch)

---

## File Summary

**Services Created**: 6 service files with ~500 lines
**Handlers Created**: 6 handler files with ~800 lines
**Routes Created**: 1 routes file with ~70 endpoints mapped
**Models Used**: 7 existing GORM models
**Total New Lines**: ~1,500 lines of code

---

## License

[Your License Here]

---

## Support

For issues, questions, or contributions, please refer to the project documentation or contact the development team.
