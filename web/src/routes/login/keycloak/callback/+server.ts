import { type RequestEvent } from '@sveltejs/kit';
import { decodeIdToken, type OAuth2Tokens } from 'arctic';

import { resolve } from '$app/paths';
import { env } from '$env/dynamic/private';
import { keycloak } from '$lib/server/keycloak';
import { createSession, generateSessionToken, setSessionTokenCookie } from '$lib/server/session';
import { createUser, getUser } from '$lib/server/user';

interface KeycloakIdTokenClaims {
	iss: string;
	sub: string;
	aud: string | string[];
	email?: string;
	name?: string;
	picture?: string;
	[key: string]: unknown;
}

function validateRequest(event: RequestEvent): { code: string; codeVerifier: string } | Response {
	const error = event.url.searchParams.get('error_description');
	if (error) {
		return new Response(error, {
			status: 400
		});
	}

	const code = event.url.searchParams.get('code');
	const state = event.url.searchParams.get('state');
	const storedState = event.cookies.get('OS_STATE') ?? null;
	const codeVerifier = event.cookies.get('OS_CODE_VERIFIER') ?? null;

	if (!code || !state || !storedState || !codeVerifier) {
		return new Response('Please restart the process.', {
			status: 400
		});
	}

	if (state !== storedState) {
		return new Response('Please restart the process.', {
			status: 400
		});
	}

	return { code, codeVerifier };
}

function validateClaims(claims: KeycloakIdTokenClaims): Response | null {
	if (!env.KEYCLOAK_REALM_URL || !env.KEYCLOAK_CLIENT_ID) {
		return new Response('Server misconfiguration.', {
			status: 500
		});
	}

	if (claims.iss !== env.KEYCLOAK_REALM_URL) {
		return new Response('Invalid issuer.', {
			status: 400
		});
	}

	const audienceValid =
		claims.aud === env.KEYCLOAK_CLIENT_ID ||
		(Array.isArray(claims.aud) && claims.aud.includes(env.KEYCLOAK_CLIENT_ID));

	if (!audienceValid) {
		return new Response('Invalid audience.', {
			status: 400
		});
	}

	return null;
}

export async function GET(event: RequestEvent): Promise<Response> {
	const validation = validateRequest(event);
	if (validation instanceof Response) return validation;

	const { code, codeVerifier } = validation;

	let tokens: OAuth2Tokens;
	try {
		tokens = await keycloak.validateAuthorizationCode(code, codeVerifier);
	} catch {
		return new Response('Invalid authorization code.', {
			status: 400
		});
	}

	const claims = decodeIdToken(tokens.idToken()) as KeycloakIdTokenClaims;
	const claimsError = validateClaims(claims);
	if (claimsError) return claimsError;

	const existingUser = await getUser(claims.sub);
	const user =
		existingUser ??
		(await createUser(
			claims.sub,
			(claims.preferred_username as string) ?? '',
			claims.email ?? '',
			claims.name ?? '',
			claims.picture ?? ''
		));

	const sessionToken = generateSessionToken();
	const session = await createSession(sessionToken, user.id);
	setSessionTokenCookie(event.cookies, sessionToken, session.expiresAt);

	return new Response(null, {
		status: 302,
		headers: { Location: resolve('/') }
	});
}
