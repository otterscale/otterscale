<script lang="ts">
	const healthData = group(healthRawData, (d) => d.type).get('error') ?? [];
	function metricColor(metric: number) {
		switch (true) {
			case metric > 62:
				return 'fill-red-800';
			case metric > 38:
				return 'fill-yellow-500';
			default:
				return 'fill-green-800';
		}
	}
	function metricBackgroundColor(metric: number) {
		switch (true) {
			case metric > 62:
				return 'fill-red-50';
			case metric > 38:
				return 'fill-yellow-50';
			default:
				return 'fill-green-50';
		}
	}

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let data = $state(healthData);

	import Autoplay from 'embla-carousel-autoplay';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Carousel from '$lib/components/ui/carousel/index.js';
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { cn } from '$lib/utils';
	import { healthRawData } from './dataset';
	import { group } from 'd3-array';
	import { AreaChart, Svg, Group, Arc, Chart, Text } from 'layerchart';
	import { Badge } from '$lib/components/ui/badge';
	import { latencies, storages, inputOutputs, usages, currentUsage } from './dataset';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Button } from '$lib/components/ui/button';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { writable } from 'svelte/store';
	import LLMError from './error/llm.svelte';
	import PrometheusError from './error/prometheus.svelte';
	import KubernetesError from './error/kubernetes.svelte';
	import CephError from './error/ceph.svelte';
	import Icon from '@iconify/svelte';
	import { Nexus, type Scope, type Error, type Application } from '$gen/api/nexus/v1/nexus_pb';
	import { goto } from '$app/navigation';
	import { MessageIterator } from '$lib/components/otterscale/ui/index';

	let { scope }: { scope: Scope } = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const errorsStore = writable<Error[]>([]);
	const errorsLoading = writable(true);
	async function fetchErrors(scopeUuid: string) {
		try {
			const response = await client.verifyEnvironment({ scopeUuid: scopeUuid });
			errorsStore.set(response.errors);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			errorsLoading.set(false);
		}
	}
	async function refreshErrors(scopeUuid: string) {
		while (page.url.searchParams.get('intervals')) {
			await new Promise((resolve) =>
				setTimeout(resolve, 1000 * Number(page.url.searchParams.get('intervals')))
			);

			try {
				const response = await client.verifyEnvironment({ scopeUuid: scopeUuid });
				errorsStore.set(response.errors);
			} catch (error) {
				console.error('Error fetching:', error);
			}
		}
	}

	const applications = [] as Application[];
	async function fetchApplicationsOverKuberneteses(scopeUuid: string) {
		try {
			const responseKubernetes = await client.listKuberneteses({
				scopeUuid: scopeUuid
			});
			for (const kubernetes of responseKubernetes.kuberneteses) {
				try {
					const responseApplications = await client.listApplications({
						scopeUuid: kubernetes.scopeUuid,
						facilityName: kubernetes.facilityName
					});
					responseApplications.applications.flat().forEach((application) => {
						applications.push(application);
					});
				} catch (error) {
					console.error(
						`Error fetching applications for kubernetes ${kubernetes.facilityName}:`,
						error
					);
				}
			}
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	function isCephError(error: Error) {
		return error.code === 'CEPH_NOT_FOUND';
	}
	function isKubernetesError(error: Error) {
		return error.code === 'KUBERNETES_NOT_FOUND';
	}
	function isPrometheusError(error: Error) {
		return error.code === 'PROMETHEUS_NOT_FOUND';
	}
	function isMachineError(error: Error) {
		return error.code === 'NO_MACHINES_DEPLOYED';
	}

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchErrors(scope.uuid);
			await fetchApplicationsOverKuberneteses(scope.uuid);

			mounted = true;
			refreshErrors(scope.uuid);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !mounted}
	<PageLoading />
{:else}
	{#if $errorsStore.some((e) => isMachineError(e))}
		<Alert.Root class="flex justify-between bg-blue-50 opacity-90">
			<span class="flex items-center gap-2">
				<Icon icon="ph:info" class="size-10" />
				<span>
					<Alert.Title>INFORMATION</Alert.Title>
					<Alert.Description>
						There are no deployed machines, please add from 'Machines'.
					</Alert.Description>
				</span>
			</span>
			<Button variant="outline" class="text-sm" onclick={() => goto('/management/machine')}
				>Go to Machines</Button
			>
		</Alert.Root>
	{/if}
	{#if $errorsStore && $errorsStore.length > 0}
		{@const criticalErrors = $errorsStore.filter((e) => Number(e.level) > 3)}
		{@const level3Errors = $errorsStore.filter((e) => Number(e.level) === 3)}
		{@const level2Errors = $errorsStore.filter((e) => Number(e.level) === 2)}
		{@const level1Errors = $errorsStore.filter((e) => Number(e.level) === 1)}
		{#each criticalErrors as error}
			{#if !isCephError(error) && !isKubernetesError(error) && !isMachineError(error)}
				<Alert.Root variant="destructive">
					<Icon icon="material-symbols:warning-rounded" class="size-7" />
					<Alert.Title class="text-sm">{error.message}</Alert.Title>
					<Alert.Description class="text-xs text-destructive">{error.details}</Alert.Description>
				</Alert.Root>
			{/if}
		{/each}
		{#if level3Errors && level3Errors.length > 0}
			<MessageIterator data={level3Errors} duration={2000} />
		{/if}
		{#if level2Errors && level2Errors.length > 0}
			<MessageIterator data={level2Errors} duration={2000} />
		{/if}
		{#if level1Errors && level1Errors.length > 0}
			<MessageIterator data={level1Errors} duration={2000} />
		{/if}
	{/if}
	{#if $errorsStore.some((e) => isCephError(e))}
		<div class="flex h-full items-center justify-center">
			<Carousel.Root
				plugins={[
					Autoplay({
						delay: 3000
					})
				]}
				opts={{
					align: 'start'
				}}
				class="w-full max-w-6xl"
			>
				<Carousel.Content>
					<Carousel.Item class="md:basis-1/2 lg:basis-1/3">
						<div class="p-1">
							<CephError />
						</div>
					</Carousel.Item>
					<Carousel.Item class="md:basis-1/2 lg:basis-1/3">
						<div class="p-1">
							<KubernetesError />
						</div>
					</Carousel.Item>
					{#each Array(5) as _, i (i)}
						<Carousel.Item class="md:basis-1/2 lg:basis-1/3">
							<div class="p-1">
								<LLMError />
							</div>
						</Carousel.Item>
					{/each}
				</Carousel.Content>
				<Carousel.Previous />
				<Carousel.Next />
			</Carousel.Root>
		</div>
	{:else if $errorsStore.some((e) => isPrometheusError(e))}
		<PrometheusError />
	{:else}
		{@const numberOfUnhealthPods = applications.reduce(
			(a, application) =>
				a + (application.pods.filter((pod) => pod.phase !== 'Running').length || 0),
			0
		)}
		{@const numberOfPods = applications.reduce(
			(a, application) => a + (application.pods.length || 0),
			0
		)}
		<Tabs.Root value="overview" class="w-full">
			<Tabs.List class="w-fit">
				<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
				<Tabs.Trigger value="details">Details</Tabs.Trigger>
			</Tabs.List>

			<Tabs.Content value="overview">
				<div class="col-span-1">
					<div class="col-span-1 grid grid-cols-2 gap-2">
						<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
							<Card.Header>
								<Card.Title class="flex justify-between">
									<h1 class="text-3xl">Health</h1>
									{#if numberOfUnhealthPods > 0}
										<DropdownMenu.Root>
											<DropdownMenu.Trigger>
												<Button variant="ghost" class="flex items-start gap-2">
													<p>Alert</p>
													<Badge
														variant="destructive"
														class="flex h-4 w-4 justify-center rounded-full p-0 text-[12px] font-light"
													>
														{numberOfUnhealthPods}
													</Badge>
												</Button>
											</DropdownMenu.Trigger>
											<DropdownMenu.Content>
												{#each applications.filter( (a) => a.pods?.some((pod) => pod.phase !== 'Running') ) as application}
													<DropdownMenu.Item>
														<a
															href={`management/application?scope=${scope.uuid}`}
															target="_blank"
															class="flex w-full items-start gap-2"
														>
															<Badge
																variant="destructive"
																class="flex h-4 w-4 justify-center rounded-full p-0 text-[12px] font-light"
															>
																{application.pods?.filter((pod) => pod.phase !== 'Running').length}
															</Badge>
															{scope.name}/{application.namespace}
														</a>
													</DropdownMenu.Item>
												{/each}
											</DropdownMenu.Content>
										</DropdownMenu.Root>
									{/if}
								</Card.Title>
								<Card.Description>
									<span class="text-sm font-light">Ratio of non-running pods over all pods</span>
								</Card.Description>
							</Card.Header>
							<Card.Content>
								<div class="flex h-full w-full items-center justify-center">
									<div class={cn(`h-[260px] w-[260px]`)}>
										<Chart>
											<Svg center>
												<Group y={150 / 4}>
													<Arc
														value={numberOfPods ? (numberOfUnhealthPods * 100) / numberOfPods : 0}
														domain={[0, 100]}
														outerRadius={150}
														innerRadius={-13}
														cornerRadius={13}
														range={[-120, 120]}
														class="fill-red-800"
														track={{ class: 'fill-red-50' }}
														let:value
													>
														<Text
															value={`${Math.round(value)}%`}
															textAnchor="middle"
															verticalAnchor="middle"
															class="text-5xl tabular-nums"
														/>
													</Arc>
												</Group>
											</Svg>
										</Chart>
									</div>
								</div>
							</Card.Content>
							<Card.Footer></Card.Footer>
						</Card.Root>
						<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
							<Card.Header>
								<Card.Title class="text-3xl">Storage</Card.Title>
								<Card.Description>
									<span class="text-sm font-light">Storage usage across all clusters</span>
								</Card.Description>
							</Card.Header>
							<Card.Content>
								<div class="flex h-full w-full items-center justify-center">
									<div class={cn(`h-[260px] w-[260px]`)}>
										<Chart>
											<Svg center>
												<Group y={150 / 4}>
													<Arc
														value={Math.round((storages.used * 100) / storages.capacity)}
														domain={[0, 100]}
														outerRadius={150}
														innerRadius={-13}
														cornerRadius={13}
														range={[-120, 120]}
														class={metricColor((storages.used * 100) / storages.capacity)}
														track={{
															class: metricBackgroundColor(
																(storages.used * 100) / storages.capacity
															)
														}}
														let:value
													>
														<Text
															value={`${value}%`}
															textAnchor="middle"
															verticalAnchor="middle"
															class="text-5xl tabular-nums"
														/>
														<Text
															value={`${storages.used} ${storages.unit} used over ${storages.capacity} ${storages.unit}`}
															textAnchor="middle"
															verticalAnchor="middle"
															dy={23 + 13}
															class="text-xs font-light"
														/>
													</Arc>
												</Group>
											</Svg>
										</Chart>
									</div>
								</div>
							</Card.Content>
						</Card.Root>
					</div>

					<div class="col-span-1 -mt-4 grid grid-cols-1 gap-2">
						<div class="col-span-1 grid grid-cols-3 gap-3">
							<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
								<Card.Header>
									<Card.Title class="text-2xl">CPU</Card.Title>
									<Card.Description
										><span class="text-sm font-light">CPU usage across all clusters</span
										></Card.Description
									>
								</Card.Header>
								<Card.Content>
									<div class="flex h-full w-full items-center justify-center">
										<div class={cn(`h-[173px] w-[173px]`)}>
											<Chart>
												<Svg center>
													<Group y={100 / 4}>
														<Arc
															value={currentUsage.CPU}
															domain={[0, 100]}
															outerRadius={100}
															innerRadius={-13}
															cornerRadius={13}
															range={[-120, 120]}
															class={metricColor(currentUsage.CPU)}
															track={{
																class: metricBackgroundColor(currentUsage.CPU)
															}}
														>
															<Text
																value={`${currentUsage.CPU}%`}
																textAnchor="middle"
																verticalAnchor="middle"
																class="text-4xl tabular-nums"
															/>
														</Arc>
													</Group>
												</Svg>
											</Chart>
										</div>
									</div>
								</Card.Content>
							</Card.Root>
							<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
								<Card.Header>
									<Card.Title class="text-2xl">GPU</Card.Title>
									<Card.Description
										><span class="text-sm font-light">GPU usage across all clusters</span
										></Card.Description
									>
								</Card.Header>
								<Card.Content>
									<div class="flex h-full w-full items-center justify-center">
										<div class={cn(`h-[173px] w-[173px]`)}>
											<Chart>
												<Svg center>
													<Group y={100 / 4}>
														<Arc
															value={currentUsage.GPU}
															domain={[0, 100]}
															outerRadius={100}
															innerRadius={-13}
															cornerRadius={13}
															range={[-120, 120]}
															class={metricColor(currentUsage.GPU)}
															track={{
																class: metricBackgroundColor(currentUsage.GPU)
															}}
														>
															<Text
																value={`${currentUsage.GPU}%`}
																textAnchor="middle"
																verticalAnchor="middle"
																class="text-4xl tabular-nums"
															/>
														</Arc>
													</Group>
												</Svg>
											</Chart>
										</div>
									</div>
								</Card.Content>
							</Card.Root>
							<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
								<Card.Header>
									<Card.Title class="text-2xl">Memory</Card.Title>
									<Card.Description
										><span class="text-sm font-light">Memory usage over all cluster</span
										></Card.Description
									>
								</Card.Header>
								<Card.Content>
									<div class="flex h-full w-full items-center justify-center">
										<div class={cn(`h-[173px] w-[173px]`)}>
											<Chart>
												<Svg center>
													<Group y={100 / 4}>
														<Arc
															value={(currentUsage.memory.usage * 100) /
																currentUsage.memory.capacity}
															domain={[0, 100]}
															outerRadius={100}
															innerRadius={-13}
															cornerRadius={13}
															range={[-120, 120]}
															class={metricColor(
																(currentUsage.memory.usage * 100) / currentUsage.memory.capacity
															)}
															track={{
																class: metricBackgroundColor(
																	(currentUsage.memory.usage * 100) / currentUsage.memory.capacity
																)
															}}
														>
															<Text
																value={`${Math.round((currentUsage.memory.usage * 100) / currentUsage.memory.capacity)}%`}
																textAnchor="middle"
																verticalAnchor="middle"
																class="text-4xl tabular-nums"
															/>
														</Arc>
													</Group>
												</Svg>
											</Chart>
										</div>
									</div>
								</Card.Content>
							</Card.Root>
						</div>
					</div>
				</div>
			</Tabs.Content>

			<Tabs.Content value="details">
				<div class="col-span-1 grid grid-cols-2 gap-2">
					<div class="col-span-2">
						<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
							<Card.Header>
								<Card.Title class="text-3xl">Latency</Card.Title>
								<Card.Description>
									<span class="text-sm font-light">Latency trends for top 3 applications</span>
								</Card.Description>
							</Card.Header>
							<Card.Content>
								<div class="h-40 w-full resize overflow-visible">
									<AreaChart
										data={latencies}
										x="date"
										series={Object.keys(latencies[0])
											.filter((key) => key !== 'date')
											.map((k) => ({ key: k }))}
										legend
										props={{
											tooltip: {
												root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
												header: { class: 'font-light' }
											}
										}}
										{renderContext}
										{debug}
									/>
								</div>
							</Card.Content>
						</Card.Root>
					</div>
					<div class="col-span-2 grid grid-cols-2 gap-2">
						<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
							<Card.Header>
								<Card.Title class="text-3xl">I/O</Card.Title>
								<Card.Description>
									<span class="text-sm font-light">Input/Output operations over time</span>
								</Card.Description>
							</Card.Header>
							<Card.Content>
								<div class="h-40 w-full resize overflow-visible">
									<AreaChart
										data={inputOutputs}
										x="date"
										series={[
											{ key: 'write', color: 'hsl(var(--color-info))' },
											{ key: 'read', color: 'hsl(var(--color-danger))' }
										]}
										legend
										props={{
											tooltip: {
												root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
												header: { class: 'font-light' }
											}
										}}
										{renderContext}
										{debug}
									/>
								</div>
							</Card.Content>
						</Card.Root>
						<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
							<Card.Header>
								<Card.Title class="text-3xl">Resource</Card.Title>
								<Card.Description>
									<span class="text-sm font-light">Historical resource usage trends</span>
								</Card.Description>
							</Card.Header>
							<Card.Content>
								<div class="h-40 w-full resize overflow-visible">
									<AreaChart
										data={usages}
										x="date"
										series={[
											{ key: 'CPU', color: 'hsl(var(--color-info))' },
											{ key: 'GPU', color: 'hsl(var(--color-accent))' },
											{ key: 'Memory', color: 'hsl(var(--color-danger))' }
										]}
										legend
										props={{
											tooltip: {
												root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
												header: { class: 'font-light' }
											}
										}}
										{renderContext}
										{debug}
									/>
								</div>
							</Card.Content>
						</Card.Root>
					</div>
				</div>
			</Tabs.Content>
		</Tabs.Root>
	{/if}
{/if}
