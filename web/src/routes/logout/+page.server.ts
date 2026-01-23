import { error, redirect } from '@sveltejs/kit';

import { env } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';
import { deleteSessionTokenCookie, invalidateSession } from '$lib/server/session';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies }) => {
	if (!locals.session) {
		throw error(401, 'Unauthorized');
	}

	await invalidateSession(locals.session.id);
	deleteSessionTokenCookie(cookies);

	throw redirect(
		302,
		`${env.KEYCLOAK_REALM_URL}/protocol/openid-connect/logout?id_token_hint=${locals.session.tokenSet.idToken}&post_logout_redirect_uri=${publicEnv.PUBLIC_WEB_URL}`
	);
};
