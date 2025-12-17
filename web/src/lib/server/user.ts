import { eq } from 'drizzle-orm';

import { db } from './db';
import { usersTable } from './db/schema';

export type User = typeof usersTable.$inferSelect;

export async function createUser(
	sub: string,
	username: string,
	email: string,
	name: string,
	picture: string
): Promise<User> {
	const users = await db
		.insert(usersTable)
		.values({
			sub,
			username,
			email,
			name,
			picture
		})
		.returning();

	if (users.length === 0) {
		throw new Error('Failed to create user');
	}

	return users[0];
}

export async function getUser(sub: string): Promise<User | null> {
	const users = await db.select().from(usersTable).where(eq(usersTable.sub, sub));

	return users.length > 0 ? users[0] : null;
}
