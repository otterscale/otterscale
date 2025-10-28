import type { Table } from '@tanstack/table-core';

import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';

interface FioDataPoint {
	name: string;
	ioBytes: number;
	bandwidthBytes: number;
	ioPerSecond: number;
	totalIos: number;
	latency: number;
	completedAt: Date;
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

	getFioOutputs(): {
		read: FioDataPoint[];
		write: FioDataPoint[];
		trim: FioDataPoint[];
	} {
		const outputMap: Record<string, FioDataPoint[]> = {
			read: [],
			write: [],
			trim: [],
		};

		this.filteredData.forEach((datum) => {
			const testResult = datum as TestResult;

			if (testResult.kind.case === 'fio' && testResult.kind.value.output && testResult.kind.value.input) {
				// Process read output
				if (testResult.kind.value.output.read && testResult.kind.value.output.read.latency) {
					outputMap['read'].push({
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
					});
				}

				// Process write output
				if (testResult.kind.value.output.write && testResult.kind.value.output.write.latency) {
					outputMap['write'].push({
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
					});
				}

				// Process trim output
				if (testResult.kind.value.output.trim && testResult.kind.value.output.trim.latency) {
					outputMap['trim'].push({
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
					});
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
		get: WarpDataPoint[];
		put: WarpDataPoint[];
		delete: WarpDataPoint[];
	} {
		const outputMap: Record<string, WarpDataPoint[]> = {
			get: [],
			put: [],
			delete: [],
		};

		this.filteredData.forEach((datum) => {
			const testResult = datum as TestResult;

			if (testResult.kind.case === 'warp' && testResult.kind.value.output && testResult.kind.value.input) {
				// Process get output
				if (testResult.kind.value.output.get) {
					outputMap['get'].push({
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
					});
				}

				// Process put output
				if (testResult.kind.value.output.put) {
					outputMap['put'].push({
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
					});
				}

				// Process delete output
				if (testResult.kind.value.output.delete) {
					outputMap['delete'].push({
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
					});
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

export { BistDashboardManager, type FioDataPoint, type WarpDataPoint, type WarpOutputGroup };
