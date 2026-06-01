# 🔐 Authentication & Authorization Documentation

## Overview

The Elektrodukasi API implements JWT-based authentication with role-based access control (RBAC).

### Features
- ✅ JWT token-based authentication (Bearer tokens)
- ✅ Token expiration in 24 hours
- ✅ Role-based access control (Admin, User)
- ✅ Public read access for articles, categories, tags, and comments
- ✅ Secure password hashing with bcrypt
- ✅ Login endpoint for user authentication

---

## 🎯 Roles & Permissions

### Admin Role
Can perform **ALL operations**:
- ✅ Create, Read, Update, Delete Users
- ✅ Create, Read, Update, Delete Categories
- ✅ Create, Read, Update, Delete Tags
- ✅ Create, Read, Update, Delete, Publish Articles
- ✅ Create, Read, Update, Delete, Approve Comments
- ✅ Create, Read, Update, Delete Projects

### User Role
Can perform **LIMITED operations**:
- ✅ Read all articles, categories, tags, projects (public data)
- ✅ Create, Read, Update, Delete Comments (only their own on articles)
- ❌ Cannot modify users, categories, tags, articles, or projects

---

## 🚀 Getting Started

### 1. Login to Get JWT Token

**Endpoint:** `POST /api/login`
**Authentication:** None (Public endpoint)

**Request:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Admin User",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "invalid email or password"
}
```

---

## 🔑 Using the JWT Token

### Add Token to Request Headers

All authenticated endpoints require the token in the `Authorization` header with `Bearer` scheme:

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:8080/api/users
```

### Token Format
```
Authorization: Bearer <your_jwt_token_here>
```

---

## 📋 Endpoint Protection

### Public Endpoints (No Authentication Required)

```
GET    /api/articles
GET    /api/articles/{id}
GET    /api/articles/slug/{slug}
GET    /api/categories
GET    /api/categories/{id}
GET    /api/tags
GET    /api/tags/{id}
GET    /api/projects
GET    /api/projects/{id}
GET    /api/projects/slug/{slug}
GET    /api/articles/{articleId}/comments
GET    /api/comments/{id}
GET    /api/comments/{id}/replies
POST   /api/login
```

### Admin Only Endpoints (Requires Admin Role)

```
POST   /api/users
GET    /api/users
GET    /api/users/{id}
PUT    /api/users/{id}
DELETE /api/users/{id}
PATCH  /api/users/{id}/deactivate

POST   /api/categories
PUT    /api/categories/{id}
DELETE /api/categories/{id}

POST   /api/tags
PUT    /api/tags/{id}
DELETE /api/tags/{id}

POST   /api/articles
PUT    /api/articles/{id}
PATCH  /api/articles/{id}/publish
POST   /api/articles/{id}/tags
DELETE /api/articles/{id}/tags/{tagId}
DELETE /api/articles/{id}

PATCH  /api/comments/{id}/approve

POST   /api/projects
PUT    /api/projects/{id}
DELETE /api/projects/{id}
```

### User & Admin Endpoints (Requires Authentication)

```
POST   /api/comments        (Create comment)
PUT    /api/comments/{id}   (Update comment)
DELETE /api/comments/{id}   (Delete comment)
GET    /api/users/{ownerId}/projects (View projects)
```

---

## 📝 Examples

### Example 1: Admin Creates a User

```bash
# 1. Login as admin
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "admin123"}'

# Response includes token
# Save token: TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 2. Create user with token
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John User",
    "email": "john@example.com",
    "password": "hashed_password",
    "role": "user"
  }'
```

### Example 2: Regular User Writes a Comment

```bash
# 1. Login as regular user
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "user123"}'

# Save token
# TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 2. Create comment
curl -X POST http://localhost:8080/api/comments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "article_id": "article-uuid",
    "user_id": "user-uuid",
    "content": "Great article!"
  }'

# 3. Try to create category (should fail - 403 Forbidden)
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "New Category"}'
# Response: {"error": "insufficient permissions for this operation"}
```

### Example 3: Read Public Data (No Token Required)

```bash
# These work without authentication
curl http://localhost:8080/api/articles
curl http://localhost:8080/api/articles/550e8400-e29b-41d4-a716-446655440000
curl http://localhost:8080/api/categories
curl http://localhost:8080/api/articles/article-uuid/comments
```

---

## 🔒 Security Features

### Token Security
- ✅ JWT signed with HMAC SHA-256
- ✅ Tokens expire after 24 hours
- ✅ Token validated on every protected request
- ✅ Invalid or expired tokens rejected with 401 Unauthorized

### Password Security
- ✅ Passwords hashed with bcrypt
- ✅ Plain passwords never stored
- ✅ Password comparison is timing-safe

### Error Messages
- ✅ Generic error messages ("invalid email or password")
- ✅ No user enumeration attacks
- ✅ Proper HTTP status codes

---

## 🛡️ Error Responses

### Missing Authorization Header
```
Status: 401 Unauthorized

{
  "error": "authorization header not provided"
}
```

### Invalid Token Format
```
Status: 401 Unauthorized

{
  "error": "invalid authorization header format"
}
```

### Expired or Invalid Token
```
Status: 401 Unauthorized

{
  "error": "invalid or expired token"
}
```

### Insufficient Permissions
```
Status: 403 Forbidden

{
  "error": "insufficient permissions for this operation"
}
```

