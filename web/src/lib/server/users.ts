import { env } from '$env/dynamic/private';

interface ClientCredentialsTokens {
	access_token: string;
	expires_in: number;
	refresh_expires_in: number;
	token_type: string;
	'not-before-policy': number;
	scope: string;
}

export async function getClientCredentialsTokens(): Promise<ClientCredentialsTokens> {
	if (!env.KEYCLOAK_ADMIN_REALM_URL) {
		console.error('KEYCLOAK_ADMIN_REALM_URL is not configured');
		throw new Error('Keycloak admin realm URL is not configured');
	}

	const response = await fetch(`${env.KEYCLOAK_REALM_URL}/protocol/openid-connect/token`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded'
		},
		body: new URLSearchParams({
			grant_type: 'client_credentials',
			client_id: env.KEYCLOAK_CLIENT_ID ?? '',
			client_secret: env.KEYCLOAK_CLIENT_SECRET ?? ''
		})
	});

	if (!response.ok) {
		throw new Error('Failed to get client credentials tokens');
	}

	const tokens = await response.json();
	return tokens;
}

export interface GetUsersOptions {
	search?: string;
	first?: number;
	max?: number;
}

export interface User {
	id: string;
	createdTimestamp: number;
	username: string;
	enabled: boolean;
	totp: boolean;
	emailVerified: boolean;
	firstName?: string;
	lastName?: string;
	email?: string;
	disableableCredentialTypes?: string[];
	requiredActions?: string[];
	notBefore?: number;
	access?: {
		manageGroupMembership: boolean;
		view: boolean;
		mapRoles: boolean;
		impersonate: boolean;
		manage: boolean;
	};
	attributes?: Record<string, string[]>;
}

export async function getUsers(options: GetUsersOptions = {}): Promise<User[]> {
	const { search = '', first = 0, max = 10 } = options;
	const tokens = await getClientCredentialsTokens();

	if (!env.KEYCLOAK_ADMIN_REALM_URL) {
		console.error('KEYCLOAK_ADMIN_REALM_URL is not configured');
		throw new Error('Keycloak admin realm URL is not configured');
	}

	if (!tokens.access_token) {
		console.error('Access token is missing');
		throw new Error('Access token is missing');
	}

	const params = new URLSearchParams({
		first: first.toString(),
		max: max.toString()
	});

	if (search) {
		// Use 'search' parameter for prefix matching (starts-with) across username, firstName, lastName, email
		// This supports Chinese characters and other Unicode text
		params.set('search', search);
	}

	const response = await fetch(`${env.KEYCLOAK_ADMIN_REALM_URL}/users?${params.toString()}`, {
		headers: {
			Authorization: `Bearer ${tokens.access_token}`
		}
	});

	// console.log(response);

	if (!response.ok) {
		const errorText = await response.text();
		console.error(
			`Failed to get users from Keycloak: ${response.status} ${response.statusText}`,
			errorText
		);
		throw new Error('Failed to get users from Keycloak');
	}

	const users = await response.json();
	return users;
}
