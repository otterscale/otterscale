DEPRECATED

<!-- <script lang="ts">
	// External dependencies
	import Icon from '@iconify/svelte';
	import { capitalizeFirstLetter } from 'better-auth';
	import * as Card from '$lib/components/ui/card/index.js';
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button';

	// Internal UI components
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Switch } from '$lib/components/ui/switch/index.js';

	// Internal utilities and types
	import { type Model } from '$gen/api/stack/v1/stack_pb';
	import { nodeIcon } from '$lib/node';
	import ListExpression from './ui/list-expression.svelte';
	import ChartMetric from './ui/chart-metric.svelte';
	import CardDictionary from './ui/card-dictionary.svelte';
	import {
		StackService,
		type Application,
		type ImportBootResourcesRequest,
		type PowerOnMachineRequest,
		type PowerOffMachineRequest,
		type CommissionMachineRequest
	} from '$gen/api/stack/v1/stack_pb';
	import { getContext, onMount } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { writable } from 'svelte/store';
	import { formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';

	const transport: Transport = getContext('transport');
	const client = createClient(StackService, transport);

	const applications: Application[] = [];
	const applicationsStore = writable<Application[]>([]);
	const applicationsIsLoading = writable(true);

	async function fetchApplications(modelUuid: string) {
		try {
			const response = await client.listApplications({
				modelUuid: modelUuid
			});
			applicationsStore.set(response.applications);
		} catch (error) {
			console.error('Error fetching machines:', error);
		} finally {
			applicationsIsLoading.set(false);
		}
	}

	const nodeType = 'Kubernetes';

	let {
		model
	}: {
		model: Model;
	} = $props();

	let mounted = $state(false);

	onMount(async () => {
		try {
			await fetchApplications(model.uuid);
			$applicationsStore.flat().forEach((application) => {
				applications.push(application);
			});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<div class="grid gap-3 px-8">
	<div class="flex items-center">
		<div class="flex items-center space-x-2">
			<Icon icon={nodeIcon(nodeType)} class="size-8" />
			<div class="flex-col p-2">
				<div class="font-bold">{model.name}</div>
				<div class="flex text-sm text-muted-foreground">
					{model.uuid}
				</div>
			</div>
		</div>
	</div>
	<div class="grid gap-3">
		<div class="flex items-center gap-1">
			<Badge variant="outline" class="text-base">
				<div class="flex items-center gap-1">
					<Icon icon="ph:wrench" class="size-5" />
					{model.status}
				</div>
			</Badge>
			<Badge variant="outline" class="text-base">
				<div class="flex items-center gap-1">
					<Icon icon="pajamas:status-health" class="size-5" />
					{model.life}
				</div>
			</Badge>
		</div>
		<div>
			{#if model.updatedAt}
				<Badge variant="outline" class="border-fit p-1 text-sm">
					<div class="flex items-center gap-1">
						<Icon icon="ph:clock" class="size-5" />
						{formatTimeAgo(timestampDate(model.updatedAt))}
					</div>
				</Badge>
			{/if}
		</div>
	</div>
	<div class="grid gap-3">
		<div class="grid grid-cols-3 gap-3">
			<ChartMetric metric={`${model.machineCount} pcs`} description="Machine" />
			<ChartMetric metric={`${model.coreCount} Cores`} description="CPU" />
			<ChartMetric metric={`${model.unitCount} pcs`} description="Unit" />
		</div>
	</div>
	{#if mounted}
		{#if applications.length > 0}
			<Card.Root>
				<Card.Header>
					<Card.Title>Applications</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid grid-cols-4 gap-3">
						{#each applications as application}
							<fieldset class="rounded-md border p-3">
								<legend
									><Button variant="ghost">
										<Icon icon="ph:app-window" />
										{application.name}
									</Button></legend
								>
								{#if application.units}
									<div class="items-center justify-between">
										{#if application.units.length > 0}
											{#each application.units as unit}
												<Button variant="ghost" class="text-xs">
													<Icon icon="ph:cube" />
													{unit.name}
												</Button>
											{/each}
										{/if}
									</div>
								{/if}
							</fieldset>
						{/each}
					</div>
				</Card.Content>
			</Card.Root>
		{/if}
	{/if}
</div> -->
