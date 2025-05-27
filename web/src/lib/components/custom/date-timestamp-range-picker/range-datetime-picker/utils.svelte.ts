export function getStartPoint(precision: 'now' | 'minute' | 'hour' | 'day' = 'now') {
    if (precision === 'now') {
        const now = new Date();
        return now;
    }
    if (precision === 'minute') {
        const now = new Date();
        now.setSeconds(0, 0);
        return now;
    }
    if (precision === 'hour') {
        const now = new Date();
        now.setSeconds(0, 0);
        now.setMinutes(0);
        return now;
    }
    if (precision === 'day') {
        const now = new Date();
        now.setSeconds(0, 0);
        now.setMinutes(0);
        now.setHours(0);
        return now;
    }
    throw new Error(`Invalid precision: ${precision}`);
}

export function setValue(
    setStart: (date: Date) => void,
    setEnd: (date: Date) => void,
    begin: Date,
    interval: number
): void {
    setEnd(begin);
    setStart(new Date(begin.getTime() - interval));
}