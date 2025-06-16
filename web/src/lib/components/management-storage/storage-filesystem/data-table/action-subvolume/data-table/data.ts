import type { Subvolume } from './types';

const data: Subvolume[] = [
    ...Array.from(
        { length: 23 },
        (_, i) =>
            ({
                name: `subvol_${(i + 1).toString().padStart(3, '0')}`,
                dataPool: `pool_${Math.floor(i / 5) + 1}`,
                usage: Math.floor(Math.random() * 1000),
                path: `/storage/subvolumes/subvol_${(i + 1).toString().padStart(3, '0')}`,
                mode: i % 2 === 0 ? 'rw' : 'ro',
                createTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000))
            }) as Subvolume
    )
]

export {
    data
}