import { type RequestEvent } from '@sveltejs/kit';
import { decodeIdToken, type OAuth2Tokens } from 'arctic';

import { resolve } from '$app/paths';
import { env } from '$env/dynamic/private';
import { keycloak } from '$lib/server/keycloak';
import { createSession, generateSessionToken, setSessionTokenCookie } from '$lib/server/session';

interface KeycloakIdTokenClaims {
	iss: string;
	sub: string;
	aud: string | string[];
	email?: string;
	name?: string;
	picture?: string;
	preferred_username?: string;
	resource_access?: {
		[key: string]: {
			roles?: string[];
		};
	};
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

	const token = generateSessionToken();
	const user = {
		sub: claims.sub,
		username: claims.preferred_username ?? '',
		name: claims.name ?? '',
		email: claims.email ?? '',
		picture: claims.picture ?? '',
		roles: claims.resource_access?.[env.KEYCLOAK_CLIENT_ID ?? '']?.roles ?? []
	};
	console.log(claims.sub);
	console.log(tokens.accessToken());
	const tokenSet = {
		idToken: tokens.idToken(),
		accessToken: tokens.accessToken(),
		accessTokenExpiresAt: tokens.accessTokenExpiresAt(),
		refreshToken: tokens.refreshToken()
	};

	const session = await createSession(token, user, tokenSet);
	setSessionTokenCookie(event.cookies, token, session.expiresAt);

	return new Response(null, {
		status: 302,
		headers: { Location: resolve('/') }
	});
}
