<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { PrometheusDriver } from 'prometheus-query';
	import { default as AreaErrorBudget } from './area-chart-error-budget.svelte';
	import { default as AreaVolumnSpace } from './area-chart-volumn-space.svelte';
	import { default as TextRunningKubelet } from './text-chart-running-kubelet.svelte';
	import { default as TextRunningKubeletPod } from './text-chart-running-kubelet-pod.svelte';
	import { default as TextRunningKubeletContainer } from './text-chart-running-kubelet-container.svelte';
	import { default as TextUpControllerManager } from './text-chart-up-controller-manager.svelte';
	import { default as TextUpETCD } from './text-chart-up-etcd.svelte';
	import { default as TextUpProxy } from './text-chart-up-proxy.svelte';
	import { default as TextUpScheduler } from './text-chart-up-scheduler.svelte';
	import { default as UsageAvailability } from './usage-rate-chart-availability.svelte';
	import { default as UsageVolumnSpace } from './usage-rate-chart-volumn-space.svelte';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<div class="flex flex-col gap-4">
	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
		<span class="col-span-1">
			<TextRunningKubelet {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextRunningKubeletPod {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextRunningKubeletContainer {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-1">
			<TextUpControllerManager {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextUpETCD {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextUpProxy {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextUpScheduler {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-2">
			<UsageAvailability {client} {scope} />
		</span>
		<span class="col-span-2">
			<AreaErrorBudget {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-2">
			<UsageVolumnSpace {client} {scope} />
		</span>
		<span class="col-span-2">
			<AreaVolumnSpace {client} {scope} />
		</span>
	</div>
</div>
