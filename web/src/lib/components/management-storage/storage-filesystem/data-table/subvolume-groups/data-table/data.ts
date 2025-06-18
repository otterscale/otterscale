import type { SubvolumeGroup } from './types';

const data: SubvolumeGroup[] = [
    ...Array.from(
        { length: 23 },
        (_, i) =>
            ({
                name: `subvolgrp_${(i + 1).toString().padStart(3, '0')}`,
                dataPool: `pool_${Math.floor(i / 5) + 1}`,
                usage: Math.floor(Math.random() * 1000),
                mode: i % 2 === 0 ? 'rw' : 'ro',
                createTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000))
            }) as SubvolumeGroup
    )
]

export {
    data
};
