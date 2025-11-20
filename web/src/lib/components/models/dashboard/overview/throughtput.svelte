<script lang="ts">
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
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

	let prompts = $state([] as SampleValue[]);
	let generations = $state([] as SampleValue[]);
	const throughputs = $derived(
		prompts.map((sample, index) => ({
			time: sample.time,
			prompt: sample.value,
			generation: generations[index]?.value ?? 0
		}))
	);

	const configuration = {
		prompt: { label: 'Prompt', color: 'var(--chart-1)' },
		generation: { label: 'Generation', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	async function fetchPrompts() {
		try {
			const response = await prometheusDriver.instantQuery(
				`max(rate(vllm:prompt_tokens_total{juju_model="${scope}"}[2m]))`
			);
			prompts = response.result[0]?.values ?? [];
		} catch (error) {
			console.error(`Fail to fetch latest prompt throughput in scope ${scope}:`, error);
		}
	}

	async function fetchGenerations() {
		try {
			const response = await prometheusDriver.instantQuery(
				`max(rate(vllm:generation_tokens_total{juju_model="${scope}"}[2m]))`
			);
			generations = response.result[0]?.values ?? [];
		} catch (error) {
			console.error(`Fail to fetch latest generation throughput in scope ${scope}:`, error);
		}
	}

	async function fetch() {
		try {
			await Promise.all([fetchPrompts(), fetchGenerations()]);
		} catch (error) {
			console.error(`Fail to fetch throughputs data in scope ${scope}:`, error);
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
			<Card.Title>{m.throughput()}</Card.Title>
			<Card.Description>{m.llm_dashboard_throughputs_tooltip()}</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={configuration} class="h-[200px] w-full">
				<AreaChart
					data={throughputs}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'prompt',
							label: configuration.prompt.label,
							color: configuration.prompt.color
						},
						{
							key: 'generation',
							label: configuration.generation.label,
							color: configuration.generation.color
						}
					]}
					props={{
						area: {
							curve: curveNatural,
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
						>
							{#snippet formatter({ item, name, value })}
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
									<p class="font-mono">{Number(value).toFixed(2)} {m.per_second()}</p>
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
{/if}
