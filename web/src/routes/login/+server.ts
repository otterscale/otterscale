import type { RequestEvent } from '@sveltejs/kit';
import { generateCodeVerifier, generateState } from 'arctic';

import { dev } from '$app/environment';
import { keycloak } from '$lib/server/keycloak';

const COOKIE_OPTIONS = {
	path: '/',
	httpOnly: true,
	secure: !dev,
	maxAge: 60 * 10, // 10 minutes
	sameSite: 'lax' as const
};

export async function GET(event: RequestEvent): Promise<Response> {
	const state = generateState();
	const codeVerifier = generateCodeVerifier();
	const url = keycloak.createAuthorizationURL(state, codeVerifier, ['openid', 'profile', 'email']);

	event.cookies.set('OS_STATE', state, COOKIE_OPTIONS);
	event.cookies.set('OS_CODE_VERIFIER', codeVerifier, COOKIE_OPTIONS);

	return new Response(null, {
		status: 307,
		headers: { Location: url.toString() }
	});
}
