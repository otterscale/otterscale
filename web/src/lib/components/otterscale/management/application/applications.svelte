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
	import { Nexus, type Application, type StorageClass } from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { CreateStorageClasses } from './index';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	let {
		scopeUuid,
		facilityName
	}: {
		scopeUuid: string;
		facilityName: string;
	} = $props();

	const applicationsStore = writable<Application[]>();
	const applicationsIsLoading = writable(true);
	async function fetchApplications() {
		try {
			const response = await client.listApplications({
				scopeUuid: scopeUuid,
				facilityName: facilityName
			});

			applicationsStore.set(response.applications);
		} catch (error) {
			console.error('Error fetching machine:', error);
		} finally {
			applicationsIsLoading.set(false);
		}
	}

	const storageClassesStore = writable<StorageClass[]>([]);
	const storageClassesLoading = writable(true);
	async function fetchStorageClasses() {
		try {
			const response = await client.listStorageClasses({
				scopeUuid: scopeUuid,
				facilityName: facilityName
			});
			storageClassesStore.set(response.storageClasses);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			storageClassesLoading.set(false);
		}
	}

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchApplications();
			await fetchStorageClasses();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if !mounted}
	<PageLoading />
{:else if $applicationsStore.length == 0}
	{@render NoApplicationHandler()}
{:else}
	{@const types = new Set($applicationsStore.map((a) => a.type))}
	{@const firstType = $applicationsStore[0].type}

	<div class="py-2">
		<Tabs.Root value={firstType}>
			<Tabs.List>
				{#each [...types] as type}
					<Tabs.Trigger value={type} class="flex items-start gap-1">
						{type}
						{@render Notification(type)}
					</Tabs.Trigger>
				{/each}
				<Tabs.Trigger value="storage_classes" class="flex items-start gap-1">
					Storage Classes
				</Tabs.Trigger>
			</Tabs.List>
			{#each [...types] as type}
				{@const applicationsByType = $applicationsStore.filter((a) => a.type === type)}
				{@const sortedApplicationsByType = applicationsByType.sort((previous, present) =>
					previous.namespace.localeCompare(present.namespace)
				)}
				<Tabs.Content value={type}>
					{@render Statistics(type, sortedApplicationsByType)}
					<div class="grid gap-2">
						<Table.Root>
							<Table.Header>
								<Table.Row class="*:text-xs *:font-light">
									<Table.Head>NAME</Table.Head>
									<Table.Head>NAMESPACE</Table.Head>
									<Table.Head>HEALTH</Table.Head>
									<Table.Head class="text-right">SERVICES</Table.Head>
									<Table.Head class="text-right">PODS</Table.Head>
									<Table.Head class="text-right">REPLICAS</Table.Head>
									<Table.Head class="text-right">CONTAINERS</Table.Head>
									<Table.Head class="text-right">VOLUMES</Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each sortedApplicationsByType as application}
									<Table.Row class="*:text-sm">
										<Table.Cell>
											<span class="flex items-center gap-1">
												<a
													href={`/management/scope/${scopeUuid}/facility/${facilityName}/namespace/${application.namespace}/application/${application.name}`}
												>
													{application.name}
												</a>
												<Icon icon="ph:arrow-square-out" />
											</span>
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
									</Table.Row>
								{/each}
							</Table.Body>
						</Table.Root>
					</div>
				</Tabs.Content>
			{/each}
			<Tabs.Content value="storage_classes">
				<div class="flex justify-end p-2">
					<CreateStorageClasses {scopeUuid} />
				</div>
				<div class="grid gap-2">
					<Table.Root>
						<Table.Header>
							<Table.Row class="*:text-xs *:font-light">
								<Table.Head>NAME</Table.Head>
								<Table.Head>PROVISIONER</Table.Head>
								<Table.Head>RECLAIM POLICY</Table.Head>
								<Table.Head>VOLUME BINDING MODE</Table.Head>
								<Table.Head>Configuration</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each $storageClassesStore as storageClass}
								<Table.Row class="*:text-sm">
									<Table.Cell>{storageClass.name}</Table.Cell>
									<Table.Cell>{storageClass.provisioner}</Table.Cell>
									<Table.Cell>{storageClass.reclaimPolicy}</Table.Cell>
									<Table.Cell>{storageClass.volumeBindingMode}</Table.Cell>
									<Table.Cell>
										<div class="flex justify-start">
											<HoverCard.Root>
												<HoverCard.Trigger>
													<Button variant="ghost" size="icon" class="h-6 w-6">
														<Icon icon="ph:info" class="size-4" />
													</Button>
												</HoverCard.Trigger>
												<HoverCard.Content class="w-fit">
													<Table.Root>
														<Table.Body>
															{#each Object.entries(storageClass.parameters) as [key, value]}
																<Table.Row class="border-none *:text-xs">
																	<Table.Head class="text-right font-light">{key}</Table.Head>
																	<Table.Cell class="font-base">{value}</Table.Cell>
																</Table.Row>
															{/each}
														</Table.Body>
													</Table.Root>
												</HoverCard.Content>
											</HoverCard.Root>
										</div>
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</div>
			</Tabs.Content>
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
	{@const healthByType = (numberOfHealthApplicationsByType * 100) / numberOApplicationsByType || 0}
	<div class="grid grid-cols-4 gap-3 *:border-none *:shadow-none">
		<Card.Root>
			<Card.Header>
				<Card.Title>Application</Card.Title>
			</Card.Header>
			<Card.Content class="text-7xl">
				{$applicationsStore.filter((a) => a.type === type).length}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>Health</Card.Title>
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
	<div class="py-2">
		<Alert.Root
			class="absolute z-10 flex w-[calc(100vw_-_theme(spacing.72))] justify-between bg-blue-50 opacity-90"
		>
			<span class="flex items-center gap-2">
				<Icon icon="ph:info" class="size-10" />
				<span>
					<Alert.Title>INFORMATION</Alert.Title>
					<Alert.Description
						>There are no applications, please install from Store.</Alert.Description
					>
				</span>
			</span>
			<Button variant="outline" class="text-sm">Go to Store</Button>
		</Alert.Root>
		<Tabs.Root>
			<Tabs.List class="flex justify-start gap-2 bg-transparent">
				{#each Array(3) as _}
					<Skeleton class="h-6 w-32 rounded-lg" />
				{/each}
			</Tabs.List>
			<div class="mt-4">
				<div class="mb-4 grid grid-cols-4 gap-3">
					{#each Array(3) as _}
						<Card.Root class="border-none shadow-none">
							<Card.Header>
								<Skeleton class="h-6 w-24" />
							</Card.Header>
							<Card.Content>
								<Skeleton class="h-16 w-full" />
							</Card.Content>
						</Card.Root>
					{/each}
				</div>
				<div class="grid gap-2">
					<Table.Root>
						<Table.Header>
							<Table.Row>
								{#each Array(5) as _}
									<Table.Head>
										<Skeleton class="h-4 w-20" />
									</Table.Head>
								{/each}
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each Array(9) as _}
								<Table.Row>
									{#each Array(8) as _}
										<Table.Cell>
											<Skeleton class="h-4 w-full" />
										</Table.Cell>
									{/each}
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</div>
			</div>
		</Tabs.Root>
	</div>
{/snippet}
