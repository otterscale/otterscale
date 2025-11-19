<script lang="ts" module>
	import { pool_name, image_name } from './../../../../paraglide/messages/zh-hant.js';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { CreateImageRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
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

	let isAdvancedOpen = $state(false);
	let isPoolsLoading = $state(true);

	const poolOptions = writable<SingleSelect.OptionType[]>([]);
	const storageClient = createClient(StorageService, transport);
	const defaults = {
		scope: scope,
		layering: true,
		exclusiveLock: true,
		objectMap: true,
		fastDiff: true,
		deepFlatten: true
	} as CreateImageRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let invalidities = $state({} as Booleanified<CreateImageRequest>);
	const invalid = $derived(invalidities.poolName || invalidities.imageName);

	async function fetchVolumeOptions() {
		try {
			const response = await storageClient.listPools({
				scope: scope
			});
			poolOptions.set(
				response.pools.map(
					(pool) =>
						({
							value: pool.name,
							label: pool.name,
							icon: 'ph:cube'
						}) as SingleSelect.OptionType
				)
			);

			isPoolsLoading = false;
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}
	onMount(async () => {
		try {
			await fetchVolumeOptions();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_rbd()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.image_name()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.imageName}
						bind:invalid={invalidities.imageName}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.pool_name()}</Form.Label>
					{#if isPoolsLoading}
						<Loading.Selection />
					{:else}
						<SingleSelect.Root
							required
							options={poolOptions}
							bind:value={request.poolName}
							bind:invalid={invalidities.poolName}
						>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.Input />
									<SingleSelect.List>
										<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $poolOptions as option}
												<SingleSelect.Item {option}>
													<Icon
														icon={option.icon ? option.icon : 'ph:empty'}
														class={cn('size-5', option.icon ? 'visible' : 'invisible')}
													/>
													{option.label}
													<SingleSelect.Check {option} />
												</SingleSelect.Item>
											{/each}
										</SingleSelect.Group>
									</SingleSelect.List>
								</SingleSelect.Options>
							</SingleSelect.Content>
						</SingleSelect.Root>
					{/if}
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.quota_size()}</Form.Label>
					<SingleInput.Measurement
						bind:value={request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>
			</Form.Fieldset>

			<Collapsible.Root bind:open={isAdvancedOpen}>
				<div class="flex items-center justify-between gap-2">
					<p class={cn('text-base font-bold', isAdvancedOpen ? 'invisible' : 'visible')}>Advance</p>
					<Collapsible.Trigger class="bg-muted rounded-full p-1 ">
						<Icon
							icon="ph:caret-left"
							class={cn('transition-all duration-300', isAdvancedOpen ? '-rotate-90' : 'rotate-0')}
						/>
					</Collapsible.Trigger>
				</div>

				<Collapsible.Content>
					<Form.Fieldset>
						<Form.Legend>{m.striping()}</Form.Legend>

						<Form.Field>
							<Form.Label>{m.object_size()}</Form.Label>
							<SingleInput.Measurement
								bind:value={request.objectSizeBytes}
								transformer={(value) => String(value)}
								units={[
									{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
									{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
								]}
							/>
						</Form.Field>

						<Form.Field>
							<Form.Label>{m.stripe_unit()}</Form.Label>
							<SingleInput.Measurement
								bind:value={request.stripeUnitBytes}
								transformer={(value) => String(value)}
								units={[
									{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
									{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
								]}
							/>
						</Form.Field>

						<Form.Field>
							<Form.Label>{m.stripe_count()}</Form.Label>
							<SingleInput.General
								bind:value={request.stripeCount}
								transformer={(value) => String(value)}
							/>
						</Form.Field>
					</Form.Fieldset>

					<Form.Fieldset>
						<Form.Legend>{m.features()}</Form.Legend>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => m.layering_description()}
								bind:value={request.layering}
							/>
						</Form.Field>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => m.exclusive_lock_description()}
								bind:value={request.exclusiveLock}
							/>
						</Form.Field>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => m.object_map_description()}
								bind:value={request.objectMap}
							/>
						</Form.Field>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => m.fast_diff_description()}
								bind:value={request.fastDiff}
							/>
						</Form.Field>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => m.deep_flatten_description()}
								bind:value={request.deepFlatten}
							/>
						</Form.Field>
					</Form.Fieldset>
				</Collapsible.Content>
			</Collapsible.Root>
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
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.createImage(request), {
							loading: `Creating ${request.imageName}...`,
							success: () => {
								reloadManager.force();
								return `Create ${request.imageName}`;
							},
							error: (error) => {
								let message = `Fail to create ${request.imageName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
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
