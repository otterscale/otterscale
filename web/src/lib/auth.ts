import * as dotenv from 'dotenv';
import { betterAuth } from "better-auth";
import { admin, jwt, openAPI, organization } from "better-auth/plugins"
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { drizzle } from "drizzle-orm/node-postgres";
import * as schema from "./auth-schema"
import pg from "pg";

dotenv.config();

const { Pool } = pg;

const pool = new Pool({
    connectionString: process.env.BETTER_AUTH_CONNECTION_STRING,
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
        apple: {
            clientId: process.env.APPLE_CLIENT_ID as string,
            clientSecret: process.env.APPLE_CLIENT_SECRET as string,
            appBundleIdentifier: process.env.APPLE_APP_BUNDLE_IDENTIFIER as string,
        },
        discord: {
            clientId: process.env.DISCORD_CLIENT_ID as string,
            clientSecret: process.env.DISCORD_CLIENT_SECRET as string,
        },
        facebook: {
            clientId: process.env.FACEBOOK_CLIENT_ID as string,
            clientSecret: process.env.FACEBOOK_CLIENT_SECRET as string,
        },
        github: {
            clientId: process.env.GITHUB_CLIENT_ID as string,
            clientSecret: process.env.GITHUB_CLIENT_SECRET as string,
        },
        google: {
            clientId: process.env.GOOGLE_CLIENT_ID as string,
            clientSecret: process.env.GOOGLE_CLIENT_SECRET as string,
        },
        microsoft: {
            clientId: process.env.MICROSOFT_CLIENT_ID as string,
            clientSecret: process.env.MICROSOFT_CLIENT_SECRET as string,
            tenantId: 'common',
            requireSelectAccount: true
        },
        tiktok: {
            clientId: process.env.TIKTOK_CLIENT_ID as string,
            clientSecret: process.env.TIKTOK_CLIENT_SECRET as string,
            clientKey: process.env.TIKTOK_CLIENT_KEY as string,
        },
        twitter: {
            clientId: process.env.TWITTER_CLIENT_ID as string,
            clientSecret: process.env.TWITTER_CLIENT_SECRET as string,
        },
        linkedin: {
            clientId: process.env.LINKEDIN_CLIENT_ID as string,
            clientSecret: process.env.LINKEDIN_CLIENT_SECRET as string,
        },
        gitlab: {
            clientId: process.env.GITLAB_CLIENT_ID as string,
            clientSecret: process.env.GITLAB_CLIENT_SECRET as string,
            issuer: process.env.GITLAB_ISSUER as string,
        },
        reddit: {
            clientId: process.env.REDDIT_CLIENT_ID as string,
            clientSecret: process.env.REDDIT_CLIENT_SECRET as string,
        },
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