import { FlagdProvider } from '@openfeature/flagd-provider';
import { OpenFeature } from '@openfeature/server-sdk';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	try {
		await OpenFeature.setProviderAndWait(new FlagdProvider({}));
	} catch (error) {
		console.error('Failed to initialize provider:', error);
	}

	const client = OpenFeature.getClient();

	const orchGPUFeatureState = await client.getBooleanValue('orch-gpu', false);

	return {
		'feature-states.orch-gpu': orchGPUFeatureState
	};
};
