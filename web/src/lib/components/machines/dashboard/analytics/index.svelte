<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import { default as CPUAverage } from './area-chart-cpu-average.svelte';
	import { default as CPUCoreProcessor } from './area-chart-cpu-core.svelte';
	import { default as DiskIOTime } from './area-chart-disk-io-time.svelte';
	import { default as BasicDisk } from './area-chart-disk-rw.svelte';
	import { default as NetworkReceived } from './area-chart-network-received.svelte';
	import { default as NetworkTransmitted } from './area-chart-network-transmitted.svelte';
	import { default as BasicRAM } from './area-chart-ram.svelte';
	import FQDNPicker from './machine-picker.svelte';
	import { default as UsageRateUptime } from './text-chart-uptime.svelte';
	import { default as UsageRateCPU } from './usage-rate-chart-cpu.svelte';
	import { default as UsageRateRAM } from './usage-rate-chart-ram.svelte';
	import { default as UsageRateRootFS } from './usage-rate-chart-root-fs.svelte';
	import { default as UsageRateSWAP } from './usage-rate-chart-swap.svelte';

	let { selectedTab, client }: { selectedTab: string; client: PrometheusDriver } = $props();

	let selectedFQDN = $state(undefined);
</script>

<div class="flex flex-col gap-4">
	<div class="mr-auto flex flex-wrap items-center gap-2">
		<FQDNPicker bind:selectedFQDN />
	</div>
	{#if selectedFQDN && selectedTab === 'analytics'}
		{#key selectedFQDN}
			<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
				<span class="col-span-1">
					<UsageRateUptime {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<UsageRateCPU {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<UsageRateRAM {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<UsageRateSWAP {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<UsageRateRootFS {client} fqdn={selectedFQDN} />
				</span>
			</div>

			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-1 lg:grid-cols-2">
				<span class="col-span-1">
					<CPUCoreProcessor {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<CPUAverage {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<BasicRAM {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<BasicDisk {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<DiskIOTime {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<NetworkReceived {client} fqdn={selectedFQDN} />
				</span>
				<span class="col-span-1">
					<NetworkTransmitted {client} fqdn={selectedFQDN} />
				</span>
			</div>
		{/key}
	{/if}
</div>
