# Postman Collection Guide

This guide explains how to import and use the Elektrodukasi API Postman collection.

## Importing the Collection

### Method 1: Import from File
1. Open **Postman**
2. Click **File** → **Import**
3. Select the **Elektrodukasi-API.postman_collection.json** file
4. Click **Import**

### Method 2: Import via URL
If the file is hosted online, paste the URL in the import dialog.

## Setting Up Variables

The collection includes the following variables that you should configure:

| Variable | Default | Description |
|----------|---------|-------------|
| `base_url` | `http://localhost:8080` | API base URL |
| `token` | (empty) | JWT token from login (auto-fill after login) |
| `user_id` | (empty) | User ID for user-specific requests |
| `category_id` | (empty) | Category ID |
| `tag_id` | (empty) | Tag ID |
| `article_id` | (empty) | Article ID |
| `article_slug` | (empty) | Article slug for slug-based queries |
| `comment_id` | (empty) | Comment ID |
| `project_id` | (empty) | Project ID |
| `project_slug` | (empty) | Project slug for slug-based queries |
| `owner_id` | (empty) | User ID for owner-based queries |

### Updating Variables

1. Click the **Variables** tab at the top
2. Edit the values in the **INITIAL VALUE** or **CURRENT VALUE** columns
3. The most important variables to set initially:
   - `base_url`: Change if your API is not running on `localhost:8080`
   - Default port in code is `8080`, change if needed

## Quick Start Workflow

### Step 1: Login
1. Go to **Authentication** → **Login**
2. Make sure the email and password match your seeded admin user:
   - Email: `risam1984@gmail.com`
   - Password: Your `USER_PASSWORD` environment variable
3. Click **Send**
4. Copy the `token` from the response

### Step 2: Save Token to Variable
1. The response contains `token` field
2. Copy the token value
3. Go to **Variables** tab and paste into `token` variable
4. Or use Postman's test scripts to auto-fill (see section below)

### Step 3: Use Protected Endpoints
All protected endpoints use `Authorization: Bearer {{token}}` header.
Once token is set, you can call admin/user endpoints.

## Endpoint Organization

The collection is organized into folders:

### Authentication
- **Login** - POST `/api/login` (Public)

### Users (Admin only)
- List Users - GET `/api/users`
- Get User - GET `/api/users/{id}`
- Create User - POST `/api/users`
- Update User - PUT `/api/users/{id}`
- Deactivate User - PATCH `/api/users/{id}/deactivate`
- Delete User - DELETE `/api/users/{id}`

### Categories (Read: Public, Write: Admin only)
- List Categories - GET `/api/categories` (Public)
- Get Category - GET `/api/categories/{id}` (Public)
- Create Category - POST `/api/categories` (Admin)
- Update Category - PUT `/api/categories/{id}` (Admin)
- Delete Category - DELETE `/api/categories/{id}` (Admin)

### Tags (Read: Public, Write: Admin only)
- List Tags - GET `/api/tags` (Public)
- Get Tag - GET `/api/tags/{id}` (Public)
- Create Tag - POST `/api/tags` (Admin)
- Update Tag - PUT `/api/tags/{id}` (Admin)
- Delete Tag - DELETE `/api/tags/{id}` (Admin)

### Articles (Read: Public, Write: Admin only)
- List Articles - GET `/api/articles` (Public)
- Get Article - GET `/api/articles/{id}` (Public)
- Get Article by Slug - GET `/api/articles/slug/{slug}` (Public)
- Get Articles by Category - GET `/api/categories/{categoryId}/articles` (Public)
- Create Article - POST `/api/articles` (Admin)
- Update Article - PUT `/api/articles/{id}` (Admin)
- Publish Article - PATCH `/api/articles/{id}/publish` (Admin)
- Add Tag - POST `/api/articles/{id}/tags` (Admin)
- Remove Tag - DELETE `/api/articles/{id}/tags/{tagId}` (Admin)
- Delete Article - DELETE `/api/articles/{id}` (Admin)

