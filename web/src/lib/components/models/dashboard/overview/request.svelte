<script lang="ts">
	import { scaleUtc } from 'd3-scale';
	import { curveStep } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let runnings = $state([] as SampleValue[]);
	let waitings = $state([] as SampleValue[]);
	const requests = $derived(
		runnings.map((sample, index) => ({
			time: sample.time,
			running: sample.value,
			waiting: waitings[index]?.value ?? 0
		}))
	);

	const configuration = {
		running: { label: 'Running', color: 'var(--chart-1)' },
		waiting: { label: 'Waiting', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	async function fetchRunnings() {
		try {
			const response = await prometheusDriver.rangeQuery(
				`sum(vllm:num_requests_running{juju_model="${scope}"})`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60
			);
			runnings = response.result[0]?.values ?? [];
		} catch (error) {
			console.error(`Fail to fetch latest running requests in scope ${scope}:`, error);
		}
	}

	async function fetchWaitings() {
		try {
			const response = await prometheusDriver.rangeQuery(
				`sum(vllm:num_requests_waiting{juju_model="${scope}"})`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60
			);
			waitings = response.result[0]?.values ?? [];
		} catch (error) {
			console.error(`Fail to fetch latest waiting requests in scope ${scope}:`, error);
		}
	}

	async function fetch() {
		try {
			await Promise.all([fetchRunnings(), fetchWaitings()]);
		} catch (error) {
			console.error(`Fail to fetch requests data in scope ${scope}:`, error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoaded = $state(false);
	onMount(async () => {
		try {
			await fetch();
			isLoaded = true;
		} catch (error) {
			console.error(`Fail to fetch data in scope ${scope}:`, error);
		}
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
</script>

{#if !isLoaded}
	Loading
{:else}
	<Card.Root class="h-full">
		<Card.Header>
			<Card.Title>{m.requests()}</Card.Title>
			<Card.Description>
				{m.llm_dashboard_requests_tooltip()}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={configuration} class="h-[200px] w-full">
				<AreaChart
					data={requests}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'running',
							label: configuration.running.label,
							color: configuration.running.color
						},
						{
							key: 'waiting',
							label: configuration.waiting.label,
							color: configuration.waiting.color
						}
					]}
					props={{
						area: {
							curve: curveStep,
							'fill-opacity': 0.4,
							line: { class: 'stroke-1' },
							motion: 'tween'
						},
						xAxis: {
							format: (v: Date) =>
								`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`
						},
						yAxis: { format: () => '' }
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
									minute: 'numeric'
								});
							}}
						/>
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
{/if}
