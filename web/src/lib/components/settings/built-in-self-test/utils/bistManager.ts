import type { Table } from '@tanstack/table-core';

import { hashCode } from './hashGroupID';

import type { FIO_Input, TestResult, Warp_Input } from '$lib/api/bist/v1/bist_pb';
import { Warp_Input_Operation } from '$lib/api/bist/v1/bist_pb';
import { formatCapacity, formatSecond } from '$lib/formatter';

interface FioDataPoint {
	name: string;
	ioBytes: number;
	bandwidthBytes: number;
	ioPerSecond: number;
	totalIos: number;
	latency: number;
	completedAt: Date;
}

interface FioOutputGroup {
	key: string;
	data: FioDataPoint[];
	color: string;
}

interface WarpDataPoint {
	name: string;
	totalBytes: number;
	totalObjects: number;
	totalOperations: number;
	bytesFastest: number;
	bytesMedian: number;
	bytesSlowest: number;
	objectsFastest: number;
	objectsMedian: number;
	objectsSlowest: number;
	completedAt: Date;
}

interface WarpOutputGroup {
	key: string;
	data: WarpDataPoint[];
	color: string;
}

class BistDashboardManager<TData = TestResult> {
	table: Table<TData>;

	constructor(table: Table<TData>) {
		this.table = table;
	}

	get filteredData() {
		return this.table.getFilteredRowModel().rows.map((row) => row.original);
	}

	private generateFioGroupName(input: FIO_Input): string {
		const runTime = formatSecond(Number(input.runTimeSeconds));
		const blockSize = formatCapacity(Number(input.blockSizeBytes));
		const fileSize = formatCapacity(Number(input.fileSizeBytes));
		return `${Number(input.jobCount)}-${runTime.value}${runTime.unit}-${blockSize.value}${blockSize.unit}-${fileSize.value}${fileSize.unit}-${Number(input.ioDepth)}`;
	}

	private generateWarpGroupName(input: Warp_Input): string {
		const duration = formatSecond(Number(input.durationSeconds));
		const objectSize = formatCapacity(Number(input.objectSizeBytes));
		return `${Warp_Input_Operation[input.operation]}-${duration.value}${duration.unit}-${objectSize.value}${objectSize.unit}-${input.objectCount}`;
	}

	// Improved: consistent color mapping per unique groupName, avoids collisions as much as possible
	private static colorMap: Map<string, string> = new Map();

	private generateColor(groupName: string): string {
		const CHART_COLORS_MANY = [
			'var(--chart-1)',
			'var(--chart-2)',
			'var(--chart-3)',
			'var(--chart-4)',
			'var(--chart-5)',
		];

		// Use a static map to ensure each groupName always gets the same color
		if (BistDashboardManager.colorMap.has(groupName)) {
			return BistDashboardManager.colorMap.get(groupName)!;
		}
		// Find the next available color in the palette
		for (const color of CHART_COLORS_MANY) {
			if (![...BistDashboardManager.colorMap.values()].includes(color)) {
				BistDashboardManager.colorMap.set(groupName, color);
				return color;
			}
		}
		// If all colors are used, fallback to hash-based assignment (may cause collisions)
		const idx = Math.abs(hashCode(groupName)) % CHART_COLORS_MANY.length;
		const color = CHART_COLORS_MANY[idx];
		BistDashboardManager.colorMap.set(groupName, color);
		return color;
	}

