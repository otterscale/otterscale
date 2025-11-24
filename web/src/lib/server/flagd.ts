import { FlagdProvider } from '@openfeature/flagd-provider';
import { OpenFeature } from '@openfeature/server-sdk';

let client: ReturnType<typeof OpenFeature.getClient>;

try {
	await OpenFeature.setProviderAndWait(new FlagdProvider({}));
} catch (error) {
	console.error('Failed to initialize provider:', error);
} finally {
	client = OpenFeature.getClient();
}

export { client };
