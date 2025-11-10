import { error, redirect } from '@sveltejs/kit';

import { resolve } from '$app/paths';
import { env } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';

import type { PageServerLoad } from './$types';

const REQUIRED_ENV_VARS = [
	{ value: publicEnv.PUBLIC_WEB_URL, name: 'PUBLIC_WEB_URL' },
	{ value: publicEnv.PUBLIC_API_URL, name: 'PUBLIC_API_URL' },
	{ value: env.KEYCLOAK_REALM_URL, name: 'KEYCLOAK_REALM_URL' },
	{ value: env.KEYCLOAK_CLIENT_ID, name: 'KEYCLOAK_CLIENT_ID' },
	{ value: env.KEYCLOAK_CLIENT_SECRET, name: 'KEYCLOAK_CLIENT_SECRET' }
] as const;

const isFlexibleBooleanTrue = (envVar: string | undefined): boolean => {
	return ['true', '1', 'yes', 'on'].includes((envVar || '').toLowerCase());
};

export const load: PageServerLoad = async ({ locals }) => {
	if (isFlexibleBooleanTrue(env.BOOTSTRAP_MODE)) {
		throw redirect(307, resolve('/setup'));
	}

	for (const { value, name } of REQUIRED_ENV_VARS) {
		if (!value) {
			throw error(503, `${name} is not set`);
		}
	}

	if (!locals.session || !locals.user) {
		throw redirect(307, resolve('/login'));
	}

	throw redirect(307, resolve('/(auth)/scopes'));
};
