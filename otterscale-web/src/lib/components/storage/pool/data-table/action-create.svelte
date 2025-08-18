<script lang="ts" module>
	import type { CreatePoolRequest } from '$lib/api/storage/v1/storage_pb';
	import { PoolType, StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import {
		Multiple as MultipleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { currentCeph } from '$lib/stores';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
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
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	let invalidName: boolean | undefined = $state();
	let invalidType: boolean | undefined = $state();
	let invalidReplicatedSize: boolean | undefined = $state();

	const requestManager = new RequestManager<CreatePoolRequest>({
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name
	} as CreatePoolRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		Create
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create Pool</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General
						id="name"
						required
						type="text"
						bind:value={requestManager.request.poolName}
						bind:invalid={invalidName}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Type</Form.Label>
					<SingleSelect.Root
						id="type"
						required
						options={poolTypes}
						bind:value={requestManager.request.poolType}
						bind:invalid={invalidType}
					>
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

				{#if requestManager.request.poolType === PoolType.ERASURE}
					<Form.Field>
						<Form.Label>EC Overwrite</Form.Label>
						<SingleInput.Boolean
							format="checkbox"
							descriptor={(value) => {
								if (value === true) {
									return 'EC Overwrites';
								} else if (value === false) {
									return 'EC Not Overwrites';
								} else {
									return 'Undetermined';
								}
							}}
							bind:value={requestManager.request.ecOverwrites}
						/>
					</Form.Field>
				{/if}

				{#if requestManager.request.poolType === PoolType.REPLICATED}
					<Form.Field>
						<Form.Label>Replcated Size</Form.Label>
						<SingleInput.General
							required
							bind:value={requestManager.request.replicatedSize}
							bind:invalid={invalidReplicatedSize}
						/>
					</Form.Field>
					<Form.Help>
						{REPLCATED_SIZE_HELP_TEXT}
					</Form.Help>
				{/if}

				<Form.Field>
					<Form.Label>Applications</Form.Label>
					<MultipleSelect.Root
						bind:value={requestManager.request.applications}
						options={applications}
					>
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

				<Form.Field>
					<Form.Label>Quota Size</Form.Label>
					<Form.Help>
						{QUOTAS_BYTES_HELP_TEXT}
					</Form.Help>
					<SingleInput.Measurement
						bind:value={requestManager.request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Quota Objects</Form.Label>
					<Form.Help>
						{QUOTAS_OBJECTS_HELP_TEXT}
					</Form.Help>
					<SingleInput.General bind:value={requestManager.request.quotaObjects} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					requestManager.reset();
				}}
			>
				Cancel
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalidName ||
						invalidType ||
						(requestManager.request.poolType === PoolType.REPLICATED && invalidReplicatedSize)}
					onclick={() => {
						toast.promise(() => storageClient.createPool(requestManager.request), {
							loading: `Creating ${requestManager.request.poolName}...`,
							success: (response) => {
								reloadManager.force();
								return `Create ${requestManager.request.poolName}`;
							},
							error: (error) => {
								let message = `Fail to create ${requestManager.request.poolName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						requestManager.reset();
						stateController.close();
					}}
				>
					Create
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
