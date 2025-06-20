import type { Bucket } from './data-table/types';

function fetchBuckets() {
    return [
        ...Array.from(
            { length: 20 },
            (_, i) => ({
                name: `bucket${(i + 1).toString().padStart(3, '0')}`,
                owner: `user${Math.floor(i / 3) + 1}`,
                usedCapacity: Math.floor(Math.random() * 500) * 1024 * 1024 * 1024, // Random GB between 0-500
                capacityLimit: 1000 * 1024 * 1024 * 1024, // 1000 GB
                objects: Math.floor(Math.random() * 500000),
                objectLimit: 1000000
            } as Bucket)
        )
    ]
}

export { fetchBuckets }