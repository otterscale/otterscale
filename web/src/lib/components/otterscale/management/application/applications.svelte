<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Button } from '$lib/components/ui/button';
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import {
		Nexus,
		type Application,
		type Facility_Info,
		type Scope,
		type StorageClass
	} from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { page } from '$app/state';
	import {
		ManagementKubernetesComboBox,
		ManagementScopeComboBox
	} from '$lib/components/otterscale/index';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';

	let scopeUuid = $state('');
	let facilityName = $state('');

	const transport: Transport = getContext('transport');
	const client = createClient(Nexus, transport);

	async function setDefaultScope() {
		client.listScopes({}).then((r) => {
			const defaultScopeUuid = r.scopes.find((s) => s.name === 'default')?.uuid;
			scopeUuid = page.url.searchParams.get('scope') ?? defaultScopeUuid ?? '';
		});
	}

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.listScopes({});
			scopesStore.set(response.scopes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}

	const kubernetesesStore = writable<Facility_Info[]>();
	async function fetchKuberneteses() {
		try {
			const response = await client.listKuberneteses({
				scopeUuid: scopeUuid
			});

			kubernetesesStore.set(response.kuberneteses);

			let defaultKubernetes = response.kuberneteses.find((s) => s.scopeUuid === scopeUuid);
			if (defaultKubernetes) {
				facilityName = defaultKubernetes.facilityName;
			}
		} catch (error) {
			console.error('Error fetching kubernetes:', error);
		}
	}

	const applicationsStore = writable<Application[]>([]);
	const applicationsIsLoading = writable(true);
	async function fetchApplications() {
		try {
			if (facilityName === '') {
				throw new Error('facilityName is empty');
			}
			const response = await client.listApplications({
				scopeUuid: scopeUuid,
				facilityName: facilityName
			});

			applicationsStore.set(response.applications);

			if (response.applications && response.applications[0]) {
				selectedValue = response.applications[0].type;
			}
		} catch (error) {
			applicationsStore.set([]);
			console.error('Error fetching applications:', error);
		} finally {
			applicationsIsLoading.set(false);
		}
	}

	async function refreshApplications() {
		while (page.url.searchParams.get('intervals')) {
			await new Promise((resolve) =>
				setTimeout(resolve, 1000 * Number(page.url.searchParams.get('intervals')))
			);
			console.log(`Refresh applications`);

			try {
				const response = await client.listApplications({
					scopeUuid: scopeUuid,
					facilityName: facilityName
				});

				applicationsStore.set(response.applications);
			} catch (error) {
				console.error('Error fetching applications:', error);
			}
		}
	}

	let selectedValue = $state('');

	let mounted = $state(false);
	onMount(async () => {
		try {
			await setDefaultScope();
			await fetchScopes()
				.then(async () => await fetchKuberneteses())
				.then(async () => await fetchApplications());
			refreshApplications();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});

	async function handleChange() {
		await fetchKuberneteses().then(async () => await fetchApplications());
	}
</script>

{#if !mounted}
	<PageLoading />
	<!-- {:else if facilityName === ''} -->
	<!-- {@render NoApplicationHandler()} -->
{:else}
	{@const applicationsByType = $applicationsStore.filter((a) => a.type === selectedValue)}
	{@const sortedApplicationsByType = applicationsByType.sort((previous, present) =>
		previous.namespace.localeCompare(present.namespace)
	)}
	{@render Statistics(selectedValue, sortedApplicationsByType)}

	{@const types = new Set($applicationsStore.map((a) => a.type))}
	<div class="py-2">
		<Tabs.Root bind:value={selectedValue}>
			<div class="flex justify-between">
				<div>
					{#if $applicationsStore.length > 0}
						<Tabs.List>
							{#each [...types] as type}
								<Tabs.Trigger value={type} class="flex items-start gap-1">
									{type}
									{@render Notification(type)}
								</Tabs.Trigger>
							{/each}
						</Tabs.List>
					{/if}
				</div>
				<ManagementScopeComboBox
					scopes={$scopesStore}
					bind:selected={scopeUuid}
					onSelect={handleChange}
				/>
			</div>
			{#if $applicationsStore.length === 0}
				{@render NoApplicationHandler()}
			{/if}
			{#each [...types] as type}
				{@const applicationsByType = $applicationsStore.filter((a) => a.type === type)}
				{@const sortedApplicationsByType = applicationsByType.sort((previous, present) =>
					previous.namespace.localeCompare(present.namespace)
				)}
				<Tabs.Content value={type}>
					<div class="grid gap-2 p-4">
						<Table.Root>
							<Table.Header class="bg-muted/50">
								<Table.Row class="*:text-xs *:font-light">
									<Table.Head>NAME</Table.Head>
									<Table.Head>NAMESPACE</Table.Head>
									<Table.Head>HEALTH</Table.Head>
									<Table.Head class="text-right">SERVICES</Table.Head>
									<Table.Head class="text-right">PODS</Table.Head>
									<Table.Head class="text-right">REPLICAS</Table.Head>
									<Table.Head class="text-right">CONTAINERS</Table.Head>
									<Table.Head class="text-right">VOLUMES</Table.Head>
									<Table.Head>NODEPORT</Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each sortedApplicationsByType as application}
									<Table.Row class="*:text-sm">
										<Table.Cell>
											<a
												href={`/management/facility/${facilityName}/namespace/${application.namespace}/application/${application.name}?scope=${scopeUuid}&intervals=15`}
											>
												<span class="flex items-center gap-1">
													{application.name}
													<Icon icon="ph:arrow-square-out" />
												</span>
											</a>
										</Table.Cell>
										<Table.Cell>
											<Badge variant="outline">{application.namespace}</Badge>
										</Table.Cell>

										<Table.Cell>
											{#if application.pods.length > 0}
												{@const healthByApplication = Math.round(
													(application.healthies / application.pods.length) * 100 || 0
												)}
												<div class="flex flex-col justify-end gap-1">
													<p class="text-right">{healthByApplication}%</p>

													<Progress
														value={healthByApplication}
														max={100}
														class={`${
															healthByApplication > 62
																? 'bg-green-50 *:bg-green-700'
																: healthByApplication > 38
																	? 'bg-yellow-50 *:bg-yellow-500'
																	: 'bg-red-50 *:bg-red-700'
														}`}
													/>
												</div>
											{/if}
										</Table.Cell>
										<Table.Cell class="text-right">{application.services.length}</Table.Cell>
										<Table.Cell class="text-right">{application.pods.length}</Table.Cell>
										<Table.Cell class="text-right">{application.replicas}</Table.Cell>
										<Table.Cell class="text-right">{application.containers.length}</Table.Cell>
										<Table.Cell class="text-right"
											>{application.persistentVolumeClaims.length}</Table.Cell
										>
										<Table.Cell>
											{#each application.services as service}
												{#each service.ports as port}
													{#if port.nodePort > 0}
														<Button
															variant="ghost"
															target="_blank"
															href={`http://${application.publicAddress}:${port.nodePort}`}
														>
															{port.targetPort}
															<Icon icon="ph:arrow-square-out" />
														</Button>
													{/if}
												{/each}
											{/each}
										</Table.Cell>
									</Table.Row>
								{/each}
							</Table.Body>
						</Table.Root>
					</div>
				</Tabs.Content>
			{/each}
		</Tabs.Root>
	</div>
{/if}

{#snippet Notification(type: string)}
	{@const applicationsByType = $applicationsStore.filter((a) => a.type === type)}
	{@const numberOfFailPods = applicationsByType.reduce(
		(a, application) => a + (application.pods.length - application.healthies),
		0
	)}
	{#if numberOfFailPods > 0}
		<Badge
			variant="destructive"
			class="-mt-1 flex h-4 w-4 animate-pulse justify-center rounded-full p-0 text-[11px] font-light"
		>
			{numberOfFailPods}
		</Badge>
	{/if}
{/snippet}

{#snippet Statistics(type: string, applicationsByType: Application[])}
	{@const numberOApplicationsByType = applicationsByType.reduce(
		(a, application) => a + application.pods.length,
		0
	)}
	{@const numberOfHealthApplicationsByType = applicationsByType.reduce(
		(a, application) => a + application.healthies,
		0
	)}
	{@const numberOfPodsByType = applicationsByType.reduce(
		(a, application) => a + application.pods.length,
		0
	)}
	{@const numberOfServicesByType = applicationsByType.reduce(
		(a, application) => a + application.services.length,
		0
	)}
	{@const healthByType = (numberOfHealthApplicationsByType * 100) / numberOApplicationsByType || 0}
	<div class="grid grid-cols-4 gap-3">
		<Card.Root>
			<Card.Header>
				<Card.Title>APPLICATION</Card.Title>
			</Card.Header>
			<Card.Content class="text-7xl">
				{$applicationsStore.filter((a) => a.type === type).length}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>SERVICE</Card.Title>
			</Card.Header>
			<Card.Content class="text-7xl">
				{numberOfServicesByType}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>POD</Card.Title>
			</Card.Header>
			<Card.Content class="text-7xl">
				{numberOfPodsByType}
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
				<p class="text-xs text-muted-foreground">
					{numberOfHealthApplicationsByType} Running over {numberOApplicationsByType} pods
				</p>
			</Card.Content>
			<Card.Footer>
				<Progress
					value={healthByType}
					max={100}
					class={`${
						healthByType > 62
							? 'bg-green-50 *:bg-green-700'
							: healthByType > 38
								? 'bg-yellow-50 *:bg-yellow-500'
								: 'bg-red-50 *:bg-red-700'
					}`}
				/>
			</Card.Footer>
		</Card.Root>
	</div>
{/snippet}

{#snippet NoApplicationHandler()}
	<div class="py-4">
		<Alert.Root class="flex justify-between bg-blue-50 opacity-90">
			<span class="flex items-center gap-2">
				<Icon icon="ph:info" class="size-10" />
				<span>
					<Alert.Title>INFORMATION</Alert.Title>
					<Alert.Description
						>There are no applications, please install from Store.</Alert.Description
					>
				</span>
			</span>
			<Button
				variant="outline"
				class="text-sm"
				href="/market"
				onclick={() =>
					toast.info(`Back to Application`, {
						duration: Number.POSITIVE_INFINITY,
						action: {
							label: 'Go',
							onClick: () => goto(`/management/application?scope=${scopeUuid}&intervals=30`)
						}
					})}>Go to Store</Button
			>
		</Alert.Root>
	</div>
{/snippet}
