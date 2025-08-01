import { redirect } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";
import { auth } from "$lib/auth";
import { applicationsPath, loginPath } from "$lib/path";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ request, url }) => {
	// Environment variables are loaded from .env in development or system environment in production
	if (!env.AUTH_SECRET) {
		return {
			error: "AUTH_SECRET is not set"
		}
	};

	if (!env.DATABASE_URL) {
		return {
			error: "DATABASE_URL is not set"
		}
	};

	// Check if the user is already authenticated
	const session = await auth.api.getSession({
		headers: request.headers,
	});

	if (session) {
		redirect(302, applicationsPath);
	}

	redirect(302, `${loginPath}${url.search}`);
};
