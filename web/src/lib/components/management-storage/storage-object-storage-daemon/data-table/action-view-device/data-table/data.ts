import type { Device } from './types';

const data = [
    ...Array.from(
        { length: 30 },
        (_, i) => ({
            id: `osd.${i}`,
            name: `device-${Math.floor(i / 3)}`,
            state: i % 3 === 0 ? 'healthy' : i % 3 === 1 ? 'warning' : 'critical',
            lifeExpectancy: i % 4 === 0 ? 'good' : i % 4 === 1 ? 'fair' : i % 4 === 2 ? 'poor' : 'critical',
            daemons: [`daemon-${i}`, `daemon-${i + 1}`]
        }) as Device
    )
]

export {
    data
}