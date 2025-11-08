<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import CPU from './cpu.svelte';
	import GPUMemory from './gpu-memory.svelte';
	import GPU from './gpu.svelte';
	import Memory from './memory.svelte';
	import NodeProportion from './node-proportion.svelte';
	import Nodes from './nodes.svelte';
	import Storage from './storage.svelte';
	import SystemLoad from './system_load.svelte';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();
</script>

<div class="grid auto-rows-auto grid-cols-3 gap-5 pt-4 md:grid-cols-6 lg:grid-cols-9">
	<div class="col-span-2">
		<CPU {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-2">
		<Memory {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-2">
		<Storage {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-3 row-span-2">
		<Nodes {scope} bind:isReloading />
	</div>
	<div class="col-span-2">
		<GPU {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-2">
		<GPUMemory {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-6">
		<SystemLoad {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-3">
		<NodeProportion {scope} bind:isReloading />
	</div>
</div>
