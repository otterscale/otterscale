import { error, redirect } from '@sveltejs/kit';

import type { PageServerLoad } from './$types';

import { resolve } from '$app/paths';
import { env } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';
import { auth } from '$lib/auth';

// Environment variables are loaded from .env in development or system environment in production
export const load: PageServerLoad = async ({ request, url }) => {
	// Validate required environment variables
	const requiredEnvVars = [
		{ value: env.AUTH_SECRET, name: 'AUTH_SECRET' },
		{ value: env.DATABASE_URL, name: 'DATABASE_URL' },
		{ value: publicEnv.PUBLIC_URL, name: 'PUBLIC_URL' },
		{ value: publicEnv.PUBLIC_API_URL, name: 'PUBLIC_API_URL' }
	];

	for (const { value, name } of requiredEnvVars) {
		if (!value) {
			throw error(503, `${name} is not set`);
		}
	}

	if (isFlexibleBooleanTrue(env.BOOTSTRAP_MODE)) {
		throw redirect(302, resolve('/setup'));
	}

	// Check if the user is already authenticated
	const session = await auth.api.getSession({
		headers: request.headers
	});

	if (session) {
		throw redirect(302, resolve('/(auth)/scopes'));
	}

	throw redirect(302, `${resolve('/login')}${url.search}`);
};

const isFlexibleBooleanTrue = (envVar: string | undefined): boolean => {
	return ['true', '1', 'yes', 'on'].includes((envVar || '').toLowerCase());
};
