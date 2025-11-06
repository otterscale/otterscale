import Keycloak from 'keycloak-js';

import { env } from '$env/dynamic/public';
import { setToken } from '$lib/jwt';

const keycloak = new Keycloak({
	url: env.PUBLIC_AUTH_URL,
	realm: env.PUBLIC_AUTH_REALM,
	clientId: env.PUBLIC_AUTH_CLIENT_ID,
});

const handleTokenRefresh = async (): Promise<void> => {
	try {
		const refreshed = await keycloak.updateToken(30);
		if (refreshed && keycloak.token) {
			await setToken(keycloak.token);
		}
	} catch (error) {
		console.error('Failed to refresh token:', error);
	}
};

export const initializeAuth = async (): Promise<void> => {
	try {
		const authenticated = await keycloak.init({
			onLoad: 'check-sso',
			checkLoginIframe: false,
			silentCheckSsoRedirectUri: `${location.origin}/silent-check-sso.html`,
		});

		if (authenticated && keycloak.token) {
			await setToken(keycloak.token);
		}

		keycloak.onTokenExpired = handleTokenRefresh;
	} catch (error) {
		console.error('Error during Keycloak initialization:', error);
	}
};

export const register = (): Promise<void> => keycloak.register();

export const login = (): Promise<void> => keycloak.login();

export const logout = (): Promise<void> => keycloak.logout();
