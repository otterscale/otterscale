import { env } from "$env/dynamic/private";
import { env as publicEnv } from "$env/dynamic/public";
import { betterAuth } from "better-auth";
import { Pool } from "pg";

export const auth = betterAuth({
	baseURL: publicEnv.PUBLIC_URL,
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
	secret: env.AUTH_SECRET,
	socialProviders: {
		apple: {
			clientId: env.APPLE_CLIENT_ID!,
			clientSecret: env.APPLE_CLIENT_SECRET!,
			appBundleIdentifier: env.APPLE_APP_BUNDLE_IDENTIFIER!,
		},
		github: {
			clientId: env.GITHUB_CLIENT_ID!,
			clientSecret: env.GITHUB_CLIENT_SECRET!,
		},
		google: {
			clientId: env.GOOGLE_CLIENT_ID!,
			clientSecret: env.GOOGLE_CLIENT_SECRET!,
		},
	},
	trustedOrigins: [publicEnv.PUBLIC_URL],
});
