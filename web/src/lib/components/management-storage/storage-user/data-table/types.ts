export type User = {
    userId: string;
    tenant: string;
    fullName: string;
    emailAddress: string;
    suspended: boolean;
    maximumBuckets: number;
    capacityLimit: number;
    objectLimit: number;
};
