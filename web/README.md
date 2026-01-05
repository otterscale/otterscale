# OtterScale Web

## Environment Variables

The following environment variables are required for the application to run properly:

| Variable                 | Description                                    | Required | Example Value                                      |
| ------------------------ | ---------------------------------------------- | -------- | -------------------------------------------------- |
| `PUBLIC_WEB_URL`         | The public URL where the application is hosted | Yes      | `http://localhost:3000`                            |
| `API_URL`                | The API endpoint URL                           | Yes      | `http://localhost:8299`                            |
| `REDIS_URL`              | Redis connection string                        | Yes      | `redis://your_secure_password_here@localhost:6379` |
| `KEYCLOAK_REALM_URL`     | Keycloak realm URL for authentication          | Yes      | `http://localhost:8080/realms/your_realm`          |
| `KEYCLOAK_CLIENT_ID`     | Keycloak client ID                             | Yes      | `your_client_id`                                   |
| `KEYCLOAK_CLIENT_SECRET` | Keycloak client secret                         | Yes      | `your_client_secret`                               |
| `BOOTSTRAP_MODE`         | Bootstrap mode flag                            | false    | `0`                                                |
