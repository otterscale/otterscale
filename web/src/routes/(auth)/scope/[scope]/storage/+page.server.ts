import { error } from '@sveltejs/kit';

import { client } from '$lib/server/flagd';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const distributedStorageEnabled = await client.getBooleanValue(
		'distributed-storage-enabled',
		false
	);
	if (!distributedStorageEnabled) {
		throw error(501, `This feature is not implemented.`);
	}
};
