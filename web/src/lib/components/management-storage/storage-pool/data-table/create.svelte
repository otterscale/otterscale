<script lang="ts" module>
	import type { CreatePoolRequest, Pool } from '$gen/api/storage/v1/storage_pb';
	import { PoolType, StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import {
		Multiple as MultipleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';
	import {
		QUOTAS_BYTES_HELP_TEXT,
		QUOTAS_OBJECTS_HELP_TEXT,
		REPLCATED_SIZE_HELP_TEXT
	} from './helper';

	export const poolTypes: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: PoolType.ERASURE,
			label: 'Erasure',
			icon: 'ph:scales'
		},
		{
			value: PoolType.REPLICATED,
			label: 'Replicated',
			icon: 'ph:copy-simple'
		}
	]);

	export const applications: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'cephfs',
			label: 'Ceph File System',
			icon: 'ph:squares-four'
		},
		{
			value: 'rbd',
			label: 'RADOS Block Device',
			icon: 'ph:squares-four'
		},
		{
			value: 'rgw',
			label: 'RADOS Gateway',
			icon: 'ph:squares-four'
		}
	]);
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		data = $bindable()
	}: { selectedScope: string; selectedFacility: string; data: Writable<Pool[]> } = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility
	} as CreatePoolRequest;
	let request: CreatePoolRequest = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<div class="flex justify-end">
		<AlertDialog.Trigger class={cn(buttonVariants({ variant: 'default', size: 'sm' }))}>
			<div class="flex items-center gap-1">
				<Icon icon="ph:plus" />
				Create
			</div>
		</AlertDialog.Trigger>
	</div>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Create Pool
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.poolName} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Type</Form.Label>
					<SingleSelect.Root required options={poolTypes} bind:value={request.poolType}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $poolTypes as type}
											<SingleSelect.Item option={type}>
												<Icon
													icon={type.icon ? type.icon : 'ph:empty'}
													class={cn('size-5', type.icon ? 'visibale' : 'invisible')}
												/>
												{type.label}
												<SingleSelect.Check option={type} />
											</SingleSelect.Item>
										{/each}
									</SingleSelect.Group>
								</SingleSelect.List>
							</SingleSelect.Options>
						</SingleSelect.Content>
					</SingleSelect.Root>
				</Form.Field>

				{#if request.poolType === PoolType.ERASURE}
					<Form.Field>
						<Form.Label>EC Overwrite</Form.Label>
						<SingleInput.Boolean
							required
							descriptor={(value) => {
								if (value === true) {
									return 'EC Overwrites';
								} else if (value === false) {
									return 'EC Not Overwrites';
								} else {
									return 'Undetermined';
								}
							}}
							bind:value={request.ecOverwrites}
						/>
					</Form.Field>
				{/if}

				{#if request.poolType === PoolType.REPLICATED}
					<Form.Field>
						<Form.Label>Replcated Size</Form.Label>
						<SingleInput.General required type="number" bind:value={request.replicatedSize} />
					</Form.Field>
					<Form.Help>
						{REPLCATED_SIZE_HELP_TEXT}
					</Form.Help>
				{/if}

				<Form.Field>
					<Form.Label>Applications</Form.Label>
					<MultipleSelect.Root bind:value={request.applications} options={applications}>
						<MultipleSelect.Viewer />
						<MultipleSelect.Controller>
							<MultipleSelect.Trigger />
							<MultipleSelect.Content>
								<MultipleSelect.Options>
									<MultipleSelect.Input />
									<MultipleSelect.List>
										<MultipleSelect.Empty>No results found.</MultipleSelect.Empty>
										<MultipleSelect.Group>
											{#each $applications as application}
												<MultipleSelect.Item option={application}>
													<Icon
														icon={application.icon ? application.icon : 'ph:empty'}
														class={cn('size-5', application.icon ? 'visibale' : 'invisible')}
													/>
													{application.label}
													<MultipleSelect.Check option={application} />
												</MultipleSelect.Item>
											{/each}
										</MultipleSelect.Group>
									</MultipleSelect.List>
									<MultipleSelect.Actions>
										<MultipleSelect.ActionAll>All</MultipleSelect.ActionAll>
										<MultipleSelect.ActionClear>Clear</MultipleSelect.ActionClear>
									</MultipleSelect.Actions>
								</MultipleSelect.Options>
							</MultipleSelect.Content>
						</MultipleSelect.Controller>
					</MultipleSelect.Root>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Quotas</Form.Legend>

				<Form.Field>
					<Form.Label>Bytes</Form.Label>
					<SingleInput.General type="number" bind:value={request.quotaBytes} />
				</Form.Field>
				<Form.Help>
					{QUOTAS_BYTES_HELP_TEXT}
				</Form.Help>

				<Form.Field>
					<Form.Label>Objects</Form.Label>
					<SingleInput.General type="number" bind:value={request.quotaObjects} />
				</Form.Field>
				<Form.Help>
					{QUOTAS_OBJECTS_HELP_TEXT}
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						stateController.close();
						storageClient
							.createPool(request)
							.then((r) => {
								toast.success(`Create ${r.name}`);
								storageClient
									.listPools({ scopeUuid: selectedScope, facilityName: selectedFacility })
									.then((r) => {
										data.set(r.pools);
									});
							})
							.catch((e) => {
								toast.error(`Fail to create pool: ${e.toString()}`);
							});
					}}>Create</AlertDialog.Action
				>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
