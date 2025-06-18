import type { ObjectStorageDaemon } from './types';

const data = [
    ...Array.from(
        { length: 30 },
        (_, i) => ({
            id: `osd.${i}`,
            host: `host-${Math.floor(i / 3)}`,
            status: i % 3 === 0
                ? ['up', 'in']
                : i % 3 === 1
                    ? ['up', 'out']
                    : ['down', 'in'],
            deviceClass: i % 2 === 0 ? 'ssd' : 'hdd',
            pgs: Math.floor(Math.random() * 100) + 50,
            size: Math.floor(Math.random() * 500 + 500), // 500-1000 GB
            flags: i % 4 === 0
                ? ['noup']
                : i % 4 === 1
                    ? ['nodown']
                    : i % 4 === 2
                        ? ['noin']
                        : [],
            usage: Math.floor(Math.random() * 100)
        }) as ObjectStorageDaemon
    )
]

export {
    data
}