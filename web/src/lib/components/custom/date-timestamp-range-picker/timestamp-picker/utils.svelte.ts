import { Time } from '@internationalized/date';

/**
 * regular expression to check for valid hour format (01-23)
 */
export function isValidHour(value: string) {
    return /^(0[0-9]|1[0-9]|2[0-3])$/.test(value);
}

/**
 * regular expression to check for valid minute format (00-59)
 */
export function isValidMinuteOrSecond(value: string) {
    return /^[0-5][0-9]$/.test(value);
}

type GetValidNumberConfig = { max: number; min?: number; loop?: boolean };

export function getValidNumber(
    value: string,
    { max, min = 0, loop = false }: GetValidNumberConfig
) {
    let numericValue = parseInt(value, 10);

    if (!isNaN(numericValue)) {
        if (!loop) {
            if (numericValue > max) numericValue = max;
            if (numericValue < min) numericValue = min;
        } else {
            if (numericValue > max) numericValue = min;
            if (numericValue < min) numericValue = max;
        }
        return numericValue.toString().padStart(2, "0");
    }

    return "00";
}

export function getValidHour(value: string) {
    if (isValidHour(value)) return value;
    return getValidNumber(value, { max: 23 });
}

export function getValidMinuteOrSecond(value: string) {
    if (isValidMinuteOrSecond(value)) return value;
    return getValidNumber(value, { max: 59 });
}

type GetValidArrowNumberConfig = {
    min: number;
    max: number;
    step: number;
};

export function getValidArrowNumber(
    value: string,
    { min, max, step }: GetValidArrowNumberConfig
) {
    let numericValue = parseInt(value, 10);
    if (!isNaN(numericValue)) {
        numericValue += step;
        return getValidNumber(String(numericValue), { min, max, loop: true });
    }
    return "00";
}

export function getValidArrowHour(value: string, step: number) {
    return getValidArrowNumber(value, { min: 0, max: 23, step });
}

export function getValidArrowMinuteOrSecond(value: string, step: number) {
    return getValidArrowNumber(value, { min: 0, max: 59, step });
}

export function setMinutes(time: Time, value: string) {
    const minutes = getValidMinuteOrSecond(value);
    return time.set({ minute: parseInt(minutes, 10) });
}

export function setSeconds(time: Time, value: string) {
    const seconds = getValidMinuteOrSecond(value);
    return time.set({ second: parseInt(seconds, 10) });
}

export function setHours(time: Time, value: string) {
    const hours = getValidHour(value);
    return time.set({ hour: parseInt(hours, 10) });
}

export type TimePickerType = "minutes" | "seconds" | "hours";

export function setDateByType(
    time: Time,
    value: string,
    type: TimePickerType,
) {
    switch (type) {
        case "minutes":
            return setMinutes(time, value);
        case "seconds":
            return setSeconds(time, value);
        case "hours":
            return setHours(time, value);
        default:
            return time;
    }
}

export function getDateByType(time: Time, type: TimePickerType) {
    switch (type) {
        case "minutes":
            return getValidMinuteOrSecond(String(time.minute));
        case "seconds":
            return getValidMinuteOrSecond(String(time.second));
        case "hours":
            return getValidHour(String(time.hour));
        default:
            return "00";
    }
}

export function getArrowByType(
    value: string,
    step: number,
    type: TimePickerType
) {
    switch (type) {
        case "minutes":
            return getValidArrowMinuteOrSecond(value, step);
        case "seconds":
            return getValidArrowMinuteOrSecond(value, step);
        case "hours":
            return getValidArrowHour(value, step);
        default:
            return "00";
    }
}
