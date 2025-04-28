import { addMinutes, startOfDay, startOfToday, subDays } from 'date-fns';

export function getRandomInteger(min: number, max: number, includeMax = true) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + (includeMax ? 1 : 0)) + min);
}

export function getRandomNumber(min: number, max: number) {
    return Math.random() * (max - min) + min;
}

export function createDateSeries<TKey extends string>(options: {
    count?: number;
    min: number;
    max: number;
    keys?: TKey[];
    value?: 'number' | 'integer';
}) {
    const now = startOfToday();

    const count = options.count ?? 10;
    const min = options.min;
    const max = options.max;
    const keys = options.keys ?? ['value'];

    return Array.from({ length: count }).map((_, i) => {
        return {
            date: subDays(now, count - i - 1),
            ...Object.fromEntries(
                keys.map((key) => {
                    return [
                        key,
                        options.value === 'integer' ? getRandomInteger(min, max) : getRandomNumber(min, max),
                    ];
                })
            ),
        } as { date: Date } & { [K in TKey]: number };
    });
}

export const healthRawData = [
    { type: 'error', namespace: 'prod/kubernetes-worker', total: (55 + 3 + 50), unhealth: (0), link: '/management/scope/2b28ecdc-51c7-4c50-8b39-ddeab98ddc14/facility/kubernetes-worker' },
    { type: 'error', namespace: 'dev/kubernetes-worker', total: (32 + 31 + 15), unhealth: (30), link: '/management/scope/db23d197-9178-4202-874c-b9374bc9987e/facility/kubernetes-worker' }
]

export const latencies = createDateSeries({
    count: 100,
    min: 0,
    max: 13,
    value: 'integer',
    keys: ['ceph', 'minio', 'x-inference']
});

export const storages = {
    capacity: 589.5,
    used: 7.2,
    unit: 'GB'
};

export const inputOutputs = createDateSeries({
    count: 100,
    min: 13,
    max: 23,
    value: 'number',
    keys: ['write', 'read']
});

export const currentUsage = {
    CPU: 13,
    GPU: 23,
    memory: {
        usage: 13,
        capacity: 23,
        unit: 'GB'
    }
}

export const usages = createDateSeries({
    count: 100,
    min: 0,
    max: 0.23,
    value: 'number',
    keys: ['CPU', 'GPU', 'Memory']
});
