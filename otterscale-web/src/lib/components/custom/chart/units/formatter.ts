// Constants for time range formatting
const HOURLY_TICKS = 24;
const DAILY_TICKS = 24;
const WEEKLY_TICKS = 7;
const MONTHLY_TICKS = 30;

// Time range type definition
export type TimeRange = '1h' | '1d' | '7d' | '30d';

/**
 * Format time range hours into appropriate string representation
 */
export const formatTimeRange = (hours: number): TimeRange => {
    if (hours === 1) {
        return '1h';
    } else if (hours === 24) {
        return '1d';
    } else if (hours === 168) { // 7 * 24
        return '7d';
    } else if (hours === 720) { // 30 * 24
        return '30d';
    } else if (hours < 24) {
        return '1h'; // Default to 1h for sub-daily ranges
    } else if (hours < 168) {
        return '1d'; // Default to 1d for sub-weekly ranges
    } else if (hours < 720) {
        return '7d'; // Default to 7d for sub-monthly ranges
    } else {
        return '30d'; // Default to 30d for longer ranges
    }
};

/**
 * Format date for x-axis display based on time range
 */
export const formatXAxisDate = (date: any, timeRange?: TimeRange): string => {
    switch (timeRange) {
        case '1h':
            return date.toLocaleTimeString('en-US', {
                hour: 'numeric',
                minute: '2-digit'
            });
        case '1d':
            return date.toLocaleTimeString('en-US', {
                hour: 'numeric'
            });
        case '7d':
            return date.toLocaleTimeString('en-US', {
                day: 'numeric',
                hour: 'numeric'
            });
        case '30d':
            return date.toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric'
            });
        default:
            return date.toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric'
            });
    }
};

/**
 * Format date for tooltip display based on time range
 */
export const formatTooltipDate = (date: any, timeRange?: TimeRange): string => {
    switch (timeRange) {
        case '1h':
            return date.toLocaleTimeString('en-US', {
                hour: 'numeric',
                minute: '2-digit'
            });
        case '1d':
            return date.toLocaleTimeString('en-US', {
                hour: 'numeric'
            });
        case '7d':
            return date.toLocaleTimeString('en-US', {
                day: 'numeric',
                hour: 'numeric'
            });
        case '30d':
            return date.toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric'
            });
        default:
            return date.toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric'
            });
    }
};

/**
 * Get number of ticks for x-axis based on time range
 */
export const getXAxisTicks = (timeRange?: TimeRange): number | undefined => {
    switch (timeRange) {
        case '1h':
            return HOURLY_TICKS;
        case '1d':
            return DAILY_TICKS;
        case '7d':
            return WEEKLY_TICKS;
        case '30d':
            return MONTHLY_TICKS;
        default:
            return undefined;
    }
};
