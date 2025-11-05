import { error, redirect } from '@sveltejs/kit';

import type { PageServerLoad } from './$types';

import { resolve } from '$app/paths';
import { env } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';

export const load: PageServerLoad = async () => {
	// Environment variables are loaded from .env in development or system environment in production
	const requiredEnvVars = [
		{ value: publicEnv.PUBLIC_URL, name: 'PUBLIC_URL' },
		{ value: publicEnv.PUBLIC_API_URL, name: 'PUBLIC_API_URL' },
		{ value: publicEnv.PUBLIC_AUTH_URL, name: 'PUBLIC_AUTH_URL' },
		{ value: publicEnv.PUBLIC_AUTH_REALM, name: 'PUBLIC_AUTH_REALM' },
		{ value: publicEnv.PUBLIC_AUTH_CLIENT_ID, name: 'PUBLIC_AUTH_CLIENT_ID' },
	];

	for (const { value, name } of requiredEnvVars) {
		if (!value) {
			throw error(503, `${name} is not set`);
		}
	}

	if (isFlexibleBooleanTrue(env.BOOTSTRAP_MODE)) {
		throw redirect(302, resolve('/setup'));
	}
};

const isFlexibleBooleanTrue = (envVar: string | undefined): boolean => {
	return ['true', '1', 'yes', 'on'].includes((envVar || '').toLowerCase());
};
