import type { Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';

import { env } from '$env/dynamic/private';
import { isFlexibleBooleanTrue } from '$lib/helper';
import { paraglideMiddleware } from '$lib/paraglide/server';
import {
	deleteSessionTokenCookie,
	setSessionTokenCookie,
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

	const token = event.cookies.get('OS_SESSION') ?? null;
	if (!token) {
		event.locals.user = null;
		event.locals.session = null;
		return resolve(event);
	}

	const { session, user } = await validateSessionToken(token);
	if (session) {
		setSessionTokenCookie(event.cookies, token, session.expiresAt);
	} else {
		deleteSessionTokenCookie(event.cookies);
	}

	event.locals.session = session;
	event.locals.user = user;
	return resolve(event);
};

export const handle: Handle = sequence(handleParaglide, handleAuth);
