import type { Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';

import { env } from '$env/dynamic/private';
import { isFlexibleBooleanTrue } from '$lib/helper';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { keycloak } from '$lib/server/keycloak';
import {
	deleteSessionTokenCookie,
	setSessionTokenCookie,
	updateSession,
	validateSessionToken
} from '$lib/server/session';

const handleParaglide: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request, locale }) => {
		event.request = request;

		return resolve(event, {
			transformPageChunk: ({ html }) => html.replace('%paraglide.lang%', locale)
		});
	});

const handleAuth: Handle = async ({ event, resolve }) => {
	if (isFlexibleBooleanTrue(env.BOOTSTRAP_MODE)) {
		return resolve(event);
	}

	const token = event.cookies.get('OS_SESSION');
	if (!token) {
		event.locals.user = null;
		event.locals.session = null;
		return resolve(event);
	}

	const { session, user } = await validateSessionToken(token);

	if (session) {
		setSessionTokenCookie(event.cookies, token, session.expiresAt);
		event.locals.session = session;
		event.locals.user = user;
	} else {
		deleteSessionTokenCookie(event.cookies);
		event.locals.session = null;
		event.locals.user = null;
	}

	return resolve(event);
};

const handleProxy: Handle = async ({ event, resolve }) => {
	const isApiProxy = event.request.headers.get('x-proxy-target') === 'api';
	const session = event.locals.session;

	if (!isApiProxy || !session?.accessToken) {
		return resolve(event);
	}

	try {
		const BUFFER_MS = 30 * 1000;
		const isNearExpiry = Date.now() >= (session.accessTokenExpiresAt?.getTime() ?? 0) - BUFFER_MS;

		if (isNearExpiry && session.refreshToken) {
			const tokens = await keycloak.refreshAccessToken(session.refreshToken);
			event.locals.session = await updateSession(
				session.id,
				tokens.accessToken(),
				tokens.accessTokenExpiresAt(),
				tokens.refreshToken()
			);
		}
	} catch (err) {
		console.error('Token refresh failed:', err);
		deleteSessionTokenCookie(event.cookies);
		return new Response('Session Expired', { status: 401 });
	}

	const targetUrl = new URL(event.url.pathname + event.url.search, env.API_URL);
	const proxyHeaders = new Headers(event.request.headers);

	proxyHeaders.delete('cookie');
	proxyHeaders.delete('x-proxy-target');
	proxyHeaders.set('Authorization', `Bearer ${event.locals.session?.accessToken}`);

	try {
		const response = await fetch(targetUrl.toString(), {
			method: event.request.method,
			headers: proxyHeaders,
			body: event.request.body,
			duplex: 'half'
		} as RequestInit);

		const responseHeaders = new Headers(response.headers);

		responseHeaders.delete('content-encoding');
		responseHeaders.delete('content-length');
		responseHeaders.delete('access-control-allow-origin');

		return new Response(response.body, {
			headers: responseHeaders,
			status: response.status,
			statusText: response.statusText
		} as ResponseInit);
	} catch (err) {
		console.error('Proxy Fetch Error:', err);
		return new Response(`Gateway Error: ${err}`, { status: 502 });
	}
};

export const handle: Handle = sequence(handleParaglide, handleAuth, handleProxy);
