import Redis from 'ioredis';

import { env } from '$env/dynamic/private';

export const redis = new Redis(env.REDIS_URL ?? '', { lazyConnect: true });
