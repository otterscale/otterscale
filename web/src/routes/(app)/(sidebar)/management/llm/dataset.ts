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

export interface LLMModel {
    name: string;
    version: string;
    parameters: string;
    metrics: {
        accuracy: number;
        speed: number;
    };
    architecture: string;
    usageStats: {
        requests: number;
        uptime: number;
    };
}

export const llmData: LLMModel[] = [
    {
        name: 'GPT-3',
        version: '1.0',
        parameters: '175B',
        metrics: { accuracy: 0.92, speed: 1.2 },
        architecture: 'Transformer',
        usageStats: { requests: 1000000, uptime: 99.9 }
    },
    {
        name: 'BERT',
        version: '2.0',
        parameters: '340M',
        metrics: { accuracy: 0.89, speed: 1.5 },
        architecture: 'Bidirectional Transformer',
        usageStats: { requests: 500000, uptime: 99.5 }
    },
    {
        name: 'LLaMA',
        version: '2.0',
        parameters: '65B',
        metrics: { accuracy: 0.91, speed: 1.3 },
        architecture: 'Transformer',
        usageStats: { requests: 800000, uptime: 99.7 }
    },
    {
        name: 'RoBERTa',
        version: '1.5',
        parameters: '355M',
        metrics: { accuracy: 0.9, speed: 1.4 },
        architecture: 'Bidirectional Transformer',
        usageStats: { requests: 600000, uptime: 99.6 }
    },
    {
        name: 'T5',
        version: '1.1',
        parameters: '11B',
        metrics: { accuracy: 0.88, speed: 1.6 },
        architecture: 'Encoder-Decoder',
        usageStats: { requests: 400000, uptime: 99.3 }
    },
    {
        name: 'BLOOM',
        version: '1.0',
        parameters: '176B',
        metrics: { accuracy: 0.91, speed: 1.1 },
        architecture: 'Transformer',
        usageStats: { requests: 300000, uptime: 99.4 }
    },
    {
        name: 'PaLM',
        version: '2.0',
        parameters: '540B',
        metrics: { accuracy: 0.93, speed: 1.0 },
        architecture: 'Transformer',
        usageStats: { requests: 900000, uptime: 99.8 }
    },
    {
        name: 'Claude',
        version: '2.0',
        parameters: '100B',
        metrics: { accuracy: 0.92, speed: 1.2 },
        architecture: 'Constitutional AI',
        usageStats: { requests: 700000, uptime: 99.6 }
    },
    {
        name: 'Falcon',
        version: '1.0',
        parameters: '40B',
        metrics: { accuracy: 0.89, speed: 1.4 },
        architecture: 'Transformer',
        usageStats: { requests: 200000, uptime: 99.2 }
    },
    {
        name: 'OPT',
        version: '1.3',
        parameters: '175B',
        metrics: { accuracy: 0.9, speed: 1.3 },
        architecture: 'Transformer',
        usageStats: { requests: 450000, uptime: 99.5 }
    }
];
