import { KeyCloak } from 'arctic';

import { resolve } from '$app/paths';
import { env } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';

export const keycloak = new KeyCloak(
	env.KEYCLOAK_REALM_URL ?? '',
	env.KEYCLOAK_CLIENT_ID ?? '',
	env.KEYCLOAK_CLIENT_SECRET ?? '',
	`${publicEnv.PUBLIC_WEB_URL}${resolve('/login/callback')}`
);
