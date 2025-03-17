import { betterAuth, type User } from "better-auth";
import { admin, jwt, openAPI, organization } from "better-auth/plugins"
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { drizzle } from "drizzle-orm/node-postgres";
import * as schema from "./auth-schema"
import pg from "pg";

const { Pool } = pg;

const pool = new Pool({
    connectionString: import.meta.env.OPENHDC_CONNECTION_STRING,
})

const db = drizzle({ client: pool });

export const auth = betterAuth({
    database: drizzleAdapter(db, {
        provider: "pg", schema: schema
    }),
    emailAndPassword: {
        enabled: true
    },
    socialProviders: {
        github: {
            clientId: import.meta.env.GITHUB_CLIENT_ID ?? '',
            clientSecret: import.meta.env.GITHUB_CLIENT_SECRET ?? '',
        }
    },
    plugins: [
        admin(),
        organization({
            teams: {
                enabled: true,
            }
        }),
        openAPI(),
        jwt(),
    ]
})