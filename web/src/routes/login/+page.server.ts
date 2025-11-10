import { redirect } from '@sveltejs/kit';

import { resolve } from '$app/paths';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	if (locals.session && locals.user) {
		throw redirect(307, resolve('/'));
	}

	return {
		user: locals.user
	};
};
