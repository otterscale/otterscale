import { error } from '@sveltejs/kit';

import { client } from '$lib/server/flagd';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const modelEnabled = await client.getBooleanValue('model-enabled', false);
	if (!modelEnabled) {
		throw error(501, `This feature is not implemented.`);
	}
};
