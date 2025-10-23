import { FlagdProvider } from '@openfeature/flagd-provider';
import { OpenFeature } from '@openfeature/server-sdk';

import type { PageServerLoad } from './$types';

OpenFeature.setProvider(new FlagdProvider({ host: 'localhost', port: 8013 }));

export const load: PageServerLoad = async () => {
	const client = OpenFeature.getClient();

	const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));
	while (client.providerStatus !== 'READY') {
		if (Date.now() - Date.now() > 5000) {
			return { enabled: null, providerStatus: client.providerStatus };
		}
		await sleep(100);
	}

	const orchGPU = await client.getBooleanValue('orch-gpu', false);

	return {
		orchestratorFeatureStates: {
			'orch-gpu': orchGPU,
		},
	};
};
