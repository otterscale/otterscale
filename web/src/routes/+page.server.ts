import { error, redirect } from '@sveltejs/kit';

import { resolve } from '$app/paths';
import {
	BOOTSTRAP_MODE,
	KEYCLOAK_CLIENT_ID,
	KEYCLOAK_CLIENT_SECRET,
	KEYCLOAK_REALM_URL
} from '$env/static/private';
import { PUBLIC_API_URL, PUBLIC_WEB_URL } from '$env/static/public';

import type { PageServerLoad } from './$types';

const REQUIRED_ENV_VARS = [
	{ value: PUBLIC_WEB_URL, name: 'PUBLIC_WEB_URL' },
	{ value: PUBLIC_API_URL, name: 'PUBLIC_API_URL' },
	{ value: KEYCLOAK_REALM_URL, name: 'KEYCLOAK_REALM_URL' },
	{ value: KEYCLOAK_CLIENT_ID, name: 'KEYCLOAK_CLIENT_ID' },
	{ value: KEYCLOAK_CLIENT_SECRET, name: 'KEYCLOAK_CLIENT_SECRET' }
] as const;

const isFlexibleBooleanTrue = (envVar: string | undefined): boolean => {
	return ['true', '1', 'yes', 'on'].includes((envVar || '').toLowerCase());
};

export const load: PageServerLoad = async ({ locals }) => {
	if (isFlexibleBooleanTrue(BOOTSTRAP_MODE)) {
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
