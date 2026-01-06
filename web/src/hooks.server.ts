import { type Handle, redirect } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';

import { env } from '$env/dynamic/private';
import { isFlexibleBooleanTrue } from '$lib/helper';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { keycloak } from '$lib/server/keycloak';
import {
	acquireRefreshLock,
	deleteSessionTokenCookie,
	getSessionTokenCookie,
	releaseRefreshLock,
	setSessionTokenCookie,
	updateSessionTokenSet,
	validateSessionToken
} from '$lib/server/session';

const HOP_BY_HOP_HEADERS = [
	'connection',
	'keep-alive',
	'proxy-authenticate',
	'proxy-authorization',
	'te',
	'trailer',
	'transfer-encoding',
	'upgrade'
];

const PROXY_REQUEST_HEADERS_TO_REMOVE = ['cookie', 'host', 'x-proxy-target'];

const PROXY_RESPONSE_HEADERS_TO_REMOVE = ['content-encoding', 'content-length'];

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

	const token = getSessionTokenCookie(event.cookies);

	if (!token) {
		event.locals.session = null;
		return resolve(event);
	}

	const { session, fresh } = await validateSessionToken(token);

	if (session) {
		if (fresh) {
			setSessionTokenCookie(event.cookies, token, session.expiresAt);
		}
		event.locals.session = session;
	} else {
		deleteSessionTokenCookie(event.cookies);
		event.locals.session = null;
	}

	return resolve(event);
};

const handleRefreshToken: Handle = async ({ event, resolve }) => {
	if (isFlexibleBooleanTrue(env.BOOTSTRAP_MODE)) {
		return resolve(event);
	}

	const session = event.locals.session;

	if (session) {
		const BUFFER_MS = 60 * 1000;
		const isNearExpiry =
			Date.now() >= (session.tokenSet.accessTokenExpiresAt.getTime() ?? 0) - BUFFER_MS;

		if (isNearExpiry) {
			const REFRESH_LOCK_TTL_MS = 10 * 1000;
			const hasLock = await acquireRefreshLock(session.id, REFRESH_LOCK_TTL_MS);

			// stale-while-revalidate
			if (hasLock) {
				try {
					const tokens = await keycloak.refreshAccessToken(session.tokenSet.refreshToken);
					const tokenSet = {
						accessToken: tokens.accessToken(),
						refreshToken: tokens.refreshToken(),
						accessTokenExpiresAt: tokens.accessTokenExpiresAt()
					};

					await updateSessionTokenSet(session.id, tokenSet);

					session.tokenSet = tokenSet;
					event.locals.session = session;
				} catch (err) {
					console.error('Token refresh failed:', err);

					deleteSessionTokenCookie(event.cookies);
					event.locals.session = null;
				} finally {
					await releaseRefreshLock(session.id);
				}
			}
		}
	}

	return resolve(event);
};

const handleGuard: Handle = async ({ event, resolve }) => {
	if (isFlexibleBooleanTrue(env.BOOTSTRAP_MODE) || event.locals.session) {
		return resolve(event);
	}

	const isApiProxy = event.request.headers.get('x-proxy-target') === 'api';

	if (isApiProxy) {
		return new Response(JSON.stringify({ error: 'Unauthorized' }), {
			status: 401,
			headers: { 'Content-Type': 'application/json' }
		});
	}

	const isPrivatePath = event.route.id?.startsWith('/(auth)/');

	if (isPrivatePath) {
		throw redirect(303, '/login');
	}

	return resolve(event);
};

const handleProxy: Handle = async ({ event, resolve }) => {
	const isApiProxy = event.request.headers.get('x-proxy-target') === 'api';
	const session = event.locals.session;

	if (!isApiProxy || !session) {
		return resolve(event);
	}

	const targetUrl = new URL(event.url.pathname + event.url.search, env.API_URL);
	const proxyHeaders = new Headers(event.request.headers);

	HOP_BY_HOP_HEADERS.forEach((header) => proxyHeaders.delete(header));
	PROXY_REQUEST_HEADERS_TO_REMOVE.forEach((header) => proxyHeaders.delete(header));

	proxyHeaders.set('Authorization', `Bearer ${session.tokenSet.accessToken}`);

	try {
		const response = await fetch(targetUrl.toString(), {
			method: event.request.method,
			headers: proxyHeaders,
			body: event.request.body,
			duplex: 'half'
		} as RequestInit);

		const responseHeaders = new Headers(response.headers);

		HOP_BY_HOP_HEADERS.forEach((header) => responseHeaders.delete(header));
		PROXY_RESPONSE_HEADERS_TO_REMOVE.forEach((header) => responseHeaders.delete(header));

		return new Response(response.body, {
			headers: responseHeaders,
			status: response.status,
			statusText: response.statusText
		} as ResponseInit);
	} catch (err) {
		console.error('Proxy Fetch Error:', err);
		return new Response(JSON.stringify({ error: 'Bad Gateway', details: 'Upstream unreachable' }), {
			status: 502,
			headers: { 'Content-Type': 'application/json' }
		});
	}
};

export const handle: Handle = sequence(
	handleParaglide,
	handleAuth,
	handleRefreshToken,
	handleGuard,
	handleProxy
);