	getFioOutputs(): {
		read: Record<string, FioOutputGroup>;
		write: Record<string, FioOutputGroup>;
		trim: Record<string, FioOutputGroup>;
	} {
		const outputMap: Record<string, Record<string, FioOutputGroup>> = {
			read: {},
			write: {},
			trim: {},
		};

		this.filteredData.forEach((datum) => {
			const testResult = datum as TestResult;

			if (testResult.kind.case === 'fio' && testResult.kind.value.output && testResult.kind.value.input) {
				const groupName = this.generateFioGroupName(testResult.kind.value.input);

				// Process read output
				if (testResult.kind.value.output.read && testResult.kind.value.output.read.latency) {
					const dataPoint: FioDataPoint = {
						name: testResult.name,
						ioBytes: Number(testResult.kind.value.output.read.ioBytes),
						bandwidthBytes: Number(testResult.kind.value.output.read.bandwidthBytes),
						ioPerSecond: testResult.kind.value.output.read.ioPerSecond,
						totalIos: Number(testResult.kind.value.output.read.totalIos),
						latency: testResult.kind.value.output.read.latency?.meanNanoseconds,
						completedAt: testResult.completedAt
							? new Date(
									Number(testResult.completedAt.seconds) * 1000 +
										Number(testResult.completedAt.nanos) / 1000000,
								)
							: new Date(),
					};

					if (!outputMap['read'][groupName]) {
						outputMap['read'][groupName] = {
							key: groupName,
							data: [],
							color: this.generateColor(groupName),
						};
					}

					outputMap['read'][groupName].data.push(dataPoint);
				}

				// Process write output
				if (testResult.kind.value.output.write && testResult.kind.value.output.write.latency) {
					const dataPoint: FioDataPoint = {
						name: testResult.name,
						ioBytes: Number(testResult.kind.value.output.write.ioBytes),
						bandwidthBytes: Number(testResult.kind.value.output.write.bandwidthBytes),
						ioPerSecond: testResult.kind.value.output.write.ioPerSecond,
						totalIos: Number(testResult.kind.value.output.write.totalIos),
						latency: testResult.kind.value.output.write.latency?.meanNanoseconds,
						completedAt: testResult.completedAt
							? new Date(
									Number(testResult.completedAt.seconds) * 1000 +
										Number(testResult.completedAt.nanos) / 1000000,
								)
							: new Date(),
					};

					if (!outputMap['write'][groupName]) {
						outputMap['write'][groupName] = {
							key: groupName,
							data: [],
							color: this.generateColor(groupName),
						};
					}

					outputMap['write'][groupName].data.push(dataPoint);
				}

				// Process trim output
				if (testResult.kind.value.output.trim && testResult.kind.value.output.trim.latency) {
					const dataPoint: FioDataPoint = {
						name: testResult.name,
						ioBytes: Number(testResult.kind.value.output.trim.ioBytes),
						bandwidthBytes: Number(testResult.kind.value.output.trim.bandwidthBytes),
						ioPerSecond: testResult.kind.value.output.trim.ioPerSecond,
						totalIos: Number(testResult.kind.value.output.trim.totalIos),
						latency: testResult.kind.value.output.trim.latency?.meanNanoseconds,
						completedAt: testResult.completedAt
							? new Date(
									Number(testResult.completedAt.seconds) * 1000 +
										Number(testResult.completedAt.nanos) / 1000000,
								)
							: new Date(),
					};

					if (!outputMap['trim'][groupName]) {
						outputMap['trim'][groupName] = {
							key: groupName,
							data: [],
							color: this.generateColor(groupName),
						};
					}

					outputMap['trim'][groupName].data.push(dataPoint);
				}
			}
		});
		return {
			read: outputMap['read'],
			write: outputMap['write'],
			trim: outputMap['trim'],
		};
	}

