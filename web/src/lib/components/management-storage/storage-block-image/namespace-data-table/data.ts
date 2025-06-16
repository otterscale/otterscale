import type { Namespace } from './types';

const namespaces = ['default', 'storage', 'system', 'backup', 'test'];

const data: Namespace[] = [
    ...Array.from(
        { length: 15 },
        (_, i): Namespace => ({
            pool: `pool-${(Math.floor(Math.random() * 5) + 1).toString().padStart(3, '0')}`,
            namespace: namespaces[Math.floor(Math.random() * namespaces.length)],
            totalImages: Math.round(Math.random() * 100),
        })
    )
];

export {
    data
}