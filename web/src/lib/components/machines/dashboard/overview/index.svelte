<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import CPU from './cpu.svelte';
	import GPU from './gpu.svelte';
	import GPUMemory from './gpu-memory.svelte';
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
	<div class="col-span-2">
		<CPU {prometheusDriver} bind:isReloading />
	</div>
	<div class="col-span-2">
		<Memory {prometheusDriver} bind:isReloading />
	</div>
	<div class="col-span-2">
		<Storage {prometheusDriver} bind:isReloading />
	</div>
	<div class="col-span-3 row-span-2">
		<Nodes bind:isReloading />
	</div>
	<div class="col-span-2">
		<GPU {prometheusDriver} bind:isReloading />
	</div>
	<div class="col-span-2">
		<GPUMemory {prometheusDriver} bind:isReloading />
	</div>
	<div class="col-span-6">
		<SystemLoad {prometheusDriver} bind:isReloading />
	</div>
	<div class="col-span-3">
		<NodeProportion bind:isReloading />
	</div>
</div>
