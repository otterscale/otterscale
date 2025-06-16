import type { BlockImageSnapshot } from './types';

const data = Array.from(
    { length: 23 },
    (_, i): BlockImageSnapshot => ({
        name: `Snapshot ${i + 1}`,
        size: Math.round(Math.random() * 1024),
        used: Math.round(Math.random() * 512),
        state: Math.random() > 0.5 ? 'available' : 'pending',
        createTime: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString()
    })
);

export {
    data
}