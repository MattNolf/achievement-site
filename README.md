# Site

This is the landing page for [achievement.dev](https://achievement.dev). It currently doesn't do anything other than store interest to a table in [Airtable](https://airtable.com).
It is intentionally minimal.

### Running Locally

Download dependancies with go modules:
```
go mod tidy
```

Run the project with required environment variables:
```
AIRTABLE_API_KEY=SECET AIRTABLE_DATABASE_ID=did AIRTABLE_TABLE_ID=tid go run main.go
```

### Build


With Docker:
```
docker build .
```

Or for a specific architecture and push:
```
docker buildx build --platform linux/amd64 --push -t {{registry}} .
```