import type { Session } from '$lib/server/session';

declare global {
	namespace App {
		interface Locals {
			session: Session | null;
		}
	}
}

export {};
