import { error, redirect } from '@sveltejs/kit';

import { resolve } from '$app/paths';
import { deleteSessionTokenCookie, invalidateSession } from '$lib/server/session';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies }) => {
	if (!locals.session) {
		throw error(401, 'Unauthorized');
	}

	await invalidateSession(locals.session.id);
	deleteSessionTokenCookie(cookies);

	throw redirect(307, resolve('/'));
};