	getWarpOutputs(): {
		get: Record<string, WarpOutputGroup>;
		put: Record<string, WarpOutputGroup>;
		delete: Record<string, WarpOutputGroup>;
	} {
		const outputMap: Record<string, Record<string, WarpOutputGroup>> = {
			get: {},
			put: {},
			delete: {},
		};

		this.filteredData.forEach((datum) => {
			const testResult = datum as TestResult;

			if (testResult.kind.case === 'warp' && testResult.kind.value.output && testResult.kind.value.input) {
				const groupName = this.generateWarpGroupName(testResult.kind.value.input);

				// Process get output
				if (testResult.kind.value.output.get) {
					const dataPoint: WarpDataPoint = {
						name: testResult.name,
						totalBytes: testResult.kind.value.output.get.totalBytes,
						totalObjects: testResult.kind.value.output.get.totalObjects,
						totalOperations: Number(testResult.kind.value.output.get.totalOperations),
						bytesFastest: testResult.kind.value.output.get.bytes?.fastestPerSecond || 0,
						bytesMedian: testResult.kind.value.output.get.bytes?.medianPerSecond || 0,
						bytesSlowest: testResult.kind.value.output.get.bytes?.slowestPerSecond || 0,
						objectsFastest: testResult.kind.value.output.get.objects?.fastestPerSecond || 0,
						objectsMedian: testResult.kind.value.output.get.objects?.medianPerSecond || 0,
						objectsSlowest: testResult.kind.value.output.get.objects?.slowestPerSecond || 0,
						completedAt: testResult.completedAt
							? new Date(
									Number(testResult.completedAt.seconds) * 1000 +
										Number(testResult.completedAt.nanos) / 1000000,
								)
							: new Date(),
					};

					if (!outputMap['get'][groupName]) {
						outputMap['get'][groupName] = {
							key: groupName,
							data: [],
							color: this.generateColor(groupName),
						};
					}

					outputMap['get'][groupName].data.push(dataPoint);
				}

				// Process put output
				if (testResult.kind.value.output.put) {
					const dataPoint: WarpDataPoint = {
						name: testResult.name,
						totalBytes: testResult.kind.value.output.put.totalBytes,
						totalObjects: testResult.kind.value.output.put.totalObjects,
						totalOperations: Number(testResult.kind.value.output.put.totalOperations),
						bytesFastest: testResult.kind.value.output.put.bytes?.fastestPerSecond || 0,
						bytesMedian: testResult.kind.value.output.put.bytes?.medianPerSecond || 0,
						bytesSlowest: testResult.kind.value.output.put.bytes?.slowestPerSecond || 0,
						objectsFastest: testResult.kind.value.output.put.objects?.fastestPerSecond || 0,
						objectsMedian: testResult.kind.value.output.put.objects?.medianPerSecond || 0,
						objectsSlowest: testResult.kind.value.output.put.objects?.slowestPerSecond || 0,
						completedAt: testResult.completedAt
							? new Date(
									Number(testResult.completedAt.seconds) * 1000 +
										Number(testResult.completedAt.nanos) / 1000000,
								)
							: new Date(),
					};

					if (!outputMap['put'][groupName]) {
						outputMap['put'][groupName] = {
							key: groupName,
							data: [],
							color: this.generateColor(groupName),
						};
					}

					outputMap['put'][groupName].data.push(dataPoint);
				}

				// Process delete output
				if (testResult.kind.value.output.delete) {
					const dataPoint: WarpDataPoint = {
						name: testResult.name,
						totalBytes: testResult.kind.value.output.delete.totalBytes,
						totalObjects: testResult.kind.value.output.delete.totalObjects,
						totalOperations: Number(testResult.kind.value.output.delete.totalOperations),
						bytesFastest: testResult.kind.value.output.delete.bytes?.fastestPerSecond || 0,
						bytesMedian: testResult.kind.value.output.delete.bytes?.medianPerSecond || 0,
						bytesSlowest: testResult.kind.value.output.delete.bytes?.slowestPerSecond || 0,
						objectsFastest: testResult.kind.value.output.delete.objects?.fastestPerSecond || 0,
						objectsMedian: testResult.kind.value.output.delete.objects?.medianPerSecond || 0,
						objectsSlowest: testResult.kind.value.output.delete.objects?.slowestPerSecond || 0,
						completedAt: testResult.completedAt
							? new Date(
									Number(testResult.completedAt.seconds) * 1000 +
										Number(testResult.completedAt.nanos) / 1000000,
								)
							: new Date(),
					};

					if (!outputMap['delete'][groupName]) {
						outputMap['delete'][groupName] = {
							key: groupName,
							data: [],
							color: this.generateColor(groupName),
						};
					}

					outputMap['delete'][groupName].data.push(dataPoint);
				}
			}
		});

		return {
			get: outputMap['get'],
			put: outputMap['put'],
			delete: outputMap['delete'],
		};
	}
}

export { BistDashboardManager, type FioDataPoint, type FioOutputGroup, type WarpDataPoint, type WarpOutputGroup };
