<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import type { Application } from '../types';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import ContentSubtitle from '$lib/components/custom/chart/content/text/text-with-subtitle.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatProgressColor } from '$lib/formatter';

	let { scope, facility }: { scope: string; facility: string } = $props();

	// Client setup
	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);
	const applications = writable<Application[]>([]);

	// State
	let selectedValue = $state('');

	// Computed values
	const filteredApplications = $derived($applications.filter((a) => a.type === selectedValue));

	const totalPods = $derived(filteredApplications.reduce((total, application) => total + application.pods.length, 0));

	const numberOfServices = $derived(
		filteredApplications.reduce((total, application) => total + application.services.length, 0),
	);

	const healthyPods = $derived(filteredApplications.reduce((total, application) => total + application.healthies, 0));

	const healthByType = $derived(totalPods > 0 ? (healthyPods * 100) / totalPods : 0);

	const healthColorClass = $derived(formatProgressColor(healthByType));

	onMount(async () => {
		try {
			const response = await client.listApplications({
				scope: scope,
				facility: facility,
			});

			applications.set(
				response.applications.map((application) => ({
					...application,
					publicAddress: response.publicAddress,
				})),
			);

			if (response.applications && response.applications[0]) {
				selectedValue = response.applications[0].type;
			}
		} catch (error) {
			console.error('Error fetching applications:', error);
		}
	});
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<Layout>
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
	</Layout>
</div>
