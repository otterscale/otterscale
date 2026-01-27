import { type Readable, readable } from 'svelte/store';

/**
 * User type representing a user from Keycloak or other identity provider
 */
export interface User {
	/** Unique identifier (e.g., OIDC subject or email) */
	subject: string;
	/** Human-readable display name */
	name: string;
	/** Optional email address */
	email?: string;
	/** Optional avatar URL */
	avatar?: string;
}

/**
 * Mock user data for testing
 * TODO: Replace with Keycloak API integration
 */
const mockUsers: User[] = [
	{
		subject: 'alice@example.com',
		name: 'Alice Chen',
		email: 'alice@example.com'
	},
	{
		subject: 'bob@example.com',
		name: 'Bob Wang',
		email: 'bob@example.com'
	},
	{
		subject: 'charlie@example.com',
		name: 'Charlie Lin',
		email: 'charlie@example.com'
	},
	{
		subject: 'diana@example.com',
		name: 'Diana Wu',
		email: 'diana@example.com'
	},
	{
		subject: 'evan@example.com',
		name: 'Evan Huang',
		email: 'evan@example.com'
	},
	{
		subject: 'fiona@example.com',
		name: 'Fiona Chang',
		email: 'fiona@example.com'
	},
	{
		subject: 'george@example.com',
		name: 'George Liu',
		email: 'george@example.com'
	},
	{
		subject: 'helen@example.com',
		name: 'Helen Yang',
		email: 'helen@example.com'
	}
];

/**
 * Readable store containing the list of available users
 * TODO: Replace with async fetch from Keycloak
 */
export const users: Readable<User[]> = readable(mockUsers);

/**
 * Fetch users from Keycloak (placeholder for future implementation)
 * @returns Promise resolving to array of users
 */
export async function fetchUsers(): Promise<User[]> {
	// TODO: Implement Keycloak API call
	// Example:
	// const response = await fetch('/api/keycloak/users');
	// return response.json();
	return mockUsers;
}

/**
 * Search users by name or email
 * @param query Search query string
 * @returns Filtered list of users matching the query
 */
export function searchUsers(query: string, userList: User[]): User[] {
	const lowerQuery = query.toLowerCase().trim();
	if (!lowerQuery) return userList;

	return userList.filter(
		(user) =>
			user.name.toLowerCase().includes(lowerQuery) ||
			user.subject.toLowerCase().includes(lowerQuery) ||
			user.email?.toLowerCase().includes(lowerQuery)
	);
}
