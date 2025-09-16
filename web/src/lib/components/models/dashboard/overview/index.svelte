<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import Latency from './latency.svelte';
	import Model from './model.svelte';
	import Request from './request.svelte';
	import Throughput from './throughtput.svelte';
	import TimeToFirstToken from './time-to-first-token.svelte';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';

	let {
		prometheusDriver: _,
		scope,
		isReloading = $bindable(),
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	const prometheusDriver = new PrometheusDriver({
		endpoint: 'http://192.168.41.100:30091',
		baseURL: '/api/v1',
	});
</script>

<div class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-8">
	<div class="col-span-2">
		<Model {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-2">
		<Latency {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-4">
		<TimeToFirstToken {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-4 row-start-2">
		<Request {prometheusDriver} {scope} bind:isReloading />
	</div>
	<div class="col-span-4 row-start-2">
		<Throughput {prometheusDriver} {scope} bind:isReloading />
	</div>
</div>
