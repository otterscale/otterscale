import { KeyCloak } from 'arctic';

import { resolve } from '$app/paths';
import {
	KEYCLOAK_CLIENT_ID,
	KEYCLOAK_CLIENT_SECRET,
	KEYCLOAK_REALM_URL
} from '$env/static/private';
import { PUBLIC_WEB_URL } from '$env/static/public';

export const keycloak = new KeyCloak(
	KEYCLOAK_REALM_URL,
	KEYCLOAK_CLIENT_ID,
	KEYCLOAK_CLIENT_SECRET,
	`${PUBLIC_WEB_URL}${resolve('/login/keycloak/callback')}`
);
