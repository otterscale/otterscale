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

export function formatBigNumber(number: Number | BigInt) {
    return number.toLocaleString('en-US');
}