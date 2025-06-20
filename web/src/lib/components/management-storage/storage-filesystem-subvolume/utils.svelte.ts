import type { Subvolume } from './data-table/types';

function getData(): Subvolume[] {
    return [
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
}

const dataset = new Map([['key1', getData()], ['key2', getData()], ['key3', getData()]]);

const volumes = new Map([['', ['key1']], ['group', ['key2', 'key3']]]);

function fetchSubvolumeGroupList() {
    return ['', 'group']
}
function fetchSubvolumeListByGroup(group: string) {
    return volumes.get(group) ?? []
}
function fetchSubvolume(group: string, subvolume: string): Subvolume[] {
    return getData()
}

export {
    fetchSubvolumeGroupList, fetchSubvolumeListByGroup, fetchSubvolume
};
