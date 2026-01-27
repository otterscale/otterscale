import { sha256 } from '@oslojs/crypto/sha2';
import { encodeBase32LowerCaseNoPadding, encodeHexLowerCase } from '@oslojs/encoding';
import type { Cookies } from '@sveltejs/kit';

import { env  } from '$env/dynamic/public';

import { redis } from './redis';

const COOKIE_NAME = isSecure() ? '__Host-OS_SESSION' : 'OS_SESSION';
const SESSION_EXPIRY_MS = 1000 * 60 * 60 * 24 * 30; // 30 days
const SESSION_REFRESH_THRESHOLD_MS = 1000 * 60 * 60 * 24 * 15; // 15 days

export async function acquireRefreshLock(sessionId: string, ttlMs: number): Promise<boolean> {
	const result = await redis.set(`refresh_lock:${sessionId}`, 'locked', 'PX', ttlMs, 'NX');
	return result === 'OK';
}

export async function releaseRefreshLock(sessionId: string): Promise<void> {
	await redis.del(`refresh_lock:${sessionId}`);
}

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

	await setSessionToRedis(session);

	return session;
}

export async function updateSessionTokenSet(sessionId: string, tokenSet: TokenSet) {
	await redis.hset(`session:${sessionId}`, {
		tokenSet: JSON.stringify(tokenSet)
	});
}

export async function validateSessionToken(
	token: string
): Promise<{ session: Session | null; fresh: boolean }> {
	const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));
	const session = await getSessionFromRedis(sessionId);

	if (session === null) {
		return { session: null, fresh: false };
	}

	if (Date.now() >= session.expiresAt.getTime()) {
		await redis.del(`session:${sessionId}`);
		return { session: null, fresh: false };
	}

	if (Date.now() >= session.expiresAt.getTime() - SESSION_REFRESH_THRESHOLD_MS) {
		session.expiresAt = new Date(Date.now() + SESSION_EXPIRY_MS);
		await setSessionToRedis(session);
		return { session: session, fresh: true };
	}

	return { session: session, fresh: false };
}

export async function invalidateSession(sessionId: string): Promise<void> {
	await redis.del(`session:${sessionId}`);
}

async function setSessionToRedis(session: Session) {
	const pipeline = redis.pipeline();
	pipeline.hset(`session:${session.id}`, {
		user: JSON.stringify(session.user),
		tokenSet: JSON.stringify(session.tokenSet),
		expiresAt: session.expiresAt.getTime().toString()
	});
	pipeline.expireat(`session:${session.id}`, Math.floor(session.expiresAt.getTime() / 1000));
	await pipeline.exec();
}

async function getSessionFromRedis(sessionId: string): Promise<Session | null> {
	const data = await redis.hgetall(`session:${sessionId}`);
	if (!data || Object.keys(data).length === 0) {
		return null;
	}

	try {
		const user = JSON.parse(data.user);
		const tokenSet = JSON.parse(data.tokenSet);
		tokenSet.accessTokenExpiresAt = new Date(tokenSet.accessTokenExpiresAt);
		const expiresAt = new Date(parseInt(data.expiresAt));

		return { id: sessionId, user, tokenSet, expiresAt };
	} catch {
		await invalidateSession(sessionId);
		return null;
	}
}

export function getSessionTokenCookie(cookies: Cookies): string | undefined {
	return cookies.get(COOKIE_NAME);
}

export function setSessionTokenCookie(cookies: Cookies, token: string, expiresAt: Date): void {
	cookies.set(COOKIE_NAME, token, {
		httpOnly: true,
		path: '/',
		secure: isSecure(),
		sameSite: 'lax',
		expires: expiresAt
	});
}

export function deleteSessionTokenCookie(cookies: Cookies): void {
	cookies.set(COOKIE_NAME, '', {
		httpOnly: true,
		path: '/',
		secure: isSecure(),
		sameSite: 'lax',
		maxAge: 0
	});
}

export function isSecure(): boolean {
	return env.PUBLIC_WEB_URL.startsWith('https');
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
	roles: string[];
};

export type TokenSet = {
	idToken: string;
	accessToken: string;
	refreshToken: string;
	accessTokenExpiresAt: Date;
};
