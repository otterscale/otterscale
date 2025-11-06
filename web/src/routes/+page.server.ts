import { error, redirect, fail } from '@sveltejs/kit';

import type { PageServerLoad, Actions } from './$types';

import { resolve } from '$app/paths';
import { env } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';

const REQUIRED_PUBLIC_ENV_VARS = [
	'PUBLIC_URL',
	'PUBLIC_API_URL',
	'PUBLIC_AUTH_URL',
	'PUBLIC_AUTH_REALM',
	'PUBLIC_AUTH_CLIENT_ID',
] as const;

const isFlexibleBooleanTrue = (value: string | undefined): boolean => {
	return ['true', '1', 'yes', 'on'].includes((value || '').toLowerCase());
};

const validateEnvironmentVariables = () => {
	for (const varName of REQUIRED_PUBLIC_ENV_VARS) {
		if (!publicEnv[varName]) {
			throw error(503, `${varName} is not set`);
		}
	}
};

export const actions: Actions = {
	setToken: async ({ request, cookies }) => {
		const data = await request.formData();
		const token = data.get('token');

		if (!token || typeof token !== 'string') {
			return fail(400, { message: 'Token is required' });
		}

		cookies.set('OS_TOKEN', token, {
			path: '/',
			httpOnly: true,
			secure: process.env.NODE_ENV === 'production',
			sameSite: 'lax',
			maxAge: 60 * 60 * 24 * 7, // 7 days
		});

		return { success: true };
	},
};

export const load: PageServerLoad = async ({ locals }) => {
	validateEnvironmentVariables();

	if (isFlexibleBooleanTrue(env.BOOTSTRAP_MODE)) {
		throw redirect(302, resolve('/setup'));
	}

	if (locals.user) {
		throw redirect(302, resolve('/scopes'));
	}
};
