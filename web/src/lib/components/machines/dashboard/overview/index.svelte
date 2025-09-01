<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import CPU from './cpu.svelte';
	import Memory from './memory.svelte';
	import NodeProportion from './node-proportion.svelte';
	import Nodes from './nodes.svelte';
	import Storage from './storage.svelte';
	import SystemLoad from './system_load.svelte';

	let {
		prometheusDriver,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; isReloading: boolean } = $props();
</script>

<div class="grid auto-rows-auto grid-cols-3 gap-5 pt-4 md:grid-cols-6 lg:grid-cols-9">
	<CPU {prometheusDriver} {isReloading} span="col-span-2" />
	<Memory {prometheusDriver} {isReloading} span="col-span-2" />
	<Storage {prometheusDriver} {isReloading} span="col-span-2" />
	<Nodes {isReloading} span="col-span-3" />
	<SystemLoad {prometheusDriver} {isReloading} span="col-span-6" />
	<NodeProportion {isReloading} span="col-span-3" />
</div>
