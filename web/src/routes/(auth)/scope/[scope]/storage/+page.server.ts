import { FlagdProvider } from '@openfeature/flagd-provider';
import { OpenFeature } from '@openfeature/server-sdk';
import { error } from '@sveltejs/kit';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	try {
		await OpenFeature.setProviderAndWait(new FlagdProvider({}));
	} catch (error) {
		console.error('Failed to initialize provider:', error);
	}

	const client = OpenFeature.getClient();

	const stgBlockFeatureState = await client.getBooleanValue('stg-block', false);
	const stgFileFeatureState = await client.getBooleanValue('stg-file', false);
	const stgObjectFeatureState = await client.getBooleanValue('stg-object', false);
	const stgGeneralFeatureState = await client.getBooleanValue('stg-general', false);

	if (!(stgBlockFeatureState && stgFileFeatureState && stgObjectFeatureState && stgGeneralFeatureState)) {
		throw error(501, `This feature is not implemented.`);
	}
};
