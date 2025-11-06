import type { User } from '$lib/jwt';

declare global {
	namespace App {
		interface Locals {
			user: User | null;
		}
	}
}

export {};
