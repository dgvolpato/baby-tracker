# Baby Tracker

A private web app for two parents to log and review a baby's feedings.

## Stack

| Layer | Tech |
|---|---|
| Frontend | SvelteKit 2 + vanilla CSS (S3 + CloudFront) |
| Backend | Go on AWS Lambda (HTTP API) |
| Database | DynamoDB single-table |
| Auth | Hardcoded credentials + JWT (env vars) |
| IaC | AWS SAM |

## Project structure

```
backend/
├── cmd/api/main.go               # Lambda entry point + routing
├── internal/
│   ├── auth/jwt.go               # JWT sign/verify + credential check
│   ├── handlers/                 # login, feedings CRUD
│   └── store/dynamodb.go         # DynamoDB operations
├── template.yaml                 # SAM template
└── Makefile

frontend/
├── src/
│   ├── lib/
│   │   ├── api.ts                # fetch wrapper with JWT
│   │   └── utils.ts              # time helpers, formatters
│   └── routes/
│       ├── login/                # sign-in page
│       ├── dashboard/            # daily totals
│       └── history/              # feeding list + add/edit modal
└── package.json
```

## Local development

### Backend

Requires [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html) and Go 1.22+.

```bash
cd backend
make build
make local        # starts API on http://localhost:3001
```

### Frontend

```bash
cd frontend
cp .env.example .env.local
# edit .env.local → set VITE_API_URL=http://localhost:3001
npm install
npm run dev       # http://localhost:5173
```

## AWS deployment

### 1. Store secrets in SSM Parameter Store

```bash
aws ssm put-parameter --name /baby-tracker/jwt-secret     --value "<random-secret>" --type SecureString
aws ssm put-parameter --name /baby-tracker/user1-name     --value "dad"             --type SecureString
aws ssm put-parameter --name /baby-tracker/user1-password --value "<password>"      --type SecureString
aws ssm put-parameter --name /baby-tracker/user2-name     --value "mom"             --type SecureString
aws ssm put-parameter --name /baby-tracker/user2-password --value "<password>"      --type SecureString
```

### 2. Deploy backend

```bash
cd backend
make deploy       # runs sam deploy --guided on first run
```

Copy the `ApiUrl` from the outputs.

### 3. Deploy frontend

```bash
cd frontend
echo "VITE_API_URL=https://<your-api-url>" > .env.production
npm run build
aws s3 sync build/ s3://<your-bucket> --delete
```

## API

All endpoints except `POST /login` require `Authorization: Bearer <token>`.

| Method | Path | Description |
|---|---|---|
| POST | `/login` | Returns JWT |
| GET | `/feedings` | List all feedings (sorted newest first) |
| POST | `/feedings` | Create a feeding |
| PUT | `/feedings/:id` | Update a feeding |
| DELETE | `/feedings/:id` | Delete a feeding |

### Feeding object

```json
{
  "id": "uuid",
  "timestamp": "2026-05-23T08:30:00Z",
  "type": "formula",
  "oz": 4.5,
  "createdBy": "dad"
}
```
