import type { BlockImage } from './types';

const namespaces = ['default', 'storage', 'system', 'backup', 'test'];

const data: BlockImage[] = [
    ...Array.from(
        { length: 15 },
        (_, i): BlockImage => ({
            name: `Block Image ${i + 1}`,
            pool: `pool-${(Math.floor(Math.random() * 5) + 1).toString().padStart(3, '0')}`,
            namespace: namespaces[Math.floor(Math.random() * namespaces.length)],
            size: Math.round(Math.random() * 1024),
            usage: Math.round(Math.random() * 100),
            objects: Math.round(Math.random() * 128),
            objectSize: Math.round(Math.random() * 1024),
            parent: '-',
            mirroring: Math.random() > 0.5 ? 'enabled' : 'disabled',
            nextScheduledSnapshot: ''
        })
    )
];

export {
    data
}