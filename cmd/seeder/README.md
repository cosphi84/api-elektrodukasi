# Seeder

Database seeder for Elektrodukasi API. This utility inserts initial data into the database.

## Usage

### Prerequisites

1. Database must be running and accessible
2. Database migrations must be run first (`migrations/` folder)
3. Set required environment variables

### Environment Variables

```bash
DATABASE_URL=host=localhost user=postgres password=postgres dbname=elektrodukasi port=5432 sslmode=disable
USER_PASSWORD=your_admin_password
```

### Running the Seeder

```bash
# Using the compiled binary
./seeder

# Or using go run
go run ./cmd/seeder

# With custom environment variables
export DATABASE_URL=host=localhost user=postgres password=postgres dbname=elektrodukasi port=5432 sslmode=disable
export USER_PASSWORD=mySecurePassword123
go run ./cmd/seeder
```

## What Gets Seeded

### Admin User
- **Name**: Risam, S.T
- **Email**: risam1984@gmail.com
- **Password**: From `USER_PASSWORD` environment variable (hashed with bcrypt)
- **Role**: admin
- **Status**: Active

### Category
- **Name**: Komponen Dasar
- **Description**: Dasar Komponen elektronika

## Behavior

- **Duplicate Prevention**: Before inserting, the seeder checks if the user (by email or name) or category (by name) already exists
- **Skip on Exists**: If data already exists, the seeder skips insertion and shows a message
- **Transaction Safe**: Uses GORM's default behavior for atomic operations

## Example Output

```
✓ Connected to database
✓ User created: Risam, S.T (risam1984@gmail.com) with role: admin
✓ Category created: Komponen Dasar

✓ Seeding completed successfully
```

Or if data exists:

```
✓ Connected to database
⊘ User 'risam1984@gmail.com' already exists, skipping
⊘ Category 'Komponen Dasar' already exists, skipping

✓ Seeding completed successfully
```

## Password Security

- Passwords are hashed using bcrypt with the default cost (10 rounds)
- Plain text passwords are never stored in the database
- The password hash is generated from the `USER_PASSWORD` environment variable

## Notes

- The seeder will not overwrite existing data
- If you need to update seed data, manually modify it using your database client
- This seeder is designed to run once to initialize the database with essential data