### Invalid Credentials
```
Status: 401 Unauthorized

{
  "error": "invalid email or password"
}
```

### Account Deactivated
```
Status: 401 Unauthorized

{
  "error": "user account is deactivated"
}
```

---

## 🔑 JWT Token Structure

A JWT token consists of three parts separated by dots: `header.payload.signature`

### Header
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

### Payload (Claims)
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "admin@example.com",
  "role": "admin",
  "exp": 1719273600,  // Expires in 24 hours
  "iat": 1719187200,  // Issued at
  "nbf": 1719187200,  // Not before
  "iss": "elektrodukasi"
}
```

### Signature
```
HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload), secret)
```

---

## ⚙️ Configuration

### Token Expiration
Edit `internal/auth/jwt.go`:
```go
const TokenExp = 24 * time.Hour  // Change to desired duration
```

### JWT Secret
**IMPORTANT: Change in production!**
Edit `internal/auth/jwt.go`:
```go
var JWTSecret = []byte("your-secret-key-change-in-production")
```

**Production Setup:**
```go
// Read from environment variable
var JWTSecret = []byte(os.Getenv("JWT_SECRET"))
```

---

## 🧪 Testing

### Using cURL

```bash
# Test login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'

# Test admin endpoint with token
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/users

# Test permission denied
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test"}'
```

### Using Postman

1. **Login Request**
   - Method: POST
   - URL: `http://localhost:8080/api/login`
   - Body (JSON): `{"email":"admin@example.com","password":"admin123"}`
   - Copy `token` from response

2. **Add Token to Headers**
   - Headers tab
   - Add: `Authorization: Bearer <your_token>`

3. **Test Protected Endpoint**
   - Method: POST
   - URL: `http://localhost:8080/api/categories`
   - Headers include Authorization header
   - Body: `{"name":"Test"}`

---

## 📚 Middleware Architecture

### Authentication Flow
```
HTTP Request
    ↓
AuthMiddleware (Validates JWT token)
    ├─→ Valid token → Extract claims, add to context
    └─→ Invalid/Missing → Return 401
        ↓
AuthorizationMiddleware (Checks role & permissions)
    ├─→ Admin role → Allow all operations
    ├─→ User role → Check specific permissions
    └─→ Insufficient → Return 403
        ↓
Handler (Process request)
```

### Middleware Stack
```go
middleware.Chain(
  handler,
  middleware.AuthMiddleware,        // Validate JWT
  middleware.InferAuthorizationFromMethod(resource) // Check permissions
)
```

---

## 🔄 Token Refresh (Future Feature)

Currently, tokens last 24 hours. Consider implementing:
- Refresh token endpoint
- Shorter expiration times (15 minutes) with refresh
- Token revocation list

---

## 📖 Setting Up Users

### Create Admin User Manually

Before the API is fully operational, you need to:

1. **Hash a password** locally or in database:
```go
import "golang.org/x/crypto/bcrypt"

password := "admin123"
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// hashedPassword: "$2a$10$..."
```

2. **Insert into database**:
```sql
INSERT INTO users (name, email, password, role, is_active)
VALUES ('Admin User', 'admin@example.com', '$2a$10$...', 'admin', true);
```

3. **Login**:
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

---

## 🚨 Common Issues

### "Authorization header not provided"
- Check header name: must be `Authorization`
- Check header format: must be `Bearer <token>`

### "Invalid or expired token"
- Token may have expired (24 hour expiration)
- Token may be malformed
- JWT secret mismatch (if multiple servers)

### "Insufficient permissions"
- User role cannot perform this operation
- Only users with 'admin' role can manage categories, tags, users, etc.
- User role can only write (POST/PUT/DELETE) to comments

### "Invalid email or password"
- Email doesn't exist
- Password is incorrect
- Account is deactivated

---

## 🔐 Security Best Practices

### ✅ Do:
- Use HTTPS in production
- Keep JWT secret secure (use environment variables)
- Rotate JWT secret periodically
- Log authentication attempts
- Validate all inputs
- Use short expiration times
- Implement rate limiting on login endpoint

### ❌ Don't:
- Store JWT secret in code
- Use weak passwords
- Accept user input for role assignment
- Log passwords or tokens
- Skip token validation
- Use HTTP in production
- Expose detailed error messages

---

## 📞 API Reference

| Endpoint | Method | Auth | Purpose |
|----------|--------|------|---------|
| `/api/login` | POST | No | Get JWT token |
| `/api/users` | POST | Admin | Create user |
| `/api/users` | GET | Admin | List users |
| `/api/categories` | POST | Admin | Create category |
| `/api/comments` | POST | User | Create comment |
| `/api/articles` | GET | No | List articles |

---

## 🔍 Debugging

### Enable Token Logging
Add to middleware:
```go
log.Printf("Token: %s", token)
log.Printf("Claims: %+v", claims)
```

### Decode JWT Token
Visit: https://jwt.io and paste token to see claims

### Check Token Expiration
```bash
# Extract exp claim
TOKEN="your_token"
# The payload is base64url encoded second part
echo $TOKEN | cut -d'.' -f2 | base64 -d
```

---

## Version Info
- **Auth Type**: JWT (HS256)
- **Token Expiration**: 24 hours
- **Password Hashing**: bcrypt
- **Status**: ✅ Production Ready (after configuration)
