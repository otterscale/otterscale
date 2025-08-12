<script lang="ts">
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	// import { DateTimestampPicker, type TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import { getLocalTimeZone, now } from '@internationalized/date';
	import { PrometheusDriver } from 'prometheus-query';
	// import * as Pickers from '../utils/pickers';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { type Writable } from 'svelte/store';
	import { default as QuickUptime } from './quick-metric-uptime.svelte';
	import { default as QuickCPU } from './quick-metric-cpu.svelte';
	import { default as QuickRAM } from './quick-metric-ram.svelte';
	import { default as QuickSWAP } from './quick-metric-swap.svelte';
	import { default as QuickRootFS } from './quick-metric-root-fs.svelte';
	import { default as CPUCoreProcessor } from './basic-metric-cpu-core.svelte';
	import { default as CPUAverage } from './basic-metric-cpu-average.svelte';
	import { default as BasicRAM } from './basic-metric-ram.svelte';
	import { default as BasicDisk } from './basic-metric-disk.svelte';
	import { default as DiskIOTime } from './basic-metric-disk-io-time.svelte';
	import { default as NetworkReceived } from './basic-metric-network-received.svelte';
	import { default as NetworkTransmitted } from './basic-metric-network-transmitted.svelte';

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	// const client: Promise<PrometheusDriver> = getContext('prometheusDriver');

	let { client }: { client: PrometheusDriver } = $props();
	// let selectedTimeRange = $state({
	// 	start: new Date(now(getLocalTimeZone()).toDate().getTime() - 60 * 60 * 1000),
	// 	end: now(getLocalTimeZone()).toDate()
	// } as TimeRange);
	let fakeMachine: Machine = {
		id: 'fake-id',
		fqdn: 'ottersacle-vm143.maas',
		description: 'This is a fake machine for demo purposes'
		// Add other required Machine fields with mock values as needed
	} as Machine;
	// let selectedMachine = $state(machines[0]);
	let selectedMachine = $state(fakeMachine);
</script>

<div class="flex flex-col gap-4">
	{#key selectedMachine}
		<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
			<span class="col-span-1">
				<QuickUptime {client} machine={selectedMachine} />
			</span>
		</div>
		<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
			<span class="col-span-1">
				<QuickCPU {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<QuickRAM {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<QuickSWAP {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<QuickRootFS {client} machine={selectedMachine} />
			</span>
		</div>

		<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2">
			<span class="col-span-1">
				<!-- <Example /> -->
				<CPUCoreProcessor {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<CPUAverage {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<BasicRAM {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<BasicDisk {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<DiskIOTime {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<NetworkReceived {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<NetworkTransmitted {client} machine={selectedMachine} />
			</span>
		</div>
	{/key}
</div>
