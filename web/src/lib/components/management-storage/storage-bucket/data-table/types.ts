export type Bucket = {
    name: string;
    owner: string;
    usedCapacity: number;
    capacityLimit: number;
    objects: number;
    objectLimit: number;
};
