import type { FlexibleIOTest } from './types';


const rwModeOptions: string[] = ['randread', 'read', 'randwrite', 'write'];


const data = [
    ...Array.from(
        { length: 12 },
        (_, i) =>
            ({
                name: `FIO_${i + 1}`,
                rwMode: rwModeOptions[Math.floor(Math.random() * rwModeOptions.length)],
                fileSize: `${Math.floor(Math.random() * 100)} ${i % 2 === 0 ? 'MB' : 'GB'}`,
                numberJobs: Math.floor(Math.random() * 100),
                blockSize: i % 2 === 0 ? '4 KB' : '1 MB',
                runtime: i % 2 === 0 ? 300 : 10,
                createTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000)),
                modifyTime: new Date(Date.now() - Math.floor(Math.random() * 10000000000))
            }) as FlexibleIOTest
    )
]

export {
    data
}