import { SampleValue } from 'prometheus-query';


export function metricColor(metric: number) {
    switch (true) {
        case metric > 62:
            return 'fill-red-800';
        case metric > 38:
            return 'fill-yellow-500';
        default:
            return 'fill-green-800';
    }
}
export function metricBackgroundColor(metric: number) {
    switch (true) {
        case metric > 62:
            return 'fill-red-50';
        case metric > 38:
            return 'fill-yellow-50';
        default:
            return 'fill-green-50';
    }
}

export function integrateSerieses(
    serieses: Map<string, SampleValue[] | undefined>,
) {
    if (!serieses) {
        return [];
    }

    const nonFalseSerieses = Object.fromEntries(
        Array.from(serieses.entries()).filter(
            ([_, v]) => v !== undefined && v.length > 0,
        ),
    );
    if (Object.entries(nonFalseSerieses).length === 0) {
        return [];
    }

    const lengths = new Set(Object.values(nonFalseSerieses).map(a => a?.length));
    if (lengths.size !== 1) {
        throw new Error('All series must have the same length');
    }

    const keys = new Set(Object.keys(nonFalseSerieses));
    if (keys.size === 0) {
        return [];
    }

    Object.values(nonFalseSerieses).forEach(series => {
        if (series) {
            series.sort((p, n) => p.time.getTime() - n.time.getTime());
        }
    });

    const anyTrueSeries = Array.from(Object.values(nonFalseSerieses))[0];
    if (!anyTrueSeries) {
        return [];
    }
    return anyTrueSeries.map((samplePoint, index) => {
        return Object.fromEntries([
            ["time", samplePoint.time],
            ...Object.entries(nonFalseSerieses)
                .map(([label, sampleSpace]) => {
                    if (sampleSpace) {
                        return [label, sampleSpace[index].value];
                    }
                })
                .filter((item) => item !== undefined),
        ]);
    });
}