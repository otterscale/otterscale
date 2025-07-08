import { env } from "$env/dynamic/private";
import { PUBLIC_BASE_URL } from "$env/static/public";
import { betterAuth } from "better-auth";
import { Pool } from "pg";

if (!env.DATABASE_URL) throw new Error("DATABASE_URL is not set");

export const auth = betterAuth({
	baseURL: PUBLIC_BASE_URL,
	database: new Pool({
		connectionString: env.DATABASE_URL,
	}),
	emailAndPassword: {
		enabled: true,
	},
	session: {
		cookieCache: {
			enabled: true,
			maxAge: 5 * 60,
		},
	},
	socialProviders: {
		apple: {
			clientId: env.APPLE_CLIENT_ID!,
			clientSecret: env.APPLE_CLIENT_SECRET!,
			appBundleIdentifier: env.APPLE_APP_BUNDLE_IDENTIFIER!,
		},
		facebook: {
			clientId: env.FACEBOOK_CLIENT_ID!,
			clientSecret: env.FACEBOOK_CLIENT_SECRET!,
		},
		github: {
			clientId: env.GITHUB_CLIENT_ID!,
			clientSecret: env.GITHUB_CLIENT_SECRET!,
		},
		google: {
			clientId: env.GOOGLE_CLIENT_ID!,
			clientSecret: env.GOOGLE_CLIENT_SECRET!,
		},
		twitter: {
			clientId: env.TWITTER_CLIENT_ID!,
			clientSecret: env.TWITTER_CLIENT_SECRET!,
		},
	},
});
