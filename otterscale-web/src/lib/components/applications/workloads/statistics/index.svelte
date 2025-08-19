<script lang="ts">
	import { ApplicationService, type Application } from '$lib/api/application/v1/application_pb';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatHealthColor } from '$lib/formatter';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	let { scopeUuid, facilityName }: { scopeUuid: string; facilityName: string } = $props();

	// Client setup
	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);
	const applications = writable<Application[]>([]);

	// State
	let selectedValue = $state('');
	let isMounted = $state(false);

	// Computed values
	const filteredApplications = $derived($applications.filter((a) => a.type === selectedValue));

	const totalPods = $derived(
		filteredApplications.reduce((total, application) => total + application.pods.length, 0)
	);

	const numberOfServices = $derived(
		filteredApplications.reduce((total, application) => total + application.services.length, 0)
	);

	const healthyPods = $derived(
		filteredApplications.reduce((total, application) => total + application.healthies, 0)
	);

	const healthByType = $derived(totalPods > 0 ? (healthyPods * 100) / totalPods : 0);

	const healthColorClass = $derived(formatHealthColor(healthByType));

	onMount(async () => {
		try {
			const response = await client.listApplications({
				scopeUuid: scopeUuid,
				facilityName: facilityName
			});

			applications.set(response.applications);

			if (response.applications && response.applications[0]) {
				selectedValue = response.applications[0].type;
			}

			isMounted = true;
		} catch (error) {
			console.error('Error fetching applications:', error);
		}
	});
</script>

<span class="grid grid-cols-4 gap-4">
	<Card.Root>
		<Card.Header>
			<Card.Title>APPLICATION</Card.Title>
		</Card.Header>
		<Card.Content class="text-7xl">
			{filteredApplications.length}
		</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header>
			<Card.Title>SERVICE</Card.Title>
		</Card.Header>
		<Card.Content class="text-7xl">
			{numberOfServices}
		</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header>
			<Card.Title>POD</Card.Title>
		</Card.Header>
		<Card.Content class="text-7xl">
			{totalPods}
		</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header>
			<Card.Title>HEALTH</Card.Title>
		</Card.Header>
		<Card.Content>
			<p class="text-3xl">
				{Math.round(healthByType)}%
			</p>
			<p class="text-muted-foreground text-xs">
				{healthyPods} Running over {totalPods} pods
			</p>
		</Card.Content>
		<Card.Footer>
			<Progress value={healthByType} max={100} class={healthColorClass} />
		</Card.Footer>
	</Card.Root>
</span>
