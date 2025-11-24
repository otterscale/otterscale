import { error } from '@sveltejs/kit';

import { client } from '$lib/server/flagd';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const applicationEnabled = await client.getBooleanValue('application-enabled', false);
	if (!applicationEnabled) {
		throw error(501, `This feature is not implemented.`);
	}
};
