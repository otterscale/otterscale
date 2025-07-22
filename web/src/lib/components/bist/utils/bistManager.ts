import type { Table } from '@tanstack/table-core';
import type { TestResult } from '$gen/api/bist/v1/bist_pb';
import { FIO_Input_AccessMode } from '$gen/api/bist/v1/bist_pb';

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

class BistDashboardManager<TData = TestResult> {
    table: Table<TData>;

    constructor(table: Table<TData>) {
        this.table = table;
    }

    get filteredData() {
        return this.table.getFilteredRowModel().rows.map((row) => row.original);
    }

    private hashCode(str: string): number {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            hash = ((hash << 5) - hash) + str.charCodeAt(i);
            hash |= 0;
        }
        return hash;
    }

    private generateGroupName(input: any): string {
        return `${FIO_Input_AccessMode[input.accessMode]}-${Number(input.jobCount)}-${input.runTime}-${input.blockSize}-${input.fileSize}-${Number(input.ioDepth)}`;
    }

    private generateColor(groupName: string): string {
        const hue = Math.abs(this.hashCode(groupName)) % 360;
        return `hsl(${hue}, 70%, 50%)`;
    }

    getFioOutputs(): { read: Record<string, FioOutputGroup>, write: Record<string, FioOutputGroup>, trim: Record<string, FioOutputGroup> } {
        const outputMap: Record<string, Record<string, FioOutputGroup>> = {
            read: {},
            write: {},
            trim: {}
        };

        this.filteredData.forEach((datum) => {
            const testResult = datum as TestResult;

            if (testResult.kind.case === 'fio' && 
                testResult.kind.value.output && 
                testResult.kind.value.input
            ) {
                const groupName = this.generateGroupName(testResult.kind.value.input);
                
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
                                Number(testResult.completedAt.nanos) / 1000000
                            )
                            : new Date(),
                    };

                    if (!outputMap['read'][groupName]) {
                        outputMap['read'][groupName] = {
                            key: groupName,
                            data: [],
                            color: this.generateColor(groupName)
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
                                Number(testResult.completedAt.nanos) / 1000000
                            )
                            : new Date(),
                    };

                    if (!outputMap['write'][groupName]) {
                        outputMap['write'][groupName] = {
                            key: groupName,
                            data: [],
                            color: this.generateColor(groupName)
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
                                Number(testResult.completedAt.nanos) / 1000000
                            )
                            : new Date(),
                    };

                    if (!outputMap['trim'][groupName]) {
                        outputMap['trim'][groupName] = {
                            key: groupName,
                            data: [],
                            color: this.generateColor(groupName)
                        };
                    }
                    
                    outputMap['trim'][groupName].data.push(dataPoint);
                }
            }
        });
        return {
            read: outputMap['read'],
            write: outputMap['write'],
            trim: outputMap['trim']
        };
    }
}

export {
    BistDashboardManager,
    type FioDataPoint,
    type FioOutputGroup
};
