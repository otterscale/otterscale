import { FlagdProvider } from '@openfeature/flagd-provider';
import { OpenFeature } from '@openfeature/server-sdk';

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

async function getFeatureStates(flags: string[]) {
	const featureStates = Object.fromEntries(flags.map((flag) => [flag, false]));

	OpenFeature.setProvider(new FlagdProvider({ host: 'localhost', port: 8013 }));

	const client = OpenFeature.getClient();
	const connectTime = Date.now();

	while (client.providerStatus !== 'READY') {
		if (Date.now() - connectTime <= 3000) {
			console.log('Retry to access features states');
			await sleep(1000);
		} else {
			return featureStates;
		}
	}

	flags.forEach((flag) => {
		client.getBooleanValue(flag, false).then((value) => {
			featureStates[flag] = value;
		});
	});
	console.log(featureStates);
	return featureStates;
}

export { getFeatureStates };
