import { env } from '$env/dynamic/private';
import { betterAuth } from "better-auth";
import { getMigrations } from "better-auth/db";
import { admin, jwt, openAPI, organization } from "better-auth/plugins"
import Database from "better-sqlite3";

export const auth = betterAuth({
    database: new Database("./db.sqlite"),
    emailAndPassword: {
        enabled: true
    },
    socialProviders: {
        apple: {
            clientId: env.APPLE_CLIENT_ID,
            clientSecret: env.APPLE_CLIENT_SECRET,
            appBundleIdentifier: env.APPLE_APP_BUNDLE_IDENTIFIER,
        },
        facebook: {
            clientId: env.FACEBOOK_CLIENT_ID,
            clientSecret: env.FACEBOOK_CLIENT_SECRET,
        },
        github: {
            clientId: env.GITHUB_CLIENT_ID,
            clientSecret: env.GITHUB_CLIENT_SECRET,
        },
        google: {
            clientId: env.GOOGLE_CLIENT_ID,
            clientSecret: env.GOOGLE_CLIENT_SECRET,
        },
        microsoft: {
            clientId: env.MICROSOFT_CLIENT_ID,
            clientSecret: env.MICROSOFT_CLIENT_SECRET,
            tenantId: 'common',
            requireSelectAccount: true
        },
        tiktok: {
            clientId: env.TIKTOK_CLIENT_ID,
            clientSecret: env.TIKTOK_CLIENT_SECRET,
            clientKey: env.TIKTOK_CLIENT_KEY,
        },
        twitter: {
            clientId: env.TWITTER_CLIENT_ID,
            clientSecret: env.TWITTER_CLIENT_SECRET,
        },
        linkedin: {
            clientId: env.LINKEDIN_CLIENT_ID,
            clientSecret: env.LINKEDIN_CLIENT_SECRET,
        },
        gitlab: {
            clientId: env.GITLAB_CLIENT_ID,
            clientSecret: env.GITLAB_CLIENT_SECRET,
            issuer: env.GITLAB_ISSUER,
        },
        reddit: {
            clientId: env.REDDIT_CLIENT_ID,
            clientSecret: env.REDDIT_CLIENT_SECRET,
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

async function runMigration() {
    try {
        const { runMigrations } = await getMigrations(auth.options);
        await runMigrations();
        console.log("Migration completed successfully.");
    } catch (error) {
        console.error("Migration failed:", error);
        process.exit(1);
    }
}

runMigration()