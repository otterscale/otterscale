import dayjs from 'dayjs';


const DIVISIONS = [
    { amount: 60, name: 'seconds' },
    { amount: 60, name: 'minutes' },
    { amount: 24, name: 'hours' },
    { amount: 7, name: 'days' },
    { amount: 4.34524, name: 'weeks' },
    { amount: 12, name: 'months' },
    { amount: Number.POSITIVE_INFINITY, name: 'years' }
] as const;

const formatter = new Intl.RelativeTimeFormat(undefined, {
    numeric: 'auto'
});

export function formatTimeAgo(date: Date) {
    let duration = (new Date(date).getTime() - new Date().getTime()) / 1000;

    for (let i = 0; i <= DIVISIONS.length; i++) {
        const division = DIVISIONS[i];
        if (Math.abs(duration) < division.amount) {
            return formatter.format(Math.round(duration), division.name);
        }
        duration /= division.amount;
    }
}

export function formatDuration(duration: number): { value: number, unit: string } {
    if (duration === 0) return { value: 0, unit: "Second" };

    const years = (duration / (365 * 24 * 3600));
    if (years >= 1) return { value: years, unit: "Year" };

    const weeks = ((duration % (365 * 24 * 3600)) / (7 * 24 * 3600));
    if (weeks >= 1) return { value: weeks, unit: "Week" };

    const days = ((duration % (7 * 24 * 3600)) / (24 * 3600));
    if (days >= 1) return { value: days, unit: "Day" };

    const hours = ((duration % (24 * 3600)) / 3600);
    if (hours >= 1) return { value: hours, unit: "Hour" };

    const minutes = ((duration % 3600) / 60);
    if (minutes >= 1) return { value: minutes, unit: "Minute" };

    const seconds = (duration % 60);
    return { value: seconds, unit: "Second" };
}

export function formatCapacity(capacity: number | bigint): { value: string, unit: string } {
    const MB = Number(capacity)
    const GB = MB / 1024;
    const TB = GB / 1024;

    if (TB >= 1) {
        return { value: `${Math.round(TB * 100) / 100}`, unit: "TB" };
    } else if (GB >= 1) {
        return { value: `${Math.round(GB * 100) / 100}`, unit: "GB" };
    } else {
        return { value: `${Math.round(MB * 100) / 100}`, unit: "MB" };
    }
}

export function formatNetworkIO(bytes: number | bigint): { value: number, unit: string } {
    const B = Number(bytes);
    const KB = B / 1024;
    const MB = KB / 1024;
    const GB = MB / 1024;
    const TB = GB / 1024;


    if (TB >= 1) {
        return { value: Math.round(TB * 100) / 100, unit: "TB/s" };
    } else if (GB >= 1) {
        return { value: Math.round(GB * 100) / 100, unit: "GB/s" };
    } else if (MB >= 1) {
        return { value: Math.round(MB * 100) / 100, unit: "MB/s" };
    } else if (KB >= 1) {
        return { value: Math.round(KB * 100) / 100, unit: "KB/s" };
    } else {
        return { value: Math.round(B * 100) / 100, unit: "B/s" };
    }
}

export function formatBigNumber(number: Number | BigInt) {
    return number.toLocaleString('en-US');
}

export const formatTime = (v: Date | number): string => {
    return dayjs(v).format('HH:mm');
};