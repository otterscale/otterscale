<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import InstancePicker from './utils/instance-picker.svelte';

	import QuickUptime from './uptime/quick-metric.svelte';
	import QuickCPU from './cpu/quick-metric.svelte';
	import BasicCPU from './cpu/basic-metric.svelte';
	import BasicRAM from './ram/basic-metric.svelte';
	import QuickRAM from './ram/quick-metric.svelte';
	import QuickSWAP from './swap/quick-metric.svelte';
	import BasicSWAP from './swap/basic-metric.svelte';
	import QuickRootFS from './root-fs/quick-metric.svelte';

	import NetworkTrafficBasic from './network/traffic-basic.svelte';

	let {
		client,
		juju_model_uuid,
		instances
	}: { client: PrometheusDriver; juju_model_uuid: string; instances: string[] } = $props();

	let selectedInstance = $state(instances[0]);
</script>

<main class="grid gap-4 p-4">
	<span class="ml-auto flex items-center gap-2">
		<InstancePicker bind:selectedInstance {instances} />
	</span>

	<p class="text-xl font-bold">Quick Metric</p>
	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-5">
		<span class="col-span-1">
			{#key selectedInstance}
				<QuickUptime {client} {juju_model_uuid} instance={selectedInstance} />
			{/key}
		</span>
		<span class="col-span-1">
			{#key selectedInstance}
				<QuickCPU {client} {juju_model_uuid} instance={selectedInstance} />
			{/key}
		</span>
		<span class="col-span-1">
			{#key selectedInstance}
				<QuickRAM {client} {juju_model_uuid} instance={selectedInstance} />
			{/key}
		</span>
		<span class="col-span-1">
			{#key selectedInstance}
				<QuickSWAP {client} {juju_model_uuid} instance={selectedInstance} />
			{/key}
		</span>
		<span class="col-span-1">
			{#key selectedInstance}
				<QuickRootFS {client} {juju_model_uuid} instance={selectedInstance} />
			{/key}
		</span>
	</div>
	<p class="text-xl font-bold">Basic Metric</p>
	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
		<span class="col-span-1">
			{#key selectedInstance}
				<BasicCPU {client} {juju_model_uuid} instance={selectedInstance} />
			{/key}
		</span>
		<span class="col-span-1">
			{#key selectedInstance}
				<BasicRAM {client} {juju_model_uuid} instance={selectedInstance} />
			{/key}
		</span>
		<span class="col-span-1">
			{#key selectedInstance}
				<BasicSWAP {client} {juju_model_uuid} instance={selectedInstance} />
			{/key}
		</span>
		<span class="col-span-1">
			{#key selectedInstance}
				<NetworkTrafficBasic {client} {juju_model_uuid} instance={selectedInstance} />
			{/key}
		</span>
	</div>
</main>
