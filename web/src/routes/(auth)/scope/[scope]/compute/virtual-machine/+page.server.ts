import { error } from '@sveltejs/kit';

import { env } from '$env/dynamic/private';
import { client } from '$lib/server/flagd';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const virtualMachineEnabled = await client.getBooleanValue('virtual-machine-enabled', false);
	if (!virtualMachineEnabled) {
		throw error(501, `This feature is not implemented.`);
	}

	return {
		url: new URL(env.API_URL ?? '')
	};
};
