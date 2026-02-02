<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import AidaptivCache from './aidaptiv-cache.svelte';
	import ClusterNodes from './cluster-nodes.svelte';
	import CpuUsage from './cpu-usage.svelte';
	import Nvidia from './nvidia.svelte';
	import MemoryUsage from './memory-usage.svelte';
	import Pods from './pods.svelte';
	import Uptime from './uptime.svelte';
	import Version from './version.svelte';
	import Health from './health.svelte';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();
</script>

<div class="grid auto-rows-[minmax(140px,auto)] grid-cols-2 gap-4 pt-4 md:gap-6 lg:grid-cols-4">
	<Health {prometheusDriver} {scope} bind:isReloading />
	<Version {prometheusDriver} {scope} bind:isReloading />
	<Uptime {prometheusDriver} {scope} bind:isReloading />
	<ClusterNodes {prometheusDriver} {scope} bind:isReloading />

	<div class="col-span-2">
		<CpuUsage {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-2">
		<MemoryUsage {prometheusDriver} {scope} bind:isReloading />
	</div>

	<div class="col-span-2">
		<Pods {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-2">
		<Nvidia {prometheusDriver} {scope} bind:isReloading />
	</div>
</div>
