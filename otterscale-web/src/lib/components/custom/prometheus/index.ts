import { PrometheusDriver, SampleValue } from 'prometheus-query';

// Define chart configuration type
export interface ChartConfig {
    [key: string]: {
        label: string;
        color: string;
    };
}

// Define data point type - supports dynamic key-value structure
export interface DataPoint {
    date: Date;
    [key: string]: any; // Allow dynamic key-value pairs, such as cpu0: 0.15, cpu1: 0.32, etc.
}

// Define Series type
export interface Series {
    key: string;
    label: string;
    color: string;
}

// One
const CHART_COLORS_ONE = 'var(--chart-5)';

// Two
const CHART_COLORS_TWO = [
    'var(--chart-4)',
    'var(--chart-7)',
];

// Two - Five
const CHART_COLORS_FEW = [
    'var(--chart-4)',
    'var(--chart-6)',
    'var(--chart-8)',
    'var(--chart-10)',
    'var(--chart-12)',
];

// Five up
const CHART_COLORS_MANY = [
    // 'var(--chart-1)',
    'var(--chart-2)',
    'var(--chart-3)',
    'var(--chart-4)',
    'var(--chart-5)',
    'var(--chart-6)',
    'var(--chart-7)',
    'var(--chart-8)',
    'var(--chart-9)',
    'var(--chart-10)',
    'var(--chart-11)',
    'var(--chart-12)',
];

/**
 * Auto-generate chartConfig based on data
 * @param data Data array
 * @returns Generated chartConfig object
 */
export function generateChartConfig(data: DataPoint[]): ChartConfig {
    if (!data || data.length === 0) return {};

    // Get all keys except 'date' from the first data point
    const samplePoint = data[0];
    const metricKeys = Object.keys(samplePoint).filter(key => key !== 'date');

    // Select color palette based on metric count
    let colorPalette: string[];

    if (metricKeys.length === 1) {
        colorPalette = [CHART_COLORS_ONE];
    } else if (metricKeys.length === 2) {
        colorPalette = CHART_COLORS_TWO;
    } else if (metricKeys.length <= 5) {
        colorPalette = CHART_COLORS_FEW;
    } else {
        colorPalette = CHART_COLORS_MANY;
    }

    // Generate chartConfig
    const chartConfig: ChartConfig = {};

    metricKeys.forEach((key, index) => {
        // Generate label (convert key to a more friendly display name)
        // const label = key.split(/(?=[A-Z])|_|-|\s+/)
        //     .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
        //     .join(' ');
        const label = key;
        // Assign color (cycle through selected color palette)
        const color = colorPalette[index % colorPalette.length];

        chartConfig[key] = {
            label,
            color
        };
    });

    return chartConfig;
}

/**
 * Generate Series array from chartConfig
 * @param config ChartConfig object
 * @returns Series array
 */
export function getSeries(config: ChartConfig): Series[] {
    return Object.keys(config).map((key) => ({
        key: key,
        label: config[key].label,
        color: config[key].color
    }));
}

/**
 * Convert old format data (time, value, metric) to new format (date, key1, key2, ...)
 * @param oldFormatData Old format data array
 * @returns New format DataPoint array
 */
export function convertToNewDataFormat(oldFormatData: Array<{ time: Date, value: number, metric: string }>): DataPoint[] {
    // Group by time
    const timeMap = new Map();

    oldFormatData.forEach(point => {
        const timeKey = point.time.getTime();
        if (!timeMap.has(timeKey)) {
            timeMap.set(timeKey, { date: point.time });
        }

        // Convert metric to valid key (remove spaces, convert to lowercase)
        // const metricKey = point.metric.replace(/\s+/g, '').toLowerCase();
        const metricKey = point.metric;
        timeMap.get(timeKey)[metricKey] = point.value;
    });

    // Convert to array and sort by time
    return Array.from(timeMap.values()).sort((a, b) => a.date.getTime() - b.date.getTime());
}

/**
 * Fetch and flatten Prometheus range query results
 * @param client PrometheusDriver instance
 * @param query Prometheus query string
 * @param timeStart Start time for the range query
 * @param timeEnd End time for the range query
 * @param step Step interval in seconds
 * @param metricName Optional custom metric name
 * @returns Flattened DataPoint array
 */
