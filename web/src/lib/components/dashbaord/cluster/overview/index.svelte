<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import CpuUsage from './cpu.svelte';
	import GPUMemorUsage from './gpu-memory-usage.svelte';
	import GPUUtilization from './gpu-utilization.svelte';
	import Health from './health.svelte';
	import MemoryUsage from './memory2.svelte';
	import ClusterNodes from './nodes.svelte';
	import Pods from './pods.svelte';
	import Uptime from './uptime.svelte';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();
</script>

<div class="grid auto-rows-[minmax(140px,auto)] grid-cols-2 gap-4 pt-4 md:gap-6 lg:grid-cols-4">
	<Health {prometheusDriver} {scope} bind:isReloading />
	<Uptime {prometheusDriver} {scope} bind:isReloading />
	<!-- <Version {prometheusDriver} {scope} bind:isReloading /> -->
	<Pods {prometheusDriver} {scope} bind:isReloading />
	<ClusterNodes {prometheusDriver} {scope} bind:isReloading />

	<div class="col-span-2">
		<CpuUsage {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-2">
		<MemoryUsage {prometheusDriver} {scope} bind:isReloading />
	</div>

	<div class="col-span-2">
		<GPUUtilization {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-2">
		<GPUMemorUsage {prometheusDriver} {scope} bind:isReloading />
	</div>
</div>
