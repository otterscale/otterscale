import { error, json } from '@sveltejs/kit';

import { getUsers } from '$lib/server/users';

import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url, locals }) => {
	if (!locals.session) {
		error(401, 'Unauthorized');
	}

	const search = url.searchParams.get('search') || '';
	const first = parseInt(url.searchParams.get('first') || '0', 10);
	const max = parseInt(url.searchParams.get('max') || '10', 10);

	const users = await getUsers({ search, first, max });
	return json(users);
};
