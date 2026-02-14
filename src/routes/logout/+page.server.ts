import { redirect } from '@sveltejs/kit';

import { env } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';
import { deleteSessionTokenCookie, invalidateSession } from '$lib/server/session';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies }) => {
	if (!locals.session) {
		throw redirect(302, '/');
	}

	await invalidateSession(locals.session.id);
	deleteSessionTokenCookie(cookies);

	const logoutUrl = new URL(`${env.KEYCLOAK_REALM_URL}/protocol/openid-connect/logout`);
	logoutUrl.searchParams.set('id_token_hint', locals.session.tokenSet.idToken);
	logoutUrl.searchParams.set('post_logout_redirect_uri', publicEnv.PUBLIC_WEB_URL ?? '');

	throw redirect(302, logoutUrl.toString());
};
