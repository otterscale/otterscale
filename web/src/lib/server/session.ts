import { sha256 } from '@oslojs/crypto/sha2';
import { encodeBase32, encodeHexLowerCase } from '@oslojs/encoding';
import type { Cookies } from '@sveltejs/kit';

import { dev } from '$app/environment';

import { db } from './db';
import type { User } from './user';

export function validateSessionToken(token: string): SessionValidationResult {
	const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));
	const row = db.queryOne(
		`
SELECT session.id, session.user_id, session.expires_at, user.id, user.sub, user.email, user.name, user.picture FROM session
INNER JOIN user ON session.user_id = user.id
WHERE session.id = ?
`,
		[sessionId]
	);

	if (row === null) {
		return { session: null, user: null };
	}
	const session: Session = {
		id: row.string(0),
		userId: row.number(1),
		expiresAt: new Date(row.number(2) * 1000)
	};
	const user: User = {
		id: row.number(3),
		sub: row.string(4),
		email: row.string(5),
		name: row.string(6),
		picture: row.string(7)
	};
	if (Date.now() >= session.expiresAt.getTime()) {
		db.execute('DELETE FROM session WHERE id = ?', [session.id]);
		return { session: null, user: null };
	}
	if (Date.now() >= session.expiresAt.getTime() - 1000 * 60 * 60 * 24 * 15) {
		session.expiresAt = new Date(Date.now() + 1000 * 60 * 60 * 24 * 30);
		db.execute('UPDATE session SET expires_at = ? WHERE session.id = ?', [
			Math.floor(session.expiresAt.getTime() / 1000),
			session.id
		]);
	}
	return { session, user };
}

export function invalidateSession(sessionId: string): void {
	db.execute('DELETE FROM session WHERE id = ?', [sessionId]);
}

// TODO: TTL
export function invalidateUserSessions(userId: number): void {
	db.execute('DELETE FROM session WHERE user_id = ?', [userId]);
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

export function generateSessionToken(): string {
	const tokenBytes = new Uint8Array(20);
	crypto.getRandomValues(tokenBytes);
	const token = encodeBase32(tokenBytes).toLowerCase();
	return token;
}

export function createSession(token: string, userId: number): Session {
	const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));
	const session: Session = {
		id: sessionId,
		userId,
		expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24 * 30)
	};
	db.execute('INSERT INTO session (id, user_id, expires_at) VALUES (?, ?, ?)', [
		session.id,
		session.userId,
		Math.floor(session.expiresAt.getTime() / 1000)
	]);
	return session;
}

export interface Session {
	id: string;
	userId: number;
	expiresAt: Date;
}

type SessionValidationResult = { session: Session; user: User } | { session: null; user: null };