### Comments (Read: Public, Write: User+Admin)
- Get Comments by Article - GET `/api/articles/{articleId}/comments` (Public)
- Get Comment - GET `/api/comments/{id}` (Public)
- Get Comment Replies - GET `/api/comments/{id}/replies` (Public)
- Create Comment - POST `/api/comments` (User/Admin)
- Update Comment - PUT `/api/comments/{id}` (User/Admin)
- Approve Comment - PATCH `/api/comments/{id}/approve` (Admin)
- Delete Comment - DELETE `/api/comments/{id}` (User/Admin)

### Projects (Read: Public, Write: Admin only)
- List Projects - GET `/api/projects` (Public)
- Get Project - GET `/api/projects/{id}` (Public)
- Get Project by Slug - GET `/api/projects/slug/{slug}` (Public)
- Get User Projects - GET `/api/users/{ownerId}/projects` (Authenticated)
- Create Project - POST `/api/projects` (Admin)
- Update Project - PUT `/api/projects/{id}` (Admin)
- Delete Project - DELETE `/api/projects/{id}` (Admin)

## Common Tasks

### Task 1: Create a Category and Add Articles

1. **Create Category** (Admin only)
   - Go to Categories → Create Category
   - Edit body with your category details
   - Send request
   - Copy `id` from response

2. **Set category_id variable**
   - Go to Variables tab
   - Paste the ID

3. **Create Article** (Admin only)
   - Go to Articles → Create Article
   - Use `{{category_id}}` in the body (or copy-paste)
   - Send request

### Task 2: Add Comments to Article

1. **Get article_id** - List articles or get specific article

2. **Create Comment** (User/Admin)
   - Go to Comments → Create Comment
   - Use `{{article_id}}` in body
   - Send request (token required)

3. **Approve Comment** (Admin only)
   - Go to Comments → Approve Comment
   - Set `{{comment_id}}` from previous response
   - Send request

### Task 3: Create and Tag Articles

1. **Create Tag** (Admin)
   - Go to Tags → Create Tag
   - Send request
   - Note `tag_id`

2. **Add Tag to Article** (Admin)
   - Go to Articles → Add Tag to Article
   - Use `{{article_id}}` and `{{tag_id}}`
   - Send request

## Authentication Notes

- **Token Expiration**: Tokens expire in 24 hours
- **Token Format**: `Authorization: Bearer <token>`
- **Re-login Required**: If you get 401 error, login again
- **Role-based Access**:
  - **Admin**: Can do anything
  - **User**: Can only write comments

## Error Codes

Common responses and meanings:

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (check your JSON body) |
| 401 | Unauthorized (token missing or invalid) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found (resource doesn't exist) |
| 500 | Server Error (contact support) |

## Tips & Tricks

### Auto-fill Token After Login
You can use Postman's test scripts to auto-fill the token:

1. Go to **Authentication** → **Login** request
2. Click **Tests** tab
3. Add:
```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.variables.set("token", jsonData.token);
}
```
4. Now token auto-fills after each login

### Create Environment File
For different environments (dev, staging, prod):

1. Click **Environments** (top left)
2. Create new environment (e.g., "Development", "Production")
3. Set `base_url` for each environment
4. Switch environments easily

### Use Pre-request Scripts
Set variables before making requests:
- Click **Pre-request Script** tab
- Add dynamic values (timestamps, UUIDs, etc.)

## Troubleshooting

**"401 Unauthorized"**
- Check your token is set in variables
- Token may have expired (24 hours) - login again

**"403 Forbidden"**
- You don't have permission for this action
- Check your role (admin vs user)
- User role can only write comments

**"404 Not Found"**
- The resource doesn't exist
- Check the ID is correct
- Try listing resources first to find valid IDs

**"400 Bad Request"**
- Check your JSON body syntax
- Verify required fields are present
- Check data types match expectations

**Connection Refused**
- API server not running
- Check `base_url` variable is correct
- Verify database is running
- Check API is listening on the correct port

## Further Reading

- See `AUTH_DOCUMENTATION.md` for authentication details
- See `API_DOCUMENTATION.md` for detailed endpoint documentation
- See `README_CRUD.md` for CRUD operation details
