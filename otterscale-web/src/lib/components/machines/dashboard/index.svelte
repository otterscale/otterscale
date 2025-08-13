<script lang="ts">
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext } from 'svelte';
	import { default as CPUAverage } from './area-chart-cpu-average.svelte';
	import { default as CPUCoreProcessor } from './area-chart-cpu-core.svelte';
	import { default as DiskIOTime } from './area-chart-disk-io-time.svelte';
	import { default as BasicDisk } from './area-chart-disk-rw.svelte';
	import { default as NetworkReceived } from './area-chart-network-received.svelte';
	import { default as NetworkTransmitted } from './area-chart-network-transmitted.svelte';
	import { default as BasicRAM } from './area-chart-ram.svelte';
	import { default as UsageRateUptime } from './text-chart-uptime.svelte';
	import { default as UsageRateCPU } from './usage-rate-chart-cpu.svelte';
	import { default as UsageRateRAM } from './usage-rate-chart-ram.svelte';
	import { default as UsageRateRootFS } from './usage-rate-chart-root-fs.svelte';
	import { default as UsageRateSWAP } from './usage-rate-chart-swap.svelte';

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
				<UsageRateUptime {client} machine={selectedMachine} />
			</span>
		</div>
		<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
			<span class="col-span-1">
				<UsageRateCPU {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<UsageRateRAM {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<UsageRateSWAP {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<UsageRateRootFS {client} machine={selectedMachine} />
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
