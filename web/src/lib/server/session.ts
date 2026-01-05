import { sha256 } from '@oslojs/crypto/sha2';
import { encodeBase32LowerCaseNoPadding, encodeHexLowerCase } from '@oslojs/encoding';
import type { Cookies } from '@sveltejs/kit';

import { dev } from '$app/environment';

import { redis } from './redis';

const SESSION_EXPIRY_MS = 1000 * 60 * 60 * 24 * 30; // 30 days
const SESSION_REFRESH_THRESHOLD_MS = 1000 * 60 * 60 * 24 * 15; // 15 days

export function generateSessionToken(): string {
	const bytes = new Uint8Array(20);
	crypto.getRandomValues(bytes);
	return encodeBase32LowerCaseNoPadding(bytes);
}

export async function createSession(
	token: string,
	user: User,
	tokenSet: TokenSet
): Promise<Session> {
	const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));

	const session: Session = {
		id: sessionId,
		user,
		tokenSet,
		expiresAt: new Date(Date.now() + SESSION_EXPIRY_MS)
	};

	await saveSessionToRedis(session);

	return session;
}

export async function updateSessionTokenSet(
	sessionId: string,
	tokenSet: TokenSet
): Promise<Session | null> {
	const item = await redis.get(`session:${sessionId}`);

	if (item === null) {
		return null;
	}

	const result = JSON.parse(item);

	const session: Session = {
		id: result.id,
		user: result.user,
		tokenSet: tokenSet,
		expiresAt: result.expiresAt
	};

	await saveSessionToRedis(session);

	return session;
}

export async function validateSessionToken(token: string): Promise<Session | null> {
	const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));
	const item = await redis.get(`session:${sessionId}`);

	if (item === null) {
		return null;
	}

	const result = JSON.parse(item);

	const session: Session = {
		id: result.id,
		user: result.user,
		tokenSet: {
			accessToken: result.tokenSet.accessToken,
			refreshToken: result.tokenSet.refreshToken,
			accessTokenExpiresAt: new Date(result.tokenSet.accessTokenExpiresAt)
		},
		expiresAt: new Date(result.expiresAt)
	};

	if (Date.now() >= session.expiresAt.getTime()) {
		await redis.del(`session:${sessionId}`);
		return null;
	}

	if (Date.now() >= session.expiresAt.getTime() - SESSION_REFRESH_THRESHOLD_MS) {
		session.expiresAt = new Date(Date.now() + SESSION_EXPIRY_MS);
		await saveSessionToRedis(session);
	}

	return session;
}

export async function invalidateSession(sessionId: string): Promise<void> {
	await redis.del(`session:${sessionId}`);
}

async function saveSessionToRedis(session: Session) {
	await redis.set(
		`session:${session.id}`,
		JSON.stringify({
			id: session.id,
			user: session.user,
			tokenSet: session.tokenSet,
			expiresAt: session.expiresAt.getTime()
		}),
		'EXAT',
		Math.floor(session.expiresAt.getTime() / 1000)
	);
}

export function setSessionTokenCookie(cookies: Cookies, token: string, expiresAt: Date): void {
	cookies.set('OS_SESSION', token, {
		httpOnly: true,
		path: '/',
		secure: !dev,
		sameSite: 'lax',
		expires: expiresAt
	});
}

export function deleteSessionTokenCookie(cookies: Cookies): void {
	cookies.set('OS_SESSION', '', {
		httpOnly: true,
		path: '/',
		secure: !dev,
		sameSite: 'lax',
		maxAge: 0
	});
}

// Types
export type Session = {
	id: string;
	user: User;
	tokenSet: TokenSet;
	expiresAt: Date;
};

export type User = {
	sub: string;
	username: string;
	name: string;
	email: string;
	picture: string;
};

export type TokenSet = {
	accessToken: string;
	refreshToken: string;
	accessTokenExpiresAt: Date;
};
