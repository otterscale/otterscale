import Keycloak from 'keycloak-js';

import { isAuthenticated, token, user } from './stores';

import { env } from '$env/dynamic/public';

const keycloak = new Keycloak({
	url: env.PUBLIC_AUTH_URL ?? '',
	realm: env.PUBLIC_AUTH_REALM ?? '',
	clientId: env.PUBLIC_AUTH_CLIENT_ID ?? '',
});

const reset = () => {
	isAuthenticated.set(false);
	user.set(undefined);
	token.set(undefined);
};

export const initializeAuth = async (): Promise<void> => {
	try {
		const authenticated = await keycloak.init({
			onLoad: 'check-sso',
			checkLoginIframe: false,
			silentCheckSsoRedirectUri: `${location.origin}/silent-check-sso.html`,
		});

		isAuthenticated.set(authenticated);

		if (authenticated) {
			user.set(keycloak.tokenParsed);
			token.set(keycloak.token);
		}

		keycloak.onTokenExpired = () => {
			keycloak.updateToken(30).then((refreshed) => {
				if (refreshed) {
					token.set(keycloak.token);
					user.set(keycloak.tokenParsed);
				}
			});
		};
	} catch (error) {
		console.error('Error during Keycloak initialization:', error);
		reset();
	}
};

export const login = () => {
	keycloak.login();
	reset();
};

export const logout = () => {
	keycloak.logout();
	reset();
};

export type User = Keycloak.KeycloakTokenParsed;
