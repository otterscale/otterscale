import type { TableRow } from './type';

async function getDataset(): Promise<TableRow[]> {
	return [
		...Array.from({ length: 10 }, (_, i) => ({
			id: i + 1,
			name: `User ${i + 1}`,
			email: `user${i + 1}@example.com`,
			role: (i % 3 === 0 ? 'Admin' : i % 2 === 0 ? 'Manager' : 'User') as
				| 'Admin'
				| 'Manager'
				| 'User',
			status: (i % 4 === 0 ? 'inactive' : 'active') as 'inactive' | 'active',
			createdAt: new Date(2024, 0, i + 1).toISOString().split('T')[0],
			isVerified: i % 2 === 0,
			loginCount: Math.floor(Math.random() * 100),
			rating: (Math.random() * 5).toFixed(1),
			lastLoginDays: Math.floor(Math.random() * 30)
		}))
	];
}

export { getDataset };
