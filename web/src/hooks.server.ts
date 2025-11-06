import type { Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';

import { verifyToken } from '$lib/jwt';
import { paraglideMiddleware } from '$lib/paraglide/server';

const handleParaglide: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request, locale }) => {
		event.request = request;

		return resolve(event, {
			transformPageChunk: ({ html }) => html.replace('%paraglide.lang%', locale),
		});
	});

const handleAuth: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('OS_TOKEN');

	if (token) {
		const payload = await verifyToken(token);
		event.locals.user = payload || null;

		if (!payload) {
			event.cookies.delete('OS_TOKEN', { path: '/' });
		}
	} else {
		event.locals.user = null;
	}

	return resolve(event);
};

export const handle = sequence(handleParaglide, handleAuth);
