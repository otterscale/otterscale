import { betterAuth } from "better-auth";
import { admin, jwt, openAPI, organization } from "better-auth/plugins"
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { drizzle } from "drizzle-orm/node-postgres";
import * as schema from "./auth-schema"
import pg from "pg";

import {
    BETTER_AUTH_CONNECTION_STRING,
    APPLE_CLIENT_ID,
    APPLE_CLIENT_SECRET,
    APPLE_APP_BUNDLE_IDENTIFIER,
    FACEBOOK_CLIENT_ID,
    FACEBOOK_CLIENT_SECRET,
    GITHUB_CLIENT_ID,
    GITHUB_CLIENT_SECRET,
    GOOGLE_CLIENT_ID,
    GOOGLE_CLIENT_SECRET,
    MICROSOFT_CLIENT_ID,
    MICROSOFT_CLIENT_SECRET,
    TIKTOK_CLIENT_ID,
    TIKTOK_CLIENT_SECRET,
    TIKTOK_CLIENT_KEY,
    TWITTER_CLIENT_ID,
    TWITTER_CLIENT_SECRET,
    LINKEDIN_CLIENT_ID,
    LINKEDIN_CLIENT_SECRET,
    GITLAB_CLIENT_ID,
    GITLAB_CLIENT_SECRET,
    GITLAB_ISSUER,
    REDDIT_CLIENT_ID,
    REDDIT_CLIENT_SECRET,
} from '$env/static/private';

const { Pool } = pg;

const pool = new Pool({
    connectionString: BETTER_AUTH_CONNECTION_STRING,
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
            clientId: APPLE_CLIENT_ID,
            clientSecret: APPLE_CLIENT_SECRET,
            appBundleIdentifier: APPLE_APP_BUNDLE_IDENTIFIER,
        },
        facebook: {
            clientId: FACEBOOK_CLIENT_ID,
            clientSecret: FACEBOOK_CLIENT_SECRET,
        },
        github: {
            clientId: GITHUB_CLIENT_ID,
            clientSecret: GITHUB_CLIENT_SECRET,
        },
        google: {
            clientId: GOOGLE_CLIENT_ID,
            clientSecret: GOOGLE_CLIENT_SECRET,
        },
        microsoft: {
            clientId: MICROSOFT_CLIENT_ID,
            clientSecret: MICROSOFT_CLIENT_SECRET,
            tenantId: 'common',
            requireSelectAccount: true
        },
        tiktok: {
            clientId: TIKTOK_CLIENT_ID,
            clientSecret: TIKTOK_CLIENT_SECRET,
            clientKey: TIKTOK_CLIENT_KEY,
        },
        twitter: {
            clientId: TWITTER_CLIENT_ID,
            clientSecret: TWITTER_CLIENT_SECRET,
        },
        linkedin: {
            clientId: LINKEDIN_CLIENT_ID,
            clientSecret: LINKEDIN_CLIENT_SECRET,
        },
        gitlab: {
            clientId: GITLAB_CLIENT_ID,
            clientSecret: GITLAB_CLIENT_SECRET,
            issuer: GITLAB_ISSUER,
        },
        reddit: {
            clientId: REDDIT_CLIENT_ID,
            clientSecret: REDDIT_CLIENT_SECRET,
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