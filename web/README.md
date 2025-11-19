# OtterScale Web

## Environment Variables

The following environment variables are required for the application to run properly:

| Variable                 | Description                                    | Required | Example Value                                                                 |
| ------------------------ | ---------------------------------------------- | -------- | ----------------------------------------------------------------------------- |
| `PUBLIC_WEB_URL`         | The public URL where the application is hosted | Yes      | `http://localhost:3000`                                                       |
| `PUBLIC_API_URL`         | The public API endpoint URL                    | Yes      | `http://localhost:8299`                                                       |
| `KEYCLOAK_REALM_URL`     | Keycloak realm URL for authentication          | Yes      | `http://localhost:8080/realms/your_realm`                                     |
| `KEYCLOAK_CLIENT_ID`     | Keycloak client ID                             | Yes      | `your_client_id`                                                              |
| `KEYCLOAK_CLIENT_SECRET` | Keycloak client secret                         | Yes      | `your_client_secret`                                                          |
| `DATABASE_URL`           | PostgreSQL database connection string          | Yes      | `postgresql://otterscale:your_secure_password_here@localhost:5432/otterscale` |
| `BOOTSTRAP_MODE`         | Bootstrap mode flag                            | false    | `0`                                                                           |
