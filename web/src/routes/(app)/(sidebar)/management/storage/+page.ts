import type { PageLoad } from './$types';
import type { FileSystem } from '$lib/components/management-storage/volume/data-table';

export const load: PageLoad = ({ params }) => {
    return {
        filesystem: [
            ...Array.from({ length: 23 }, (_, i) => ({
                id: `fs-${(i + 4).toString().padStart(3, '0')}`,
                name: `Sample_${i + 1}${i % 2 === 0 ? '.txt' : ''}`,
                type: i % 2 === 0 ? 'file' : 'directory',
                size: Math.floor(Math.random() * 10000000),
                enabled: Math.random() > 0.2,
                permission: i % 2 === 0 ? '-rw-r--r--' : 'drwxr-xr-x',
                createTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000)),
                modifyTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000))
            } as FileSystem))
        ]
    };;
};