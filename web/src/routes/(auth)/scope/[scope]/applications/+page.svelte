<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { FacilityService, type Facility } from '$lib/api/facility/v1/facility_pb';
	import { Dashboard } from '$lib/components/applications/dashboard';
	import Loading from '$lib/components/custom/loading/report.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import { formatCapacity, formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb, currentKubernetes } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { ReloadManager, Reloader } from '$lib/components/custom/reloader';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.applications(page.params.scope) });

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);
	const facilityClient = createClient(FacilityService, transport);
	const prometheusDriver: Writable<PrometheusDriver | null> = writable(null);

	let isMounted = $state(false);

	const facilities = writable<Facility[]>([]);

	const controlPlane = $derived(
		$facilities.find((facility) => facility.name.includes('kubernetes-control-plane') && facility.units.length > 0),
	);
	const controlPlaneUnits = $derived(controlPlane?.units ?? []);
	const activeControlPlaneUnits = $derived(
		controlPlaneUnits.filter((unit) => unit.workloadStatus?.state === 'active') ?? [],
	);

	const worker = $derived(
		$facilities.find((facility) => facility.name.includes('kubernetes-worker') && facility.units.length > 0),
	);
	const workerUnits = $derived(worker?.units ?? []);
	const activeWorkerUnits = $derived(workerUnits.filter((unit) => unit.workloadStatus?.state === 'active') ?? []);
	let runningPods = $state(0);
	let runningContainers = $state(0);

	let cpuUsages: SampleValue[] = $state([]);
	const cpuUsagesConfiguration = {
		usage: { label: 'Usage', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;
	let cpuRequests = $state(0);
	let cpuLimits = $state(0);

	let memoryUsages: SampleValue[] = $state([]);
	const memoryUsagesConfiguration = {
		usage: { label: 'Usage', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig;
	let memoryRequests = $state(0);
	let memoryLimits = $state(0);

	let receives = $state([] as SampleValue[]);
	let transmits = $state([] as SampleValue[]);
	const traffics = $derived(
		receives.map((sample, index) => ({
			time: sample.time,
			receive: sample.value,
			transmit: transmits[index]?.value ?? 0,
		})),
	);
	const trafficsConfigurations = {
		receive: { label: 'Receive', color: 'var(--chart-1)' },
		transmit: { label: 'Transmit', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;

	let reads = $state([] as SampleValue[]);
	let writes = $state([] as SampleValue[]);
	const throughputs = $derived(
		reads.map((sample, index) => ({
			time: sample.time,
			read: sample.value,
			write: writes[index]?.value ?? 0,
		})),
	);
	const throughputsConfigurations = {
		read: { label: 'Read', color: 'var(--chart-1)' },
		write: { label: 'Write', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;

	function fetch() {
		environmentService.getPrometheus({}).then((response) => {
			prometheusDriver.set(
				new PrometheusDriver({
					endpoint: `${env.PUBLIC_API_URL}/prometheus`,
					baseURL: response.baseUrl,
				}),
			);

			if ($prometheusDriver) {
				$prometheusDriver
					.instantQuery(
						`
						sum(
							kubelet_running_pods{job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics"}
						)
						or
						sum(
							kubelet_running_pod_count{job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics"}
						)
						`,
					)
					.then((response) => {
						runningPods = response.result[0].value.value;
					});
				$prometheusDriver
					.instantQuery(
						`
						sum(
							kubelet_running_containers{job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics"}
						)
						or
						sum(
							kubelet_running_container_count{job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics"}
						)
						`,
					)
					.then((response) => {
						runningContainers = response.result[0].value.value;
					});
				$prometheusDriver
					.rangeQuery(
						`
						sum(
						node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{juju_model_uuid="${$currentKubernetes?.scopeUuid}"}
						)
						`,
						Date.now() - 60 * 60 * 1000,
						Date.now(),
						2 * 60,
					)
					.then((response) => {
						cpuUsages = response.result[0].values;
					});
				$prometheusDriver
					.instantQuery(
						`
						sum(
							namespace_cpu:kube_pod_container_resource_requests:sum{juju_model_uuid="${$currentKubernetes?.scopeUuid}"}
						)
						/
						sum(
							kube_node_status_allocatable{job="kube-state-metrics",juju_model_uuid="${$currentKubernetes?.scopeUuid}",resource="cpu"}
						)
						`,
					)
					.then((response) => {
						cpuRequests = response.result[0].value.value;
					});
				$prometheusDriver
					.instantQuery(
						`
						sum(
							namespace_cpu:kube_pod_container_resource_limits:sum{juju_model_uuid="${$currentKubernetes?.scopeUuid}"}
						)
						/
						sum(
							kube_node_status_allocatable{job="kube-state-metrics",juju_model_uuid="${$currentKubernetes?.scopeUuid}",resource="cpu"}
						)
						`,
					)
					.then((response) => {
						cpuLimits = response.result[0].value.value;
					});

				$prometheusDriver
					.rangeQuery(
						`
						sum(
						container_memory_rss{container!="",job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics/cadvisor"}
						)
						`,
						Date.now() - 60 * 60 * 1000,
						Date.now(),
						2 * 60,
					)
					.then((response) => {
						memoryUsages = response.result[0].values;
					});
				$prometheusDriver
					.instantQuery(
						`
						sum(
							namespace_memory:kube_pod_container_resource_requests:sum{juju_model_uuid="${$currentKubernetes?.scopeUuid}"}
						)
						/
						sum(
							kube_node_status_allocatable{job="kube-state-metrics",juju_model_uuid="${$currentKubernetes?.scopeUuid}",resource="memory"}
						)
						`,
					)
					.then((response) => {
						memoryRequests = response.result[0].value.value;
					});
				$prometheusDriver
					.instantQuery(
						`
						sum(
							namespace_memory:kube_pod_container_resource_limits:sum{juju_model_uuid="${$currentKubernetes?.scopeUuid}"}
						)
						/
						sum(
							kube_node_status_allocatable{job="kube-state-metrics",juju_model_uuid="${$currentKubernetes?.scopeUuid}",resource="memory"}
						)
						`,
					)
					.then((response) => {
						memoryLimits = response.result[0].value.value;
					});

				$prometheusDriver
					.rangeQuery(
						`
						sum(
						irate(
							container_network_receive_bytes_total{job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics/cadvisor",namespace=~".+"}[4m]
						)
						)
						`,
						new Date().setMinutes(0, 0, 0) - 1 * 60 * 60 * 1000,
						new Date().setMinutes(0, 0, 0),
						2 * 60,
					)
					.then((response) => {
						receives = response.result[0].values;
					});
				$prometheusDriver
					.rangeQuery(
						`
						sum(
						irate(
							container_network_transmit_bytes_total{job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics/cadvisor",namespace=~".+"}[4m]
						)
						)
						`,
						new Date().setMinutes(0, 0, 0) - 1 * 60 * 60 * 1000,
						new Date().setMinutes(0, 0, 0),
						2 * 60,
					)
					.then((response) => {
						transmits = response.result[0].values;
					});

				$prometheusDriver
					.rangeQuery(
						`
						sum(
						rate(
							container_fs_reads_bytes_total{container!="",device=~"(/dev/)?(mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|md.+|dasd.+)",job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics/cadvisor",namespace!=""}[4m]
						)
						)
						`,
						new Date().setMinutes(0, 0, 0) - 1 * 60 * 60 * 1000,
						new Date().setMinutes(0, 0, 0),
						2 * 60,
					)
					.then((response) => {
						reads = response.result[0].values;
					});
				$prometheusDriver
					.rangeQuery(
						`
						sum(
						rate(
							container_fs_writes_bytes_total{container!="",job="kubelet",juju_model_uuid="${$currentKubernetes?.scopeUuid}",metrics_path="/metrics/cadvisor",namespace!=""}[4m]
						)
						)
						`,
						new Date().setMinutes(0, 0, 0) - 1 * 60 * 60 * 1000,
						new Date().setMinutes(0, 0, 0),
						2 * 60,
					)
					.then((response) => {
						writes = response.result[0].values;
					});
			}
		});

		facilityClient.listFacilities({ scopeUuid: $currentKubernetes?.scopeUuid }).then((response) => {
			facilities.set(response.facilities);
		});
	}

	const reloadManager = new ReloadManager(fetch);

	onMount(async () => {
		try {
			await fetch();
			isMounted = true;
			reloadManager.start();
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		}
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<div class="mx-auto grid w-full gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.dashboard()}</h1>
		<p class="text-muted-foreground">
			{m.management_dashboard_description()}
		</p>
	</div>

	<Tabs.Root value="overview">
		<div class="flex justify-between gap-2">
			<Tabs.List>
				<Tabs.Trigger value="overview">{m.overview()}</Tabs.Trigger>
				<Tabs.Trigger value="analytics" disabled>{m.analytics()}</Tabs.Trigger>
			</Tabs.List>
			<Reloader
				bind:checked={reloadManager.state}
				onCheckedChange={() => {
					if (reloadManager.state) {
						reloadManager.restart();
					} else {
						reloadManager.stop();
					}
				}}
			/>
		</div>
		<Tabs.Content
			value="overview"
			class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-10"
		>
			<Card.Root class="relative col-span-2 gap-2 overflow-hidden">
				<Icon
					icon="ph:compass"
					class="text-primary/5 absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				/>
				<Card.Header>
					<Card.Title>{m.control_planes()}</Card.Title>
					<Card.Description>{m.ready()}</Card.Description>
				</Card.Header>
				<Card.Content class="text-3xl">
					{activeControlPlaneUnits.length} / {controlPlaneUnits.length}
				</Card.Content>
			</Card.Root>

			<Card.Root class="relative col-span-2 gap-2 overflow-hidden">
				<Icon
					icon="ph:cube"
					class="text-primary/5 absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				/>
				<Card.Header>
					<Card.Title>{m.workers()}</Card.Title>
					<Card.Description>{m.ready()}</Card.Description>
				</Card.Header>
				<Card.Content class="text-3xl">
					{activeWorkerUnits.length} / {workerUnits.length}
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-3 row-span-2 gap-2">
				<Card.Header>
					<Card.Title>{m.cpu_usage()}</Card.Title>
					<Card.Action class="text-muted-foreground flex flex-col gap-0.5 text-sm">
						<div class="flex justify-between gap-2">
							<p>{m.requests()}</p>
							<p class="font-mono">{Math.round(cpuRequests * 100)}%</p>
						</div>
						<div class="flex justify-between gap-2">
							<p>{m.limits()}</p>
							<p class="font-mono">{Math.round(cpuLimits * 100)}%</p>
						</div>
					</Card.Action>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={cpuUsagesConfiguration}>
						<AreaChart
							data={cpuUsages}
							x="time"
							xScale={scaleUtc()}
							yPadding={[0, 25]}
							series={[
								{
									key: 'value',
									label: cpuUsagesConfiguration.usage.label,
									color: cpuUsagesConfiguration.usage.color,
								},
							]}
							seriesLayout="stack"
							props={{
								area: {
									curve: curveNatural,
									'fill-opacity': 0.4,
									line: { class: 'stroke-1' },
									motion: 'tween',
								},
								xAxis: {
									format: (v: Date) =>
										`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`,
								},
								yAxis: { format: () => '' },
							}}
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
											minute: 'numeric',
										});
									}}
								>
									{#snippet formatter({ item, name, value })}
										<div
											style="--color-bg: {item.color}"
											class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
										></div>
										<div
											class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
										>
											<div class="grid gap-1.5">
												<span class="text-muted-foreground">{name}</span>
											</div>
											<p class="font-mono">{(Number(value) * 100).toFixed(2)}%</p>
										</div>
									{/snippet}
								</Chart.Tooltip>
							{/snippet}
							{#snippet marks({ series, getAreaProps })}
								{#each series as s, i (s.key)}
									<LinearGradient
										stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
										vertical
									>
										{#snippet children({ gradient })}
											<Area {...getAreaProps(s, i)} fill={gradient} />
										{/snippet}
									</LinearGradient>
								{/each}
							{/snippet}
						</AreaChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-3 row-span-2 gap-2">
				<Card.Header>
					<Card.Title>{m.memory_usage()}</Card.Title>
					<Card.Action class="text-muted-foreground flex flex-col gap-0.5 text-sm">
						<div class="flex justify-between gap-2">
							<p>{m.requests()}</p>
							<p class="font-mono">{Math.round(memoryRequests * 100)}%</p>
						</div>
						<div class="flex justify-between gap-2">
							<p>{m.limits()}</p>
							<p class="font-mono">{Math.round(memoryLimits * 100)}%</p>
						</div>
					</Card.Action>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={memoryUsagesConfiguration}>
						<AreaChart
							data={memoryUsages}
							x="time"
							xScale={scaleUtc()}
							yPadding={[0, 25]}
							series={[
								{
									key: 'value',
									label: memoryUsagesConfiguration.usage.label,
									color: memoryUsagesConfiguration.usage.color,
								},
							]}
							seriesLayout="stack"
							props={{
								area: {
									curve: curveNatural,
									'fill-opacity': 0.4,
									line: { class: 'stroke-1' },
									motion: 'tween',
								},
								xAxis: {
									format: (v: Date) =>
										`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`,
								},
								yAxis: { format: () => '' },
							}}
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
											minute: 'numeric',
										});
									}}
								>
									{#snippet formatter({ item, name, value })}
										{@const { value: capacity, unit } = formatCapacity(Number(value))}
										<div
											style="--color-bg: {item.color}"
											class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
										></div>
										<div
											class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
										>
											<div class="grid gap-1.5">
												<span class="text-muted-foreground">{name}</span>
											</div>
											<p class="font-mono">{capacity} {unit}</p>
										</div>
									{/snippet}
								</Chart.Tooltip>
							{/snippet}
							{#snippet marks({ series, getAreaProps })}
								{#each series as s, i (s.key)}
									<LinearGradient
										stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
										vertical
									>
										{#snippet children({ gradient })}
											<Area {...getAreaProps(s, i)} fill={gradient} />
										{/snippet}
									</LinearGradient>
								{/each}
							{/snippet}
						</AreaChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="relative col-span-2 col-start-1 gap-2 overflow-hidden">
				<Icon
					icon="ph:squares-four"
					class="text-primary/5 absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				/>
				<Card.Header>
					<Card.Title>{m.pods()}</Card.Title>
					<Card.Description>{m.running()}</Card.Description>
				</Card.Header>
				<Card.Content class="text-3xl">
					{runningPods}
				</Card.Content>
			</Card.Root>

			<Card.Root class="relative col-span-2 gap-2 overflow-hidden">
				<Icon
					icon="ph:shipping-container"
					class="text-primary/5 absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				/>
				<Card.Header>
					<Card.Title>{m.containers()}</Card.Title>
					<Card.Description>{m.running()}</Card.Description>
				</Card.Header>
				<Card.Content class="text-3xl">
					{runningContainers}
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-3 col-start-5 gap-2">
				<Card.Header>
					<Card.Title>{m.network_bandwidth()}</Card.Title>
					<Card.Description>{m.receive_and_transmit()}</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={trafficsConfigurations}>
						<AreaChart
							data={traffics}
							x="time"
							xScale={scaleUtc()}
							yPadding={[0, 25]}
							series={[
								{
									key: 'receive',
									label: trafficsConfigurations.receive.label,
									color: trafficsConfigurations.receive.color,
								},
								{
									key: 'transmit',
									label: trafficsConfigurations.transmit.label,
									color: trafficsConfigurations.transmit.color,
								},
							]}
							seriesLayout="stack"
							props={{
								area: {
									curve: curveNatural,
									'fill-opacity': 0.4,
									line: { class: 'stroke-1' },
									motion: 'tween',
								},
								xAxis: {
									format: (v: Date) =>
										`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`,
								},
								yAxis: { format: () => '' },
							}}
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
											minute: 'numeric',
										});
									}}
								>
									{#snippet formatter({ item, name, value })}
										{@const { value: io, unit } = formatIO(Number(value))}
										<div
											style="--color-bg: {item.color}"
											class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
										></div>
										<div
											class="flex flex-1 shrink-0 items-center justify-between gap-2 text-xs leading-none"
										>
											<div class="grid gap-1.5">
												<span class="text-muted-foreground">{name}</span>
											</div>
											<p class="font-mono">{io} {unit}</p>
										</div>
									{/snippet}
								</Chart.Tooltip>
							{/snippet}
							{#snippet marks({ series, getAreaProps })}
								{#each series as s, i (s.key)}
									<LinearGradient
										stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
										vertical
									>
										{#snippet children({ gradient })}
											<Area {...getAreaProps(s, i)} fill={gradient} />
										{/snippet}
									</LinearGradient>
								{/each}
							{/snippet}
						</AreaChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-3 gap-2">
				<Card.Header>
					<Card.Title>{m.storage_throughPut()}</Card.Title>
					<Card.Description>{m.read_and_write()}</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={throughputsConfigurations}>
						<AreaChart
							data={throughputs}
							x="time"
							xScale={scaleUtc()}
							yPadding={[0, 25]}
							series={[
								{
									key: 'read',
									label: throughputsConfigurations.read.label,
									color: throughputsConfigurations.read.color,
								},
								{
									key: 'write',
									label: throughputsConfigurations.write.label,
									color: throughputsConfigurations.write.color,
								},
							]}
							seriesLayout="stack"
							props={{
								area: {
									curve: curveNatural,
									'fill-opacity': 0.4,
									line: { class: 'stroke-1' },
									motion: 'tween',
								},
								xAxis: {
									format: (v: Date) =>
										`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`,
								},
								yAxis: { format: () => '' },
							}}
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
											minute: 'numeric',
										});
									}}
								>
									{#snippet formatter({ item, name, value })}
										{@const { value: io, unit } = formatIO(Number(value))}
										<div
											style="--color-bg: {item.color}"
											class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
										></div>
										<div
											class="flex flex-1 shrink-0 items-center justify-between gap-2 text-xs leading-none"
										>
											<div class="grid gap-1.5">
												<span class="text-muted-foreground">{name}</span>
											</div>
											<p class="font-mono">{io} {unit}</p>
										</div>
									{/snippet}
								</Chart.Tooltip>
							{/snippet}
							{#snippet marks({ series, getAreaProps })}
								{#each series as s, i (s.key)}
									<LinearGradient
										stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
										vertical
									>
										{#snippet children({ gradient })}
											<Area {...getAreaProps(s, i)} fill={gradient} />
										{/snippet}
									</LinearGradient>
								{/each}
							{/snippet}
						</AreaChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>
		</Tabs.Content>
		<Tabs.Content value="analytics">
			<!-- {#if isMounted && $prometheusDriver && $activeScope}
				<Dashboard client={$prometheusDriver} scope={$activeScope} />
			{:else}
				<Loading />
			{/if} -->
		</Tabs.Content>
	</Tabs.Root>
</div>
