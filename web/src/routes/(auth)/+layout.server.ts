import { redirect } from '@sveltejs/kit';

import type { LayoutServerLoad } from './$types';

import { resolve } from '$app/paths';

export const load: LayoutServerLoad = async ({ locals }) => {
	if (!locals.user) {
		throw redirect(302, resolve('/'));
	}

	return {
		user: locals.user,
	};
};
