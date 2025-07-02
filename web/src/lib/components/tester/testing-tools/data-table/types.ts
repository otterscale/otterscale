export type FlexibleIOTest = {
    name: string;
    rwMode: string,             // rw: randread, read, randwrite, write
    fileSize: string,           // filesize
    numberJobs: number,         // numjobs
    blockSize: string,          // bs
    runtime: number,            // runtime
    createTime: Date;
    modifyTime: Date,
};
