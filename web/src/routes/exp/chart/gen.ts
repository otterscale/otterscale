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

export const longData = [
    { year: 2019, basket: 1, fruit: 'apples', value: 3840 },
    { year: 2019, basket: 1, fruit: 'bananas', value: 1920 },
    { year: 2019, basket: 2, fruit: 'cherries', value: 960 },
    { year: 2019, basket: 2, fruit: 'grapes', value: 400 },

    { year: 2018, basket: 1, fruit: 'apples', value: 1600 },
    { year: 2018, basket: 1, fruit: 'bananas', value: 1440 },
    { year: 2018, basket: 2, fruit: 'cherries', value: 960 },
    { year: 2018, basket: 2, fruit: 'grapes', value: 400 },

    { year: 2017, basket: 1, fruit: 'apples', value: 820 },
    { year: 2017, basket: 1, fruit: 'bananas', value: 1000 },
    { year: 2017, basket: 2, fruit: 'cherries', value: 640 },
    { year: 2017, basket: 2, fruit: 'grapes', value: 400 },

    { year: 2016, basket: 1, fruit: 'apples', value: 820 },
    { year: 2016, basket: 1, fruit: 'bananas', value: 560 },
    { year: 2016, basket: 2, fruit: 'cherries', value: 720 },
    { year: 2016, basket: 2, fruit: 'grapes', value: 400 },
];