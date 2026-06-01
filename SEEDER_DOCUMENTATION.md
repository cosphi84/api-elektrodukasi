# Seeder Documentation

## Overview

A database seeder has been created to initialize essential data for the Elektrodukasi API. The seeder runs from `cmd/seeder/main.go` and automatically checks for existing data before inserting.

## Features

✓ **Duplicate Prevention**: Checks if user (by email or name) exists before inserting
✓ **Duplicate Prevention**: Checks if category (by name) exists before inserting
✓ **Password Hashing**: Uses bcrypt for secure password storage
✓ **Environment Variables**: Reads `USER_PASSWORD` from environment
✓ **Error Handling**: Comprehensive error messages for debugging
✓ **Database Connection**: Uses same database setup as main API

## Seeded Data

### Admin User
```
Name: Risam, S.T
Email: risam1984@gmail.com
Password: Hashed from USER_PASSWORD environment variable
Role: admin
Status: Active
```

### Category
```
Name: Komponen Dasar
Description: Dasar Komponen elektronika
```

## How to Use

### Step 1: Set Environment Variables

```bash
export DATABASE_URL=host=localhost user=postgres password=postgres dbname=elektrodukasi port=5432 sslmode=disable
export USER_PASSWORD=your_secure_password
```

### Step 2: Run Database Migrations

Ensure database schema is created before running seeder:

```bash
# Run migrations if not already done
# See migrations/ folder for SQL files
```

### Step 3: Run the Seeder

```bash
# Using go run
go run ./cmd/seeder

# Or after building
go build ./cmd/seeder
./seeder
```

## Expected Output

**First run (fresh database):**
```
✓ Connected to database
✓ User created: Risam, S.T (risam1984@gmail.com) with role: admin
✓ Category created: Komponen Dasar

✓ Seeding completed successfully
```

**Subsequent runs (data already exists):**
```
✓ Connected to database
⊘ User 'risam1984@gmail.com' already exists, skipping
⊘ Category 'Komponen Dasar' already exists, skipping

✓ Seeding completed successfully
```

## Implementation Details

### Files Created
- `cmd/seeder/main.go` - Main seeder application
- `cmd/seeder/README.md` - Detailed seeder documentation

### Key Functions

#### `seedUser(db *gorm.DB, password string)`
- Checks if user exists by email or name
- Skips insertion if found
- Hashes password using bcrypt
- Creates user with admin role

#### `seedCategory(db *gorm.DB)`
- Checks if category exists by name
- Skips insertion if found
- Creates category with description

### Dependencies Used
- `gorm.io/gorm` - ORM operations
- `gorm.io/driver/postgres` - PostgreSQL driver
- `golang.org/x/crypto/bcrypt` - Password hashing

## Security Considerations

1. **Password Hashing**: All passwords are hashed using bcrypt with default cost (10 rounds)
2. **Environment Variables**: Sensitive data (passwords, database URL) passed via environment, not hardcoded
3. **Duplicate Prevention**: Prevents accidental data duplication
4. **No Plain Text Passwords**: Password from environment is only used to generate hash

## Future Enhancements

Optional improvements:
- Add command-line flags for custom seed data
- Create additional seeders for tags, articles, etc.
- Add rollback functionality
- Add seed verification/validation
- Add batch seeding from CSV/JSON files

## Troubleshooting

**Error: "Failed to connect to database"**
- Verify PostgreSQL is running
- Check DATABASE_URL is correct
- Verify database exists and migrations are applied

**Error: "USER_PASSWORD environment variable is not set"**
- Set the USER_PASSWORD environment variable before running seeder
- Example: `export USER_PASSWORD=myPassword123`

**Error: "Failed to create user"**
- Check database migrations have been run
- Verify user table exists with correct schema
- Check for database constraints (email uniqueness)

## Related Documentation

- See `AUTH_DOCUMENTATION.md` for login and authentication details
- See `QUICK_REFERENCE.md` for API endpoint quick reference
- See `README_CRUD.md` for CRUD operations documentation
