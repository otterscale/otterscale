import { client } from '$lib/server/flagd';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const modelEnabled = await client.getBooleanValue('model-enabled', false);
	return {
		'model-enabled': modelEnabled
	};
};
