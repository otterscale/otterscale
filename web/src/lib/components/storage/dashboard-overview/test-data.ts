import * as Chart from '$lib/components/ui/chart/index.js';
import ActivityIcon from '@lucide/svelte/icons/activity';


// Static test data moved out from index.svelte
export const chartData1 = [{ browser: 'safari', visitors: 1260, color: 'var(--color-safari)' }];
export const chartConfig1 = {
    visitors: { label: 'Visitors' },
    safari: { label: 'Safari', color: 'var(--chart-2)' }
} satisfies Chart.ChartConfig;


export const chartData2 = [
  { date: new Date('2024-01-01'), desktop: 10 },
  { date: new Date('2024-02-01'), desktop: 12 },
  { date: new Date('2024-03-01'), desktop: 24 },
  { date: new Date('2024-04-01'), desktop: 28 },
  { date: new Date('2024-05-01'), desktop: 36 },
  { date: new Date('2024-06-01'), desktop: 48 }
];
export const chartConfig2 = {
    desktop: { label: 'Desktop', color: 'var(--chart-1)', icon: ActivityIcon }
} satisfies Chart.ChartConfig;


export const chartData3 = [
  { date: new Date('2024-04-01'), desktop: 222, mobile: 150 },
  { date: new Date('2024-04-02'), desktop: 97, mobile: 180 },
  { date: new Date('2024-04-03'), desktop: 167, mobile: 120 },
  { date: new Date('2024-04-04'), desktop: 242, mobile: 260 },
  { date: new Date('2024-04-05'), desktop: 373, mobile: 290 },
  { date: new Date('2024-04-06'), desktop: 301, mobile: 340 },
  { date: new Date('2024-04-07'), desktop: 245, mobile: 180 },
  { date: new Date('2024-04-08'), desktop: 409, mobile: 320 },
  { date: new Date('2024-04-09'), desktop: 59, mobile: 110 },
  { date: new Date('2024-04-10'), desktop: 261, mobile: 190 },
  { date: new Date('2024-04-11'), desktop: 327, mobile: 350 },
  { date: new Date('2024-04-12'), desktop: 292, mobile: 210 },
  { date: new Date('2024-04-13'), desktop: 342, mobile: 380 },
  { date: new Date('2024-04-14'), desktop: 137, mobile: 220 },
  { date: new Date('2024-04-15'), desktop: 120, mobile: 170 },
  { date: new Date('2024-04-16'), desktop: 138, mobile: 190 },
  { date: new Date('2024-04-17'), desktop: 446, mobile: 360 },
  { date: new Date('2024-04-18'), desktop: 364, mobile: 410 },
  { date: new Date('2024-04-19'), desktop: 243, mobile: 180 },
  { date: new Date('2024-04-20'), desktop: 89, mobile: 150 },
  { date: new Date('2024-04-21'), desktop: 137, mobile: 200 },
  { date: new Date('2024-04-22'), desktop: 224, mobile: 170 },
  { date: new Date('2024-04-23'), desktop: 138, mobile: 230 },
  { date: new Date('2024-04-24'), desktop: 387, mobile: 290 },
  { date: new Date('2024-04-25'), desktop: 215, mobile: 250 },
  { date: new Date('2024-04-26'), desktop: 75, mobile: 130 },
  { date: new Date('2024-04-27'), desktop: 383, mobile: 420 },
  { date: new Date('2024-04-28'), desktop: 122, mobile: 180 },
  { date: new Date('2024-04-29'), desktop: 315, mobile: 240 },
  { date: new Date('2024-04-30'), desktop: 454, mobile: 380 },
  { date: new Date('2024-05-01'), desktop: 165, mobile: 220 },
  { date: new Date('2024-05-02'), desktop: 293, mobile: 310 },
  { date: new Date('2024-05-03'), desktop: 247, mobile: 190 },
  { date: new Date('2024-05-04'), desktop: 385, mobile: 420 },
  { date: new Date('2024-05-05'), desktop: 481, mobile: 390 },
  { date: new Date('2024-05-06'), desktop: 498, mobile: 520 },
  { date: new Date('2024-05-07'), desktop: 388, mobile: 300 },
  { date: new Date('2024-05-08'), desktop: 149, mobile: 210 },
  { date: new Date('2024-05-09'), desktop: 227, mobile: 180 },
  { date: new Date('2024-05-10'), desktop: 293, mobile: 330 },
  { date: new Date('2024-05-11'), desktop: 335, mobile: 270 },
  { date: new Date('2024-05-12'), desktop: 197, mobile: 240 },
  { date: new Date('2024-05-13'), desktop: 197, mobile: 160 },
  { date: new Date('2024-05-14'), desktop: 448, mobile: 490 },
  { date: new Date('2024-05-15'), desktop: 473, mobile: 380 },
  { date: new Date('2024-05-16'), desktop: 338, mobile: 400 },
  { date: new Date('2024-05-17'), desktop: 499, mobile: 420 },
  { date: new Date('2024-05-18'), desktop: 315, mobile: 350 },
  { date: new Date('2024-05-19'), desktop: 235, mobile: 180 },
  { date: new Date('2024-05-20'), desktop: 177, mobile: 230 },
  { date: new Date('2024-05-21'), desktop: 82, mobile: 140 },
  { date: new Date('2024-05-22'), desktop: 81, mobile: 120 },
  { date: new Date('2024-05-23'), desktop: 252, mobile: 290 },
  { date: new Date('2024-05-24'), desktop: 294, mobile: 220 },
  { date: new Date('2024-05-25'), desktop: 201, mobile: 250 },
  { date: new Date('2024-05-26'), desktop: 213, mobile: 170 },
  { date: new Date('2024-05-27'), desktop: 420, mobile: 460 },
  { date: new Date('2024-05-28'), desktop: 233, mobile: 190 },
  { date: new Date('2024-05-29'), desktop: 78, mobile: 130 },
  { date: new Date('2024-05-30'), desktop: 340, mobile: 280 },
  { date: new Date('2024-05-31'), desktop: 178, mobile: 230 },
  { date: new Date('2024-06-01'), desktop: 178, mobile: 200 },
  { date: new Date('2024-06-02'), desktop: 470, mobile: 410 },
  { date: new Date('2024-06-03'), desktop: 103, mobile: 160 },
  { date: new Date('2024-06-04'), desktop: 439, mobile: 380 },
  { date: new Date('2024-06-05'), desktop: 88, mobile: 140 },
  { date: new Date('2024-06-06'), desktop: 294, mobile: 250 },
  { date: new Date('2024-06-07'), desktop: 323, mobile: 370 },
  { date: new Date('2024-06-08'), desktop: 385, mobile: 320 },
  { date: new Date('2024-06-09'), desktop: 438, mobile: 480 },
  { date: new Date('2024-06-10'), desktop: 155, mobile: 200 },
  { date: new Date('2024-06-11'), desktop: 92, mobile: 150 },
  { date: new Date('2024-06-12'), desktop: 492, mobile: 420 },
  { date: new Date('2024-06-13'), desktop: 81, mobile: 130 },
  { date: new Date('2024-06-14'), desktop: 426, mobile: 380 },
  { date: new Date('2024-06-15'), desktop: 307, mobile: 350 },
  { date: new Date('2024-06-16'), desktop: 371, mobile: 310 },
  { date: new Date('2024-06-17'), desktop: 475, mobile: 520 },
  { date: new Date('2024-06-18'), desktop: 107, mobile: 170 },
  { date: new Date('2024-06-19'), desktop: 341, mobile: 290 },
  { date: new Date('2024-06-20'), desktop: 408, mobile: 450 },
  { date: new Date('2024-06-21'), desktop: 169, mobile: 210 },
  { date: new Date('2024-06-22'), desktop: 317, mobile: 270 },
  { date: new Date('2024-06-23'), desktop: 480, mobile: 530 },
  { date: new Date('2024-06-24'), desktop: 132, mobile: 180 },
  { date: new Date('2024-06-25'), desktop: 141, mobile: 190 },
  { date: new Date('2024-06-26'), desktop: 434, mobile: 380 },
  { date: new Date('2024-06-27'), desktop: 448, mobile: 490 },
  { date: new Date('2024-06-28'), desktop: 149, mobile: 200 },
  { date: new Date('2024-06-29'), desktop: 103, mobile: 160 },
  { date: new Date('2024-06-30'), desktop: 446, mobile: 400 }
];
export const chartConfig3 = {
    views: { label: 'Page Views', color: '' },
    desktop: { label: 'Desktop', color: 'var(--chart-1)' },
    mobile: { label: 'Mobile', color: 'var(--chart-3)' }
} satisfies Chart.ChartConfig;

