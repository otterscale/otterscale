<script lang="ts">
	import CreateScope from './create.svelte';

	import { Progress } from '$lib/components/ui/progress/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { formatTimeAgo } from '$lib/formatter';
	import { toast } from 'svelte-sonner';
	import { getContext, onMount } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Nexus, type CreateScopeRequest, type Scope } from '$gen/api/nexus/v1/nexus_pb';
	import * as Table from '$lib/components/ui/table';

	let createScopeRequest = $state({} as CreateScopeRequest);
	function resetCreateScopeRequest() {
		createScopeRequest = {} as CreateScopeRequest;
	}
	let createScopeConfirm = $state(false);

	let { scopes }: { scopes: Scope[] } = $props();

	let isScopeConfigurationOpen = $state(
		Object.fromEntries(scopes.map((scope) => [scope.uuid, false]))
	);

	const transport: Transport = getContext('transportNEW');

	const client = createClient(Nexus, transport);
</script>

<main>
	<div class="grid grid-cols-4 gap-3 *:border-none *:shadow-none">
		{@render StatisticScopes(scopes)}
	</div>
	<div class="p-4">
		<span class="flex justify-end py-2">
			<CreateScope />
		</span>
		<Table.Root>
			<Table.Header class="bg-muted/50">
				<Table.Row class="*:text-xs *:font-light">
					<Table.Head>NAME</Table.Head>
					<Table.Head>STATUS</Table.Head>
					<Table.Head>LIFE</Table.Head>
					<Table.Head class="text-end">UNITS</Table.Head>
					<Table.Head class="text-end">CORES</Table.Head>
					<Table.Head class="text-end">MACHINES</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each scopes as scope}
					<Table.Row>
						<Table.Cell>
							<div class="flex justify-between">
								<span class="flex items-center gap-1">
									{scope.name}
									<a href={`/management/scope/${scope.uuid}`}>
										<Icon icon="ph:arrow-square-out" />
									</a>
								</span>
								<!-- {@render GetModelDetail(scope)} -->
							</div>
						</Table.Cell>

						<Table.Cell>
							<Badge variant="outline">
								{scope.status}
							</Badge>
						</Table.Cell>
						<Table.Cell>
							<Badge variant="outline">
								{scope.life}
							</Badge>
						</Table.Cell>
						<Table.Cell class="text-end">{scope.unitCount}</Table.Cell>
						<Table.Cell class="text-end">{scope.coreCount}</Table.Cell>
						<Table.Cell class="text-end">{scope.machineCount}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
</main>
<!-- {#snippet GetModelDetail(scope: Scope)}
	<DropdownMenu.Root>
		<DropdownMenu.Trigger>
			<Icon icon="ph:dots-three-vertical" class="hover:scale-105" />
		</DropdownMenu.Trigger>
		<DropdownMenu.Content>
			<DropdownMenu.Item
				onclick={() => {
					isScopeConfigurationOpen[scope.uuid] = true;
				}}
			>
				<Icon icon="ph:info" /> Configuration
			</DropdownMenu.Item>
			<DropdownMenu.Item
				onclick={() => {
					isModelIntegrationOpen[model.uuid] = true;
				}}
			>
				<Icon icon="ph:info" /> Integration
			</DropdownMenu.Item>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
	<ManagementScopeConfiguration
		model_uuid={scope.uuid}
		bind:open={isScopeConfigurationOpen[scope.uuid]}
	/>
	<ManagementModelIntegrations
		model_uuid={model.uuid}
		bind:open={isModelIntegrationOpen[model.uuid]}
	/>
{/snippet} -->

{#snippet StatisticScopes(scopes: Scope[])}
	{@const numberOfScopes = scopes.length}
	{@const numberOfAvailableFacilities = scopes.filter((m) => m.status === 'available').length}
	{@const availability = (numberOfAvailableFacilities * 100) / numberOfScopes || 0}
	{@const numberOfAliveFacilities = scopes.filter((m) => m.life === 'alive').length}
	{@const livability = (numberOfAliveFacilities * 100) / numberOfScopes || 0}
	<Card.Root>
		<Card.Header>
			<Card.Title>SCOPE</Card.Title>
		</Card.Header>
		<Card.Content class="text-7xl">
			{numberOfScopes}
		</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header>
			<Card.Title>AVAILABLE</Card.Title>
		</Card.Header>
		<Card.Content>
			<p class="text-3xl">
				{Math.round(availability)}%
			</p>
			<p class="text-xs text-muted-foreground">
				{numberOfAvailableFacilities} Available over {numberOfScopes} scopes
			</p>
		</Card.Content>
		<Card.Footer>
			<Progress
				value={availability}
				max={100}
				class={`${
					availability > 62
						? 'bg-green-50 *:bg-green-700'
						: availability > 38
							? 'bg-yellow-50 *:bg-yellow-500'
							: 'bg-red-50 *:bg-red-700'
				}`}
			/>
		</Card.Footer>
	</Card.Root>
	<Card.Root>
		<Card.Header>
			<Card.Title>LIVABILITY</Card.Title>
		</Card.Header>
		<Card.Content>
			<p class="text-3xl">
				{Math.round(livability)}%
			</p>
			<p class="text-xs text-muted-foreground">
				{numberOfAliveFacilities} Active over {numberOfScopes} scopes
			</p>
		</Card.Content>
		<Card.Footer>
			<Progress
				value={livability}
				max={100}
				class={`${
					livability > 62
						? 'bg-green-50 *:bg-green-700'
						: livability > 38
							? 'bg-yellow-50 *:bg-yellow-500'
							: 'bg-red-50 *:bg-red-700'
				}`}
			/>
		</Card.Footer>
	</Card.Root>
{/snippet}
