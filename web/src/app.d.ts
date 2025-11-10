import type { Session, User } from '$lib/server';

declare global {
	namespace App {
		interface Locals {
			user: User | null;
			session: Session | null;
		}
	}
}

export {};
