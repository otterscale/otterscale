<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { CreatePoolRequest } from '$lib/api/storage/v1/storage_pb';
	import { PoolType, StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Multiple as MultipleSelect, Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { currentCeph } from '$lib/stores';
	import { cn } from '$lib/utils';

	export const poolTypes: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: PoolType.ERASURE,
			label: 'Erasure',
			icon: 'ph:scales',
		},
		{
			value: PoolType.REPLICATED,
			label: 'Replicated',
			icon: 'ph:copy-simple',
		},
	]);

	export const applications: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: 'cephfs',
			label: 'Ceph File System',
			icon: 'ph:squares-four',
		},
		{
			value: 'rbd',
			label: 'RADOS Block Device',
			icon: 'ph:squares-four',
		},
		{
			value: 'rgw',
			label: 'RADOS Gateway',
			icon: 'ph:squares-four',
		},
	]);
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	let invalidName: boolean | undefined = $state();
	let invalidType: boolean | undefined = $state();
	let invalidReplicatedSize: boolean | undefined = $state();

	const defaults = {
		scope: $currentCeph?.scope,
		facility: $currentCeph?.name,
	} as CreatePoolRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_pool()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.poolName}
						bind:invalid={invalidName}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.type()}</Form.Label>
					<SingleSelect.Root
						required
						options={poolTypes}
						bind:value={request.poolType}
						bind:invalid={invalidType}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
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
						<!-- <Form.Label>{m.ec_overwrite()}</Form.Label> -->
						<SingleInput.Boolean descriptor={() => m.ec_overwrite()} bind:value={request.ecOverwrites} />
					</Form.Field>
				{/if}
				{#if request.poolType === PoolType.REPLICATED}
					<Form.Field>
						<Form.Label>{m.replicated_size()}</Form.Label>
						<SingleInput.General
							required
							bind:value={request.replicatedSize}
							bind:invalid={invalidReplicatedSize}
						/>
					</Form.Field>
					<Form.Help>
						{m.pool_replicated_size_direction()}
					</Form.Help>
				{/if}
				<Form.Field>
					<Form.Label>{m.applications()}</Form.Label>
					<MultipleSelect.Root bind:value={request.applications} options={applications}>
						<MultipleSelect.Viewer />
						<MultipleSelect.Controller>
							<MultipleSelect.Trigger />
							<MultipleSelect.Content>
								<MultipleSelect.Options>
									<MultipleSelect.Input />
									<MultipleSelect.List>
										<MultipleSelect.Empty>{m.no_result()}</MultipleSelect.Empty>
										<MultipleSelect.Group>
											{#each $applications as application}
												<MultipleSelect.Item option={application}>
													<Icon
														icon={application.icon ? application.icon : 'ph:empty'}
														class={cn(
															'size-5',
															application.icon ? 'visibale' : 'invisible',
														)}
													/>
													{application.label}
													<MultipleSelect.Check option={application} />
												</MultipleSelect.Item>
											{/each}
										</MultipleSelect.Group>
									</MultipleSelect.List>
									<MultipleSelect.Actions>
										<MultipleSelect.ActionAll>{m.all()}</MultipleSelect.ActionAll>
										<MultipleSelect.ActionClear>{m.clear()}</MultipleSelect.ActionClear>
									</MultipleSelect.Actions>
								</MultipleSelect.Options>
							</MultipleSelect.Content>
						</MultipleSelect.Controller>
					</MultipleSelect.Root>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.quota_size()}</Form.Label>
					<Form.Help>
						{m.pool_quota_objects_direction()}
					</Form.Help>
					<SingleInput.Measurement
						bind:value={request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType,
						]}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.quota_objects()}</Form.Label>
					<Form.Help>
						{m.pool_quota_objects_direction()}
					</Form.Help>
					<SingleInput.General bind:value={request.quotaObjects} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalidName ||
						invalidType ||
						(request.poolType === PoolType.REPLICATED && invalidReplicatedSize)}
					onclick={() => {
						toast.promise(() => storageClient.createPool(request), {
							loading: `Creating ${request.poolName}...`,
							success: () => {
								reloadManager.force();
								return `Create ${request.poolName}`;
							},
							error: (error) => {
								let message = `Fail to create ${request.poolName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
								});
								return message;
							},
						});
						reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
