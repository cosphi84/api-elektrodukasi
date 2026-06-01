# Elektrodukasi API - CRUD Operations Documentation

This document outlines all available API endpoints for the Elektrodukasi platform.

## Base URL
```
http://localhost:8080/api
```

## Authentication
⚠️ **Note**: Authentication middleware should be implemented based on your requirements.

---

## Users API

### Create User
```http
POST /api/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "hashed_password",
  "role": "user"
}
```

**Response**: `201 Created`
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "is_active": true,
  "role": "user",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### List Users
```http
GET /api/users?page=1&page_size=10
```

**Response**: `200 OK`
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "page_size": 10
}
```

### Get User by ID
```http
GET /api/users/{id}
```

### Update User
```http
PUT /api/users/{id}
Content-Type: application/json

{
  "name": "Jane Doe",
  "avatar": "https://example.com/avatar.jpg"
}
```

### Deactivate User
```http
PATCH /api/users/{id}/deactivate
```

### Delete User
```http
DELETE /api/users/{id}
```

---

## Categories API

### Create Category
```http
POST /api/categories
Content-Type: application/json

{
  "name": "Programming",
  "description": "Programming tutorials and articles"
}
```

### List Categories
```http
GET /api/categories?page=1&page_size=10
```

### Get Category by ID
```http
GET /api/categories/{id}
```

### Update Category
```http
PUT /api/categories/{id}
Content-Type: application/json

{
  "name": "Advanced Programming",
  "description": "Advanced programming topics"
}
```

### Delete Category
```http
DELETE /api/categories/{id}
```

---

## Tags API

### Create Tag
```http
POST /api/tags
Content-Type: application/json

{
  "name": "javascript"
}
```

### List Tags
```http
GET /api/tags?page=1&page_size=10
```

### Get Tag by ID
```http
GET /api/tags/{id}
```

### Update Tag
```http
PUT /api/tags/{id}
Content-Type: application/json

{
  "name": "typescript"
}
```

### Delete Tag
```http
DELETE /api/tags/{id}
```

---

## Articles API

### Create Article
```http
POST /api/articles
Content-Type: application/json

{
  "author_id": "uuid",
  "category_id": "uuid",
  "title": "Getting Started with Go",
  "slug": "getting-started-with-go",
  "summary": "Learn the basics of Go programming",
  "content_html": "<p>HTML content here</p>",
  "content_json": {/* tiptap JSON */},
  "image": "https://example.com/image.jpg"
}
```

### List Articles
```http
GET /api/articles?page=1&page_size=10&published=true
```

**Query Parameters**:
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 10)
- `published`: Filter by published status (true/false)

### Get Article by ID
```http
GET /api/articles/{id}
```
✨ **Note**: Increments view count automatically

### Get Article by Slug
```http
GET /api/articles/slug/{slug}
```

### Get Articles by Category
```http
GET /api/categories/{categoryId}/articles?page=1&page_size=10
```

### Update Article
```http
PUT /api/articles/{id}
Content-Type: application/json

{
  "title": "Advanced Go Programming",
  "summary": "Master advanced Go concepts"
}
```

### Publish Article
```http
PATCH /api/articles/{id}/publish
```

### Add Tag to Article
```http
POST /api/articles/{id}/tags
Content-Type: application/json

{
  "tag_id": "uuid"
}
```

### Remove Tag from Article
```http
DELETE /api/articles/{id}/tags/{tagId}
```

### Delete Article
```http
DELETE /api/articles/{id}
```

---

## Comments API

### Create Comment
```http
POST /api/comments
Content-Type: application/json

{
  "article_id": "uuid",
  "user_id": "uuid",
  "parent_id": null,
  "content": "Great article!"
}
```

### Get Comment by ID
```http
GET /api/comments/{id}
```

### Get Comments by Article
```http
GET /api/articles/{articleId}/comments?page=1&page_size=10&approved=false
```

**Query Parameters**:
- `approved`: Filter by approval status (true/false)

### Update Comment
```http
PUT /api/comments/{id}
Content-Type: application/json

{
  "content": "Updated comment"
}
```

### Approve Comment
```http
PATCH /api/comments/{id}/approve
```

### Get Comment Replies
```http
GET /api/comments/{id}/replies
```

### Delete Comment
```http
DELETE /api/comments/{id}
```

---

## Projects API

### Create Project
```http
POST /api/projects
Content-Type: application/json

{
  "owner_id": "uuid",
  "title": "My Awesome Project",
  "slug": "my-awesome-project",
  "summary": "A brief description",
  "link": "https://github.com/user/project",
  "metadata": {/* JSON data */}
}
```

### List Projects
```http
GET /api/projects?page=1&page_size=10
```

### Get Project by ID
```http
GET /api/projects/{id}
```

### Get Project by Slug
```http
GET /api/projects/slug/{slug}
```

### Get Projects by Owner
```http
GET /api/users/{ownerId}/projects?page=1&page_size=10
```

### Update Project
```http
PUT /api/projects/{id}
Content-Type: application/json

{
  "title": "Updated Project Title",
  "link": "https://github.com/user/updated-project"
}
```

### Delete Project
```http
DELETE /api/projects/{id}
```

---

## Response Codes

| Code | Description |
|------|-------------|
| 200 | OK - Successful GET/PUT/PATCH request |
| 201 | Created - Successful POST request |
| 204 | No Content - Successful DELETE request |
| 400 | Bad Request - Invalid input |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error |

---

## Error Response Format

```json
{
  "error": "error message"
}
```

---

## Environment Variables

```bash
DATABASE_URL=host=localhost user=postgres password=postgres dbname=elektrodukasi port=5432 sslmode=disable
PORT=8080
```

---

## Notes

- All UUIDs should be valid UUID format (e.g., `550e8400-e29b-41d4-a716-446655440000`)
- Timestamps are in ISO 8601 format
- Article view count is incremented automatically when fetched by ID or slug
- Comments support nested replies via `parent_id`
- Soft deletes are not implemented by default (use DELETE for hard deletes)
