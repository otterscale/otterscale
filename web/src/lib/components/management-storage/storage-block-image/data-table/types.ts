export type BlockImage = {
    name: string;
    pool: string;
    namespace: string;
    size: number;
    usage: number;
    objects: number;
    objectSize: number;
    parent: string;
    mirroring: string;
    nextScheduledSnapshot: string;
};