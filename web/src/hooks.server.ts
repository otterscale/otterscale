import type { Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';

import { env } from '$env/dynamic/private';
import { isFlexibleBooleanTrue } from '$lib/helper';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { keycloak } from '$lib/server/keycloak';
import {
	deleteSessionTokenCookie,
	setSessionTokenCookie,
	updateSessionTokenSet,
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
		event.locals.session = null;
		return resolve(event);
	}

	const session = await validateSessionToken(token);

	if (session) {
		setSessionTokenCookie(event.cookies, token, session.expiresAt);
		event.locals.session = session;

		try {
			const BUFFER_MS = 30 * 1000;
			const isNearExpiry =
				Date.now() >= (session.tokenSet.accessTokenExpiresAt.getTime() ?? 0) - BUFFER_MS;

			if (isNearExpiry) {
				const tokens = await keycloak.refreshAccessToken(session.tokenSet.refreshToken);
				const tokenSet = {
					accessToken: tokens.accessToken(),
					refreshToken: tokens.refreshToken(),
					accessTokenExpiresAt: tokens.accessTokenExpiresAt()
				};

				event.locals.session = await updateSessionTokenSet(session.id, tokenSet);
			}
		} catch (err) {
			console.error('Token refresh failed:', err);

			deleteSessionTokenCookie(event.cookies);
			event.locals.session = null;
		}
	} else {
		deleteSessionTokenCookie(event.cookies);
		event.locals.session = null;
	}

	return resolve(event);
};

const handleProxy: Handle = async ({ event, resolve }) => {
	if (isFlexibleBooleanTrue(env.BOOTSTRAP_MODE)) {
		return resolve(event);
	}

	const isApiProxy = event.request.headers.get('x-proxy-target') === 'api';
	const session = event.locals.session;

	if (!isApiProxy || !session) {
		return resolve(event);
	}

	const targetUrl = new URL(event.url.pathname + event.url.search, env.API_URL);

	const proxyHeaders = new Headers();
	const headersToForward = ['accept', 'content-type', 'user-agent'];

	event.request.headers.forEach((value, key) => {
		if (headersToForward.includes(key.toLowerCase())) {
			proxyHeaders.set(key, value);
		}
	});

	proxyHeaders.set('Authorization', `Bearer ${session.tokenSet.accessToken}`);

	try {
		const response = await fetch(targetUrl.toString(), {
			method: event.request.method,
			headers: proxyHeaders,
			body: event.request.body,
			duplex: 'half'
		} as RequestInit);

		const responseHeaders = new Headers(response.headers);
		const headersToClean = ['access-control-allow-origin', 'content-encoding', 'content-length'];

		headersToClean.forEach((h) => responseHeaders.delete(h));

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
