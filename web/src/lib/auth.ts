import Keycloak from 'keycloak-js';

import { isAuthenticated } from './stores/auth';

import { browser } from '$app/environment';
import { env } from '$env/dynamic/public';

// Initialize Keycloak instance
const keycloak = new Keycloak({
	url: env.PUBLIC_AUTH_URL ?? '',
	realm: env.PUBLIC_AUTH_REALM ?? '',
	clientId: env.PUBLIC_AUTH_CLIENT_ID ?? '',
});

// Setup token refresh handler
const setupTokenRefresh = () => {
	keycloak.onTokenExpired = () => {
		keycloak.updateToken(30).catch((error) => {
			console.error('Failed to refresh token:', error);
			isAuthenticated.set(false);
		});
	};
};

export const initializeAuth = async (): Promise<void> => {
	if (!browser) return;

	try {
		const authenticated = await keycloak.init({
			onLoad: 'login-required',
			pkceMethod: 'S256',
			checkLoginIframe: false,
			redirectUri: `${env.PUBLIC_URL}`,
		});

		isAuthenticated.set(authenticated);
		if (authenticated) {
			setupTokenRefresh();
		}
	} catch (error) {
		console.error('Error during Keycloak initialization:', error);
		isAuthenticated.set(false);
	}
};

export const login = (): void => {
	if (browser) {
		keycloak.login();
	}
};

export const logout = (): void => {
	if (browser) {
		keycloak.logout();
	}
};

export const getUser = (): User | undefined => keycloak.tokenParsed;

export const getToken = (): string | undefined => keycloak.token;

export type User = Keycloak.KeycloakTokenParsed;
