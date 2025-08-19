import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import { auth } from '$lib/auth';
import { staticPaths } from '$lib/path';
import type { PageServerLoad } from './$types';

const isProviderConfigured = (clientId: any, clientSecret: any, ...additionalKeys: any[]) => {
	return Boolean(clientId && clientSecret && additionalKeys.every((key) => key));
};

export const load: PageServerLoad = async ({ request, url }) => {
	const session = await auth.api.getSession({
		headers: request.headers
	});

	if (session) {
		redirect(302, staticPaths.scopes.url);
	}

	const nextPath = url.searchParams.get('next') || staticPaths.scopes.url;

	return {
		nextPath,
		apple: isProviderConfigured(
			env.APPLE_CLIENT_ID,
			env.APPLE_CLIENT_SECRET,
			env.APPLE_APP_BUNDLE_IDENTIFIER
		),
		github: isProviderConfigured(env.GITHUB_CLIENT_ID, env.GITHUB_CLIENT_SECRET),
		google: isProviderConfigured(env.GOOGLE_CLIENT_ID, env.GOOGLE_CLIENT_SECRET),
		oidcProvider: env.AUTH_OIDC_PROVIDER,
		ssoLoginPrompt: env.SSO_LOGIN_PROMPT
	};
};
