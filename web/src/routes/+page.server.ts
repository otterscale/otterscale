import { error, redirect } from '@sveltejs/kit';

import { resolve } from '$app/paths';
import { env } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';
import { isFlexibleBooleanTrue } from '$lib/helper';

import type { PageServerLoad } from './$types';

interface EnvVar {
	value: string | undefined;
	name: string;
}

const REQUIRED_ENV_VARS_ALL: readonly EnvVar[] = [
	{ value: publicEnv.PUBLIC_API_URL, name: 'PUBLIC_API_URL' }
];

const REQUIRED_ENV_VARS_NORMAL: readonly EnvVar[] = [
	{ value: publicEnv.PUBLIC_WEB_URL, name: 'PUBLIC_WEB_URL' },
	{ value: env.KEYCLOAK_REALM_URL, name: 'KEYCLOAK_REALM_URL' },
	{ value: env.KEYCLOAK_CLIENT_ID, name: 'KEYCLOAK_CLIENT_ID' },
	{ value: env.KEYCLOAK_CLIENT_SECRET, name: 'KEYCLOAK_CLIENT_SECRET' },
	{ value: env.DATABASE_URL, name: 'DATABASE_URL' }
];

const checkRequiredEnvVars = (envVars: readonly EnvVar[]): void => {
	for (const { value, name } of envVars) {
		if (!value) {
			console.error(`Missing required environment variable: ${name}`);
			throw error(503, `${name} is not set`);
		}
	}
};

export const load: PageServerLoad = async ({ locals }) => {
	checkRequiredEnvVars(REQUIRED_ENV_VARS_ALL);

	if (isFlexibleBooleanTrue(env.BOOTSTRAP_MODE)) {
		throw redirect(307, resolve('/setup'));
	}

	checkRequiredEnvVars(REQUIRED_ENV_VARS_NORMAL);

	if (!locals.session || !locals.user) {
		throw redirect(307, resolve('/login'));
	}

	throw redirect(307, resolve('/(auth)/scopes'));
};
