import type { Volume } from './data-table/types';

function fetchVolumes() {
    return [
        ...Array.from(
            { length: 23 },
            (_, i) =>
                ({
                    id: `fs-${(i + 4).toString().padStart(3, '0')}`,
                    name: `Sample_${i + 1}${i % 2 === 0 ? '.txt' : ''}`,
                    type: i % 2 === 0 ? 'file' : 'directory',
                    size: Math.floor(Math.random() * 10000000),
                    enabled: Math.random() > 0.2,
                    permission: i % 2 === 0 ? '-rw-r--r--' : 'drwxr-xr-x',
                    createTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000)),
                    modifyTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000))
                }) as Volume
        )
    ]
}

export { fetchVolumes }