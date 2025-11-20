<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatBigNumber, formatPercentage, formatProgressColor } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	import type { Application } from '../types';

	let {
		table,
		scope
	}: {
		table: Table<Application>;
		scope: string;
	} = $props();

	const filteredApplications = $derived(
		table.getFilteredRowModel().rows.map((row) => row.original)
	);

	// // Client setup
	// const transport: Transport = getContext('transport');
	// const client = createClient(ApplicationService, transport);
	// const applications = writable<Application[]>([]);

	// // State
	// let selectedValue = $state('');

	// // Computed values
	// const filteredApplications = $derived($applications.filter((a) => a.type === selectedValue));

	// const totalPods = $derived(
	// 	filteredApplications.reduce((total, application) => total + application.pods.length, 0)
	// );

	// const numberOfServices = $derived(
	// 	filteredApplications.reduce((total, application) => total + application.services.length, 0)
	// );

	// const healthyPods = $derived(
	// 	filteredApplications.reduce((total, application) => total + application.healthies, 0)
	// );

	// const healthByType = $derived(totalPods > 0 ? (healthyPods * 100) / totalPods : 0);

	// const healthColorClass = $derived(formatProgressColor(healthByType));

	// onMount(async () => {
	// 	try {
	// 		const response = await client.listApplications({
	// 			scope: scope
	// 		});

	// 		applications.set(
	// 			response.applications.map((application) => ({
	// 				...application,
	// 				hostname: response.hostname
	// 			}))
	// 		);

	// 		if (response.applications && response.applications[0]) {
	// 			selectedValue = response.applications[0].type;
	// 		}
	// 	} catch (error) {
	// 		console.error('Error fetching applications:', error);
	// 	}
	// });
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<!-- <Layout>
		{#snippet title()}
			<Title title="APPLICATION" />
		{/snippet}

		{#snippet content()}
			<Content value={filteredApplications.length} />
		{/snippet}
	</Layout>

	<Layout>
		{#snippet title()}
			<Title title="SERVICE" />
		{/snippet}

		{#snippet content()}
			<Content value={numberOfServices} />
		{/snippet}
	</Layout>

	<Layout>
		{#snippet title()}
			<Title title="POD" />
		{/snippet}

		{#snippet content()}
			<Content value={totalPods} />
		{/snippet}
	</Layout>

	<Layout>
		{#snippet title()}
			<Title title="HEALTH" />
		{/snippet}

		{#snippet content()}
			<ContentSubtitle
				value={Math.round(healthByType)}
				unit="%"
				subtitle={`${healthyPods} Running over ${totalPods} pods`}
			/>
		{/snippet}

		{#snippet footer()}
			<Progress value={healthByType} max={100} class={healthColorClass} />
		{/snippet}
	</Layout> -->

	{#snippet Applications()}
		{@const title = m.applications()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:cube'}
		{@const applications = filteredApplications.length}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title class="flex items-center gap-2 font-medium">
					<div
						class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
					>
						<Icon icon={titleIcon} class="size-5" />
					</div>
					<p class="font-bold">{title}</p>
					<Badge>
						{scope}
					</Badge>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<p class="text-7xl font-semibold">{applications}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Applications()}

	{#snippet Services()}
		{@const title = m.services()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:squares-four'}
		{@const services = filteredApplications.reduce(
			(total, application) => total + application.services.length,
			0
		)}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title class="flex items-center gap-2 font-medium">
					<div
						class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
					>
						<Icon icon={titleIcon} class="size-5" />
					</div>
					<p class="font-bold">{title}</p>
					<Badge>
						{scope}
					</Badge>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<p class="text-7xl font-semibold">{services}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Services()}

	{#snippet HealthPods()}
		{@const title = m.health()}
		{@const titleIcon = 'ph:chart-pie-bold'}
		{@const backgroundIcon = 'ph:check'}
		{@const healthPods = filteredApplications.reduce(
			(a, application) => a + application.healthies,
			0
		)}
		{@const totalPods = filteredApplications.reduce(
			(total, application) => total + application.pods.length,
			0
		)}
		{@const percentage = formatPercentage(healthPods, totalPods)}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title>
					<div class="flex items-center gap-2 font-medium">
						<div
							class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
						>
							<Icon icon={titleIcon} class="size-5" />
						</div>
						<p class="font-bold">{title}</p>
						<Badge>
							{scope}
						</Badge>
					</div>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<div class="space-y-1">
					<p class="text-5xl font-semibold">{percentage ? `${percentage} %` : 'NaN'}</p>
					<p class="text-3xl text-muted-foreground">
						{formatBigNumber(healthPods)}/{formatBigNumber(totalPods)}
					</p>
				</div>
			</Card.Content>
			<Progress
				value={Number(percentage ?? 0)}
				max={100}
				class={cn(
					formatProgressColor(Number(percentage ?? 0)),
					'absolute top-0 left-0 h-2 rounded-none'
				)}
			/>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render HealthPods()}
</div>