export async function fetchFlattenedRange(
    client: PrometheusDriver,
    query: string,
    timeStart?: Date,
    timeEnd?: Date,
    step: number = 15,
    metricName?: string
): Promise<DataPoint[]> {
    const start = timeStart || new Date(Date.now() - 60 * 60 * 1000); // 1 hour ago
    const end = timeEnd || new Date(); // Now

    try {
        const response = await client.rangeQuery(
            query,
            start.getTime(),
            end.getTime(),
            step
        );

        const oldFormatData = response.result.flatMap((series) => {
            const resolvedMetricName = resolveMetricName(series, metricName);
            return series.values.map((sampleValue: SampleValue) => ({
                time: sampleValue.time,
                value: sampleValue.value,
                metric: resolvedMetricName
            }));
        });

        return convertToNewDataFormat(oldFormatData);
    } catch (error) {
        console.error('Error fetching flattened range:', error);
        return [];
    }
}

/**
 * Fetch multiple queries and combine results into flattened range data
 * @param client PrometheusDriver instance
 * @param queries Object containing multiple named queries
 * @param timeStart Start time for the range query
 * @param timeEnd End time for the range query
 * @param step Step interval in seconds
 * @returns Combined DataPoint array with all metrics
 */
export async function fetchMultipleFlattenedRange(
    client: PrometheusDriver,
    queries: Record<string, string>,
    timeStart?: Date,
    timeEnd?: Date,
    step: number = 15
): Promise<DataPoint[]> {
    const start = timeStart || new Date(Date.now() - 60 * 60 * 1000); // 1 hour ago
    const end = timeEnd || new Date(); // Now

    try {
        const queryPromises = Object.entries(queries).map(([metricName, query]) =>
            executeQueryWithMetricName(client, query, start, end, step, metricName)
        );

        const allResults = await Promise.all(queryPromises);
        const combinedOldFormatData = allResults.flat();

        return convertToNewDataFormat(combinedOldFormatData);
    } catch (error) {
        console.error('Error fetching multiple flattened range:', error);
        return [];
    }
}

/**
 * Execute a single query and return formatted data
 * @param client PrometheusDriver instance
 * @param query Prometheus query string
 * @param start Start time
 * @param end End time
 * @param step Step interval
 * @param metricName Base metric name
 * @returns Formatted metric data array
 */
async function executeQueryWithMetricName(
    client: PrometheusDriver,
    query: string,
    start: Date,
    end: Date,
    step: number,
    metricName: string
): Promise<Array<{ time: Date; value: number; metric: string }>> {
    const response = await client.rangeQuery(
        query,
        start.getTime(),
        end.getTime(),
        step
    );

    return response.result.flatMap((series) => {
        const finalMetricName = buildUniqueMetricName(metricName, series, response.result.length);
        return series.values.map((sampleValue: SampleValue) => ({
            time: sampleValue.time,
            value: sampleValue.value,
            metric: finalMetricName
        }));
    });
}

/**
 * Resolve metric name from series or use provided name
 * @param series Prometheus series data
 * @param customName Optional custom metric name
 * @returns Resolved metric name
 */
function resolveMetricName(series: any, customName?: string): string {
    if (customName) {
        return customName;
    }

    const labels = series.metric.labels;
    if (labels && Object.keys(labels).length > 0) {
        const firstLabelValue = Object.values(labels)[0];
        return String(firstLabelValue);
    }

    return 'unknown';
}

/**
 * Build unique metric name for multi-series results
 * @param baseName Base metric name
 * @param series Prometheus series data
 * @param seriesCount Total number of series
 * @returns Unique metric name
 */
function buildUniqueMetricName(baseName: string, series: any, seriesCount: number): string {
    if (seriesCount === 1 || !series.metric.labels) {
        return baseName;
    }

    const labelEntries = Object.entries(series.metric.labels);
    if (labelEntries.length > 0) {
        const [, labelValue] = labelEntries[0];
        return `${baseName} ${labelValue}`;
    }

    return baseName;
}
