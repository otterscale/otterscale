# Otterscale Web

## Environment Variables

The following environment variables are required for the application to run properly:

| Variable | Description | Required | Example Value |
|----------|-------------|----------|---------------|
| `PUBLIC_URL` | The public URL where the application is hosted | Yes | `https://otterscale.example.com` |
| `PUBLIC_API_URL` | The public API endpoint URL | Yes | `https://otterscale-service.example.com` |
| `AUTH_SECRET` | Secret key used for authentication token signing | Yes | `Qf6XiQEthdq2d8uQJqLvcQtg7QEz3JUe` |
| `AUTH_TRUSTED_PROVIDERS` | Comma-separated list of trusted authentication providers | No | `otterscale,otterscale-oidc` |
| `AUTH_OIDC_PROVIDER` | The OIDC provider identifier | No | `otterscale-oidc` |
| `DATABASE_URL` | PostgreSQL database connection string | Yes | `postgresql://USER:PASSWORD@localhost:5432/postgres` |
| `GITHUB_CLIENT_ID` | GitHub OAuth application client ID | No | `-` |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth application client secret | No | `-` |
| `SSO_LOGIN_PROMPT` | Enable/disable SSO login prompt | No | `true` |

## Developing

Before starting development, you need to set up the database schema:

```bash
pnpm migrate
```

This command will create the necessary database tables and schema required for the application.

Once you've installed dependencies and set up the database, start the development server:

```bash
pnpm dev

# or start the server and open the app in a new browser tab
pnpm dev -- --open
```

### Setup Instructions

1. Copy the environment variables above to a `.env` file in your project
2. Replace the example values with your actual configuration values
3. Ensure the database is running and accessible at the specified `DATABASE_URL`
4. (Optional) Configure your GitHub OAuth application with the appropriate client ID and secret
5. (Optional) Set up your OIDC provider according to your authentication requirements
