import { env } from '$env/dynamic/private';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	return {
		url: env.API_URL
	};
};
