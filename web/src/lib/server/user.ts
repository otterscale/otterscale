import type { Row } from '@pilcrowjs/db-query';

import { db } from './db';

export interface User {
	id: number;
	sub: string;
	name: string;
	picture: string;
	email: string;
}

function mapRowToUser(row: Row, offset = 0): User {
	return {
		id: row.number(offset),
		sub: row.string(offset + 1),
		email: row.string(offset + 2),
		name: row.string(offset + 3),
		picture: row.string(offset + 4)
	};
}

export function createUser(sub: string, email: string, name: string, picture: string): User {
	const row = db.queryOne(
		'INSERT INTO user (sub, email, name, picture) VALUES (?, ?, ?, ?) RETURNING id, sub, email, name, picture',
		[sub, email, name, picture]
	);

	if (!row) {
		throw new Error('Failed to create user');
	}

	return mapRowToUser(row);
}

export function getUser(sub: string): User | null {
	const row = db.queryOne('SELECT id, sub, email, name, picture FROM user WHERE sub = ?', [sub]);
	return row ? mapRowToUser(row) : null;
}
