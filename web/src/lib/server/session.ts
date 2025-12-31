import { sha256 } from '@oslojs/crypto/sha2';
import { encodeBase32, encodeHexLowerCase } from '@oslojs/encoding';
import type { Cookies } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';

import { dev } from '$app/environment';

import { db } from './db';
import { sessionsTable, usersTable } from './db/schema';
import type { User } from './user';

export async function validateSessionToken(token: string): Promise<SessionValidationResult> {
	const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));

	const sessions = await db
		.select()
		.from(sessionsTable)
		.innerJoin(usersTable, eq(sessionsTable.userId, usersTable.id))
		.where(eq(sessionsTable.id, sessionId));

	if (sessions.length === 0) {
		return { session: null, user: null };
	}

	const session = sessions[0].sessions;
	const user = sessions[0].users;

	if (Date.now() >= session.expiresAt.getTime()) {
		await db.delete(sessionsTable).where(eq(sessionsTable.id, session.id));
		return { session: null, user: null };
	}

	if (Date.now() >= session.expiresAt.getTime() - 1000 * 60 * 60 * 24 * 15) {
		const newExpiresAt = new Date(Date.now() + 1000 * 60 * 60 * 24 * 30);

		await db
			.update(sessionsTable)
			.set({
				expiresAt: newExpiresAt
			})
			.where(eq(sessionsTable.id, session.id));

		session.expiresAt = newExpiresAt;
	}

	return { session, user };
}

export async function invalidateSession(sessionId: string): Promise<void> {
	await db.delete(sessionsTable).where(eq(sessionsTable.id, sessionId));
}

// TODO: TTL
export async function invalidateUserSessions(userId: number): Promise<void> {
	await db.delete(sessionsTable).where(eq(sessionsTable.userId, userId));
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

export async function createSession(
	token: string,
	userId: number,
	accessToken: string,
	accessTokenExpiresAt: Date,
	refreshToken: string
): Promise<Session> {
	const sessionId = encodeHexLowerCase(sha256(new TextEncoder().encode(token)));
	const session: Session = {
		id: sessionId,
		userId,
		expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24 * 30),
		accessToken,
		accessTokenExpiresAt,
		refreshToken
	};

	const sessions = await db
		.insert(sessionsTable)
		.values({
			id: session.id,
			userId: session.userId,
			expiresAt: session.expiresAt,
			accessToken: session.accessToken,
			accessTokenExpiresAt: session.accessTokenExpiresAt,
			refreshToken: session.refreshToken
		})
		.returning();

	if (sessions.length === 0) {
		throw new Error('Failed to create session');
	}

	return sessions[0];
}

export async function updateSession(
	sessionId: string,
	accessToken: string,
	accessTokenExpiresAt: Date,
	refreshToken: string
): Promise<Session> {
	const sessions = await db
		.update(sessionsTable)
		.set({
			accessToken: accessToken,
			accessTokenExpiresAt: accessTokenExpiresAt,
			refreshToken: refreshToken
		})
		.where(eq(sessionsTable.id, sessionId))
		.returning();

	if (sessions.length === 0) {
		throw new Error('Failed to update session');
	}

	return sessions[0];
}

export type Session = typeof sessionsTable.$inferSelect;

type SessionValidationResult = { session: Session; user: User } | { session: null; user: null };