export const chartData4 = [
  // reuse same shape as chartData3 — duplicated values intentionally preserved
  ...chartData3
];
export const chartConfig4 = {
  views: { label: 'Page Views', color: '' },
  desktop: { label: 'Desktop', color: 'var(--chart-2)' },
  mobile: { label: 'Mobile', color: 'var(--chart-4)' }
};

export const chartData5 = [
  { date: new Date('2024-01-01'), desktop: 186 },
  { date: new Date('2024-02-01'), desktop: 305 },
  { date: new Date('2024-03-01'), desktop: 237 },
  { date: new Date('2024-04-01'), desktop: 73 },
  { date: new Date('2024-05-01'), desktop: 209 },
  { date: new Date('2024-06-01'), desktop: 214 }
];
export const chartConfig5 = {
  desktop: { label: 'Desktop', color: 'var(--chart-5)' }
};

export const chartData6 = [
  { browser: 'PCIe', visitors: 275, color: 'var(--color-chrome)' },
  { browser: 'SSD', visitors: 200, color: 'var(--color-safari)' }
];
export const chartConfig6 = {
  visitors: { label: 'Visitors' },
  chrome: { label: 'Chrome', color: 'var(--chart-1)' },
  safari: { label: 'Safari', color: 'var(--chart-2)' },
  firefox: { label: 'Firefox', color: 'var(--chart-3)' },
  edge: { label: 'Edge', color: 'var(--chart-4)' },
  other: { label: 'Other', color: 'var(--chart-5)' }
};

export const totalVisitors6 = chartData6.reduce((acc, curr) => acc + curr.visitors, 0);
