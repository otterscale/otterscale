import type { Snapshot } from './data-table/types';

function fetchSnapshots() {
    return [
        ...Array.from(
            { length: 23 },
            (_, i) =>
                ({
                    name: `snapshot_${(i + 1).toString().padStart(3, '0')}`,
                    createTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000)),
                    pendingClones: i % 3 === 0 ? ['clone1', 'clone2']
                        : i % 3 === 1 ? ['clone3']
                            : [],
                }) as Snapshot
        )
    ]
}

export { fetchSnapshots }