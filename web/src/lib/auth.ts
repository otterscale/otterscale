import { betterAuth } from 'better-auth';
import { getMigrations } from 'better-auth/db';
import { Pool } from 'pg';
import { sso } from '@better-auth/sso';
import { env } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';

export const auth = betterAuth({
	account: {
		accountLinking: {
			enabled: true,
			trustedProviders: env.AUTH_TRUSTED_PROVIDERS?.split(',') || [],
		},
	},
	baseURL: env.PUBLIC_URL,
	database: new Pool({
		connectionString: env.DATABASE_URL,
	}),
	emailAndPassword: {
		enabled: true,
	},
	plugins: [sso()],
	secret: env.AUTH_SECRET,
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
		github: {
			clientId: env.GITHUB_CLIENT_ID!,
			clientSecret: env.GITHUB_CLIENT_SECRET!,
		},
		google: {
			clientId: env.GOOGLE_CLIENT_ID!,
			clientSecret: env.GOOGLE_CLIENT_SECRET!,
		},
	},
	telemetry: { enabled: false },
	trustedOrigins: [publicEnv.PUBLIC_URL].filter(Boolean),
});

async function initializeDatabase() {
	const { runMigrations } = await getMigrations(auth.options);
	await runMigrations();
	console.log('Database migrations completed successfully');
}

if (!process.env.BUILD) {
	initializeDatabase().catch((error) => {
		console.error('Failed to initialize database:', error);
		process.exit(1);
	});
}
