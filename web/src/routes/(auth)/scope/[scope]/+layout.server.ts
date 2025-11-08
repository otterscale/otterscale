import { FlagdProvider } from '@openfeature/flagd-provider';
import { OpenFeature } from '@openfeature/server-sdk';
import { redirect } from '@sveltejs/kit';
import type { User } from 'better-auth';

import { auth } from '$lib/auth';

import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ request, url }) => {
	const session = await auth.api.getSession({
		headers: request.headers
	});

	if (!session) {
		redirect(302, `/?next=${url.pathname}`);
	}

	try {
		await OpenFeature.setProviderAndWait(new FlagdProvider({}));
	} catch (error) {
		console.error('Failed to initialize provider:', error);
	}

	const client = OpenFeature.getClient();

	const appGeneralFeatureState = await client.getBooleanValue('app-general', false);
	const appHelmChartFeatureState = await client.getBooleanValue('app-helm-chart', false);
	const vmGeneralFeatureState = await client.getBooleanValue('vm-general', false);
	const mdlGeneralFeatureState = await client.getBooleanValue('mdl-general', false);
	const stgGeneralFeatureState = await client.getBooleanValue('stg-general', false);
	const stgBlockFeatureState = await client.getBooleanValue('stg-block', false);
	const stgFileFeatureState = await client.getBooleanValue('stg-file', false);
	const stgObjectFeatureState = await client.getBooleanValue('stg-object', false);

	return {
		user: session.user as User,
		'feature-states.app-general': appGeneralFeatureState,
		'feature-states.app-helm-chart': appHelmChartFeatureState,
		'feature-states.vm-general': vmGeneralFeatureState,
		'feature-states.mdl-general': mdlGeneralFeatureState,
		'feature-states.stg-general': stgGeneralFeatureState,
		'feature-states.stg-block': stgBlockFeatureState,
		'feature-states.stg-file': stgFileFeatureState,
		'feature-states.stg-object': stgObjectFeatureState
	};
};
