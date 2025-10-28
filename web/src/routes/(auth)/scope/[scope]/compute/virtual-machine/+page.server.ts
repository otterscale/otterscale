import { FlagdProvider } from '@openfeature/flagd-provider';
import { OpenFeature } from '@openfeature/server-sdk';
import { error } from '@sveltejs/kit';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	try {
		await OpenFeature.setProviderAndWait(new FlagdProvider({ host: 'localhost', port: 8013 }));
	} catch (error) {
		console.error('Failed to initialize provider:', error);
	}

	const client = OpenFeature.getClient();

	const vmGeneralFeatureState = await client.getBooleanValue('vm-general', false);

	if (!vmGeneralFeatureState) {
		throw error(501, `This feature is not implemented.`);
	}
};
