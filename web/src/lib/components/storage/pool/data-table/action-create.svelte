<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { CreatePoolRequest } from '$lib/api/storage/v1/storage_pb';
	import { Pool_Application, Pool_Type, StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	export const types: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: Pool_Type.ERASURE,
			label: 'Erasure',
			icon: 'ph:scales'
		},
		{
			value: Pool_Type.REPLICATED,
			label: 'Replicated',
			icon: 'ph:copy-simple'
		}
	]);

	// Fix as Block temperary
	export const applications: Writable<SingleSelect.OptionType[]> = writable([
		{
			value: Pool_Application.FILE,
			label: 'Ceph File System',
			icon: 'ph:squares-four'
		},
		{
			value: Pool_Application.BLOCK,
			label: 'RADOS Block Device',
			icon: 'ph:squares-four'
		},
		{
			value: Pool_Application.OBJECT,
			label: 'RADOS Gateway',
			icon: 'ph:squares-four'
		}
	]);
</script>

<script lang="ts">
	let {
		scope,
		reloadManager
	}: {
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let isAdvancedOpen = $state(false);
	function initAdvanced() {
		isAdvancedOpen = false;
	}

	let request = $state({} as CreatePoolRequest);
	function init() {
		request = {
			scope: scope,
			applications: [Pool_Application.BLOCK]
		} as CreatePoolRequest;

		initAdvanced();
	}

	let invalidity = $state({} as Booleanified<CreatePoolRequest>);
	const invalid = $derived(
		invalidity.poolName ||
			invalidity.type ||
			(request.type === Pool_Type.REPLICATED && invalidity.replicatedSize)
	);

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
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
					<SingleInput.GeneralRule
						required
						type="text"
						bind:value={request.poolName}
						bind:invalid={invalidity.poolName}
						validateRule="lower-alphanum-dash-dot"
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.type()}</Form.Label>
					<SingleSelect.Root
						required
						options={types}
						bind:value={request.type}
						bind:invalid={invalidity.type}
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $types as type (type.value)}
											<SingleSelect.Item option={type}>
												<Icon
													icon={type.icon ? type.icon : 'ph:empty'}
													class={cn('size-5', type.icon ? 'visible' : 'invisible')}
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
				{#if request.type === Pool_Type.ERASURE}
					<Form.Field>
						<!-- <Form.Label>{m.ec_overwrite()}</Form.Label> -->
						<Form.Help>{m.pool_erasure_limit_direction()}</Form.Help>
						<SingleInput.Boolean
							descriptor={() => m.ec_overwrite()}
							bind:value={request.ecOverwrites}
						/>
					</Form.Field>
				{/if}
				{#if request.type === Pool_Type.REPLICATED}
					<Form.Field>
						<Form.Label>{m.replicated_size()}</Form.Label>
						<SingleInput.General
							required
							bind:value={request.replicatedSize}
							bind:invalid={invalidity.replicatedSize}
						/>
					</Form.Field>
					<Form.Help>
						{m.pool_replicated_size_direction()}
					</Form.Help>
				{/if}
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
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>
				<Form.Field>
					<Collapsible.Root bind:open={isAdvancedOpen}>
						<div class="flex items-center justify-between gap-2">
							<p class={cn('text-base font-bold', isAdvancedOpen ? 'invisible' : 'visible')}>
								Advance
							</p>
							<Collapsible.Trigger class="rounded-full bg-muted p-1 ">
								<Icon
									icon="ph:caret-left"
									class={cn(
										'transition-all duration-300',
										isAdvancedOpen ? '-rotate-90' : 'rotate-0'
									)}
								/>
							</Collapsible.Trigger>
						</div>

						<Collapsible.Content>
							<Form.Fieldset>
								<Form.Legend>{m.quota_objects()}</Form.Legend>

								<Form.Field>
									<Form.Label>{m.quota_objects()}</Form.Label>
									<Form.Help>
										{m.pool_quota_objects_direction()}
									</Form.Help>
									<SingleInput.General bind:value={request.quotaObjects} />
								</Form.Field>
							</Form.Fieldset>
						</Collapsible.Content>
					</Collapsible.Root>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
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
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
