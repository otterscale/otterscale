import { FlagdProvider } from '@openfeature/flagd-provider';
import { OpenFeature } from '@openfeature/server-sdk';

import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	try {
		await OpenFeature.setProviderAndWait(new FlagdProvider({ host: 'localhost', port: 8013 }));
	} catch (error) {
		console.error('Failed to initialize provider:', error);
	}

	const client = OpenFeature.getClient();

	const appContainerFeatureState = await client.getBooleanValue('app-container', false);

	return {
		'feature-states.app-container': appContainerFeatureState,
	};
};
