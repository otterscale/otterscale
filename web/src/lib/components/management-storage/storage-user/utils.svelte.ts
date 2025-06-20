import type { User } from './data-table/types';

function fetchUsers() {
    return [
        ...Array.from(
            { length: 20 },
            (_, i) => ({
                userId: `user${(i + 1).toString().padStart(3, '0')}`,
                tenant: `tenant${Math.floor(i / 5) + 1}`,
                fullName: `Test User ${i + 1}`,
                emailAddress: `user${i + 1}@example.com`,
                suspended: Math.random() > 0.8,
                maximumBuckets: Math.floor(Math.random() * 10) + 1,
                capacityLimit: Math.floor(Math.random() * 1000) * 1024 * 1024 * 1024, // Random GB between 0-1000
                objectLimit: Math.floor(Math.random() * 1000000)
            } as User)
        )
    ]
}

export { fetchUsers }