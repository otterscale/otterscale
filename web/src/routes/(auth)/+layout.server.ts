import { redirect } from '@sveltejs/kit';

import { resolve } from '$app/paths';

import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	if (!locals.session) {
		throw redirect(307, resolve('/'));
	}

	return {
		user: locals.session.user
	};
};
