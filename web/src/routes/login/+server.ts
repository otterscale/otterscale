import type { RequestEvent } from '@sveltejs/kit';
import { generateCodeVerifier, generateState } from 'arctic';

import { keycloak } from '$lib/server/keycloak';
import { isSecure } from '$lib/server/session';

export async function GET(event: RequestEvent): Promise<Response> {
	const state = generateState();
	const codeVerifier = generateCodeVerifier();
	const url = keycloak.createAuthorizationURL(state, codeVerifier, ['openid', 'profile', 'email']);

	event.cookies.set('OS_STATE', state, {
		httpOnly: true,
		maxAge: 60 * 10, // 10 minutes
		path: '/',
		sameSite: 'lax' as const,
		secure: isSecure()
	});

	event.cookies.set('OS_CODE_VERIFIER', codeVerifier, {
		httpOnly: true,
		maxAge: 60 * 10, // 10 minutes
		path: '/',
		sameSite: 'lax' as const,
		secure: isSecure()
	});

	return new Response(null, {
		status: 307,
		headers: { Location: url.toString() }
	});
}
