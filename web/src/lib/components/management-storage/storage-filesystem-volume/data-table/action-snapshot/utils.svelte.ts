import type { Snapshot } from './data-table/types';

function fetchSnapshots() {
    return [
        ...Array.from(
            { length: 23 },
            (_, i) =>
                ({
                    name: `snapshot_${(i + 1).toString().padStart(3, '0')}`,
                    path: i % 3 === 0 ? 'path1'
                        : i % 3 === 1 ? 'path2'
                            : 'path3',
                    createTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000)),
                }) as Snapshot
        )
    ]
}

export { fetchSnapshots }