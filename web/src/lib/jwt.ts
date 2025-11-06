import { jwtVerify, createRemoteJWKSet, type JWTPayload } from 'jose';

import { env } from '$env/dynamic/public';

export type User = JWTPayload;

const getJwksUri = () => `${env.PUBLIC_AUTH_URL}/realms/${env.PUBLIC_AUTH_REALM}/protocol/openid-connect/certs`;

const getIssuer = () => `${env.PUBLIC_AUTH_URL}/realms/${env.PUBLIC_AUTH_REALM}`;

export const setToken = async (token: string): Promise<void> => {
	const response = await fetch('/?/setToken', {
		method: 'POST',
		headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
		body: new URLSearchParams({
			jwksUri: getJwksUri(),
			token,
		}),
	});

	if (!response.ok) {
		throw new Error(`Failed to set token: ${response.statusText}`);
	}
};

export const verifyToken = async (token: string): Promise<User | null> => {
	try {
		const jwksUrl = new URL(getJwksUri());
		const remoteJWKSet = createRemoteJWKSet(jwksUrl);

		const { payload } = await jwtVerify(token, remoteJWKSet, {
			issuer: getIssuer(),
			audience: env.PUBLIC_AUTH_CLIENT_ID,
		});
		return payload;
	} catch (err) {
		console.error('Token verification failed:', err);
		return null;
	}
};
