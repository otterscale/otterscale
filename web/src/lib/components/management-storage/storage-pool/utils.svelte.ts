import type { Pool } from './data-table/types';

function fetchPools() {
    return [
        ...Array.from(
            { length: 23 },
            (_, i) =>
                ({
                    id: `pool-${(i + 1).toString().padStart(3, '0')}`,
                    name: `Storage Pool ${i + 1}`,
                    dataProtection: i % 3 === 0 ? 'Replicated' : i % 3 === 1 ? 'Erasure Coded' : 'None',
                    applications: i % 2 === 0 ? 'Block, File' : 'Object',
                    PGStatus: i % 4 === 0 ? 'Active+Clean' : i % 4 === 1 ? 'Active+Degraded' : i % 4 === 2 ? 'Active+Remapped' : 'Active',
                    usage: Math.floor(Math.random() * 100)
                }) as Pool
        )
    ]
}

export { fetchPools };

