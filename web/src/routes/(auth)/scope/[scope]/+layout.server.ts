import { redirect } from '@sveltejs/kit';
import type { User } from 'better-auth';

import type { LayoutServerLoad } from './$types';

import { auth } from '$lib/auth';

export const load: LayoutServerLoad = async ({ request, url }) => {
	const session = await auth.api.getSession({
		headers: request.headers,
	});

	if (!session) {
		redirect(302, `/?next=${url.pathname}`);
	}

	return {
		user: session.user as User,
	};
};
