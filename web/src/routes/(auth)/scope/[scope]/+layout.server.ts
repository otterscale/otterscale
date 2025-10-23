import { OpenFeature } from '@openfeature/server-sdk';
import { redirect } from '@sveltejs/kit';
import type { User } from 'better-auth';

import type { LayoutServerLoad } from './$types';

import { auth } from '$lib/auth';

export const load: LayoutServerLoad = async ({ request, url }) => {
	const session = await auth.api.getSession({
		headers: request.headers,
	});

	if (!session) {
		redirect(302, `/?next=${url.pathname}`);
	}

	const client = OpenFeature.getClient();

	const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));
	while (client.providerStatus !== 'READY') {
		if (Date.now() - Date.now() > 5000) {
			return {
				user: session.user as User,
				isAppGeneralOn: true,
				isAppHelmChartOn: false,
			};
		}
		await sleep(1000);
	}

	const appGeneral = await client.getBooleanValue('app-general', false);
	const appHelmChart = await client.getBooleanValue('app-helm-chart', false);
	const vmGeneral = await client.getBooleanValue('vm-general', false);
	const mdlGeneral = await client.getBooleanValue('mdl-general', false);
	const stgGeneral = await client.getBooleanValue('stg-general', false);
	const stgBlock = await client.getBooleanValue('stg-block', false);
	const stgFile = await client.getBooleanValue('stg-file', false);
	const stgObject = await client.getBooleanValue('stg-object', false);

	return {
		user: session.user as User,
		featureStates: {
			'app-general': appGeneral,
			'app-helm-chart': appHelmChart,
			'vm-general': vmGeneral,
			'mdl-general': mdlGeneral,
			'stg-general': stgGeneral,
			'stg-block': stgBlock,
			'stg-file': stgFile,
			'stg-object': stgObject,
		},
	};
};
