import { error, redirect } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";
import { env as publicEnv } from "$env/dynamic/public";
import { auth } from "$lib/auth";
import { staticPaths } from "$lib/path";
import type { PageServerLoad } from "./$types";

// Environment variables are loaded from .env in development or system environment in production
export const load: PageServerLoad = async ({ request, url }) => {
	// Validate required environment variables
	const requiredEnvVars = [
		{ value: env.AUTH_SECRET, name: "AUTH_SECRET" },
		{ value: env.DATABASE_URL, name: "DATABASE_URL" },
		{ value: publicEnv.PUBLIC_URL, name: "PUBLIC_URL" },
		{ value: publicEnv.PUBLIC_API_URL, name: "PUBLIC_API_URL" }
	];

	for (const { value, name } of requiredEnvVars) {
		if (!value) {
			error(503, `${name} is not set`)
		}
	}

	// Check if the user is already authenticated
	const session = await auth.api.getSession({
		headers: request.headers,
	});

	if (session) {
		redirect(302, staticPaths.scopes.url);
	}

	redirect(302, `${staticPaths.login.url}${url.search}`);
};
