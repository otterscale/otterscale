<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import Latency from './latency.svelte';
	import Model from './model.svelte';
	import Request from './request.svelte';
	import Throughput from './throughtput.svelte';
	import TimeToFirstToken from './time-to-first-token.svelte';
	import VGPU from './vgpu.svelte';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();
</script>

<div class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-8">
	<div class="col-span-2 grid gap-5">
		<div class="col-span-2">
			<Model {prometheusDriver} {scope} bind:isReloading />
		</div>
		<div class="col-span-2">
			<Latency {prometheusDriver} {scope} bind:isReloading />
		</div>
	</div>

	<div class="col-span-2 h-full">
		<VGPU {prometheusDriver} {scope} bind:isReloading />
	</div>

	<div class="col-span-4">
		<TimeToFirstToken {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-4 row-start-2">
		<Throughput {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-4 row-start-2">
		<Request {prometheusDriver} {scope} bind:isReloading />
	</div>
</div>
