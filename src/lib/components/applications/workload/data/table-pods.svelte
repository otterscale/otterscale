<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { AreaChart } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext } from 'svelte';
	import { type Writable } from 'svelte/store';

	import { browser } from '$app/environment';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { Application, Application_Pod } from '$lib/api/application/v1/application_pb';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Empty } from '$lib/components/custom/table';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Chart from '$lib/components/ui/chart';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	import Actions from './cell-actions.svelte';
</script>

<script lang="ts">
	let {
		application,
		scope,
		namespace,
		reloadManager
	}: {
		application: Writable<Application>;
		scope: string;
		namespace: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let prometheusDriver = $state<PrometheusDriver | null>(null);
	async function fetchCPUMetrics(pod: Application_Pod) {
		if (!prometheusDriver) {
			const response = await environmentService.getPrometheus({});
			prometheusDriver = new PrometheusDriver({
				endpoint: '/prometheus',
				baseURL: response.baseUrl,
				headers: {
					'x-proxy-target': 'api'
				}
			});
		}

		const cpuResponse = await prometheusDriver.rangeQuery(
			`
			sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{namespace="${namespace}", pod="${pod.name}", juju_model="${scope}", container!=""}) by (container)
			`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);

		return (cpuResponse?.result[0]?.values as any[]) ?? [];
	}

	async function fetchMemoryMetrics(pod: Application_Pod) {
		if (!prometheusDriver) {
			const response = await environmentService.getPrometheus({});
			prometheusDriver = new PrometheusDriver({
				endpoint: '/prometheus',
				baseURL: response.baseUrl,
				headers: {
					'x-proxy-target': 'api'
				}
			});
		}

		const memoryResponse = await prometheusDriver.rangeQuery(
			`
			sum(container_memory_working_set_bytes{job="kubelet", metrics_path="/metrics/cadvisor", juju_model="${scope}", namespace="${namespace}", pod="${pod.name}", container!="", image!=""}) by (container)
			`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);

		return (memoryResponse?.result[0]?.values as any[]) ?? [];
	}

	async function fetchGPUMetrics(pod: Application_Pod) {
		if (!prometheusDriver) {
			const response = await environmentService.getPrometheus({});
			prometheusDriver = new PrometheusDriver({
				endpoint: '/prometheus',
				baseURL: response.baseUrl,
				headers: {
					'x-proxy-target': 'api'
				}
			});
		}

		const gpuResponse = await prometheusDriver.rangeQuery(
			`
			sum(Device_utilization_desc_of_container{juju_model="${scope}", namespace="${namespace}", podname="${pod.name}"}) by (deviceuuid, vdeviceid, podname, podnamespace)
			`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);

		return (gpuResponse?.result[0]?.values as any[]) ?? [];
	}

	function openTerminal(pod: Application_Pod) {
		if (!browser) {
			return;
		}

		const searchParams = new URLSearchParams({
			scope,
			namespace: page.params.namespace ?? '',
			pod: pod.name,
			container: '',
			command: '/bin/sh'
		});

		const terminalUrl = `${resolve('/tty')}?${searchParams.toString()}`;
		const windowName = m.tty();

		const features = [
			'width=800',
			'height=600',
			'toolbar=no',
			'location=no',
			'menubar=no',
			'status=no',
			'scrollbars=no',
			'resizable=yes'
		].join(',');

		const newWindow = window.open(terminalUrl, windowName, features);

		if (newWindow) {
			newWindow.focus();
		}
	}
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>
				{m.name()}
			</Table.Head>
			<Table.Head>
				{m.phase()}
			</Table.Head>
			<Table.Head>
				{m.ready()}
			</Table.Head>
			<Table.Head>
				{m.restarts()}
			</Table.Head>
			<Table.Head>
				{m.conditions()}
			</Table.Head>
			<Table.Head>
				{m.cpu()}
			</Table.Head>
			<Table.Head>
				{m.memory()}
			</Table.Head>
			<Table.Head>
				{m.gpu()}
			</Table.Head>
			<Table.Head>
				{m.terminal()}
			</Table.Head>
			<Table.Head></Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each $application.pods as pod, index (index)}
			<Table.Row>
				<Table.Cell>{pod.name}</Table.Cell>
				<Table.Cell>
					<Badge variant="outline">{pod.phase}</Badge>
				</Table.Cell>
				<Table.Cell>
					{pod.ready}
				</Table.Cell>
				<Table.Cell>{pod.restarts}</Table.Cell>
				<Table.Cell>
					{#if pod.conditions}
						{@const trueConditions = pod.conditions.filter(
							(condition) => condition.status === 'True'
						)}
						<div class="flex flex-wrap gap-1">
							{#each trueConditions as trueCondition, index (index)}
								<Tooltip.Provider>
									<Tooltip.Root>
										<Tooltip.Trigger>
											<Badge
												variant={['Failed', 'FailureTarget'].includes(trueCondition.type)
													? 'destructive'
													: 'outline'}
											>
												{trueCondition.type}
											</Badge>
										</Tooltip.Trigger>
										<Tooltip.Content>
											{#if trueCondition.message}
												{trueCondition.message}
											{:else}
												{trueCondition.type}
											{/if}
										</Tooltip.Content>
									</Tooltip.Root>
								</Tooltip.Provider>
							{/each}
						</div>
					{/if}
				</Table.Cell>

				<Table.Cell>
					{#await fetchCPUMetrics(pod) then usages}
						{@const maximumValue = Math.max(...usages.map((usage) => Number(usage.value)))}
						{@const minimumValue = Math.min(...usages.map((usage) => Number(usage.value)))}
						{@const configuration = {
							value: {
								label: 'usage',
								color: maximumValue > 0.5 ? 'var(--warning)' : 'var(--healthy)'
							}
						} satisfies Chart.ChartConfig}
						{#if usages.length > 0}
							<Layout.Cell class="relative justify-center">
								<Chart.Container config={configuration} class="h-10 w-full">
									<AreaChart
										data={usages}
										x="time"
										series={[
											{
												key: 'value',
												label: configuration['value'].label,
												color: configuration['value'].color
											}
										]}
										props={{
											area: {
												curve: curveNatural,
												'fill-opacity': 0.1,
												line: { class: 'stroke-1' },
												motion: 'tween'
											}
										}}
										axis={false}
										xScale={scaleUtc()}
										yDomain={[minimumValue, maximumValue]}
										grid={false}
									>
										{#snippet tooltip()}
											<Chart.Tooltip
												indicator="dot"
												labelFormatter={(v: Date) => {
													return v.toLocaleDateString('en-US', {
														year: 'numeric',
														month: 'short',
														day: 'numeric',
														hour: 'numeric',
														minute: 'numeric'
													});
												}}
											>
												{#snippet formatter({ item, name, value })}
													<div
														class="flex flex-1 shrink-0 items-center justify-start gap-1 font-mono text-xs leading-none"
														style="--color-bg: {item.color}"
													>
														<Icon icon="ph:square-fill" class="text-(--color-bg)" />
														<h1 class="font-semibold text-muted-foreground">{name}</h1>
														<p class="ml-auto">{(Number(value) * 100).toFixed(2)} %</p>
													</div>
												{/snippet}
											</Chart.Tooltip>
										{/snippet}
									</AreaChart>
								</Chart.Container>
							</Layout.Cell>
						{/if}
					{/await}
				</Table.Cell>
				<Table.Cell>
					{#await fetchMemoryMetrics(pod) then usages}
						{@const maximumValue = Math.max(...usages.map((usage) => Number(usage.value)))}
						{@const minimumValue = Math.min(...usages.map((usage) => Number(usage.value)))}
						{@const configuration = {
							value: {
								label: 'usage',
								color: maximumValue > 0.5 ? 'var(--warning)' : 'var(--healthy)'
							}
						} satisfies Chart.ChartConfig}
						{#if usages.length > 0}
							<Layout.Cell class="relative justify-center">
								<Chart.Container config={configuration} class="h-10 w-full">
									<AreaChart
										data={usages}
										x="time"
										series={[
											{
												key: 'value',
												label: configuration['value'].label,
												color: configuration['value'].color
											}
										]}
										props={{
											area: {
												curve: curveNatural,
												'fill-opacity': 0.1,
												line: { class: 'stroke-1' },
												motion: 'tween'
											}
										}}
										axis={false}
										xScale={scaleUtc()}
										yDomain={[minimumValue, maximumValue]}
										grid={false}
									>
										{#snippet tooltip()}
											<Chart.Tooltip
												indicator="dot"
												labelFormatter={(v: Date) => {
													return v.toLocaleDateString('en-US', {
														year: 'numeric',
														month: 'short',
														day: 'numeric',
														hour: 'numeric',
														minute: 'numeric'
													});
												}}
											>
												{#snippet formatter({ item, name, value })}
													{@const { value: capacity, unit } = formatCapacity(Number(value))}
													<div
														class="flex flex-1 shrink-0 items-center justify-start gap-1 font-mono text-xs leading-none"
														style="--color-bg: {item.color}"
													>
														<Icon icon="ph:square-fill" class="text-(--color-bg)" />
														<h1 class="font-semibold text-muted-foreground">{name}</h1>
														<p class="ml-auto">{capacity} {unit}</p>
													</div>
												{/snippet}
											</Chart.Tooltip>
										{/snippet}
									</AreaChart>
								</Chart.Container>
							</Layout.Cell>
						{/if}
					{/await}
				</Table.Cell>
				<Table.Cell>
					{#await fetchGPUMetrics(pod) then usages}
						{@const maximumValue = Math.max(...usages.map((usage) => Number(usage.value)))}
						{@const minimumValue = Math.min(...usages.map((usage) => Number(usage.value)))}
						{@const configuration = {
							value: {
								label: 'usage',
								color: maximumValue > 0.5 ? 'var(--warning)' : 'var(--healthy)'
							}
						} satisfies Chart.ChartConfig}
						{#if usages.length > 0}
							<Layout.Cell class="relative justify-center">
								<Chart.Container config={configuration} class="h-10 w-full">
									<AreaChart
										data={usages}
										x="time"
										series={[
											{
												key: 'value',
												label: configuration['value'].label,
												color: configuration['value'].color
											}
										]}
										props={{
											area: {
												curve: curveNatural,
												'fill-opacity': 0.1,
												line: { class: 'stroke-1' },
												motion: 'tween'
											}
										}}
										axis={false}
										xScale={scaleUtc()}
										yDomain={[minimumValue, maximumValue]}
										grid={false}
									>
										{#snippet tooltip()}
											<Chart.Tooltip
												indicator="dot"
												labelFormatter={(v: Date) => {
													return v.toLocaleDateString('en-US', {
														year: 'numeric',
														month: 'short',
														day: 'numeric',
														hour: 'numeric',
														minute: 'numeric'
													});
												}}
											>
												{#snippet formatter({ item, name, value })}
													<div
														class="flex flex-1 shrink-0 items-center justify-start gap-1 font-mono text-xs leading-none"
														style="--color-bg: {item.color}"
													>
														<Icon icon="ph:square-fill" class="text-(--color-bg)" />
														<h1 class="font-semibold text-muted-foreground">{name}</h1>
														<p class="ml-auto">{(Number(value) * 100).toFixed(2)} %</p>
													</div>
												{/snippet}
											</Chart.Tooltip>
										{/snippet}
									</AreaChart>
								</Chart.Container>
							</Layout.Cell>
						{/if}
					{/await}
				</Table.Cell>
				<Table.Cell>
					<Button variant="secondary" size="icon" onclick={() => openTerminal(pod)}>
						<Icon icon="ph:terminal-window" />
					</Button>
				</Table.Cell>
				<Table.Cell>
					<Actions {pod} {scope} {namespace} {reloadManager} />
				</Table.Cell>
			</Table.Row>
		{/each}
		{#if $application.pods.length === 0}
			<Table.Row>
				<Table.Cell colspan={6}>
					<Empty />
				</Table.Cell>
			</Table.Row>
		{/if}
	</Table.Body>
</Table.Root>
