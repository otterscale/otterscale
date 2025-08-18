<script lang="ts" module>
	import type { CreateImageRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { currentCeph } from '$lib/stores';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable } from 'svelte/store';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let isAdvancedOpen = $state(false);
	let isPoolsLoading = $state(true);
	let isImageNameInvalid = $state(false);
	let isPoolNameInvalid = $state(false);
	let isMounted = $state(false);
	let poolOptions = $state(writable<SingleSelect.OptionType[]>([]));
	const storageClient = createClient(StorageService, transport);
	const requestManager = new RequestManager<CreateImageRequest>({
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		layering: true,
		exclusiveLock: true,
		objectMap: true,
		fastDiff: true,
		deepFlatten: true
	} as CreateImageRequest);
	const stateController = new StateController(false);

	async function fetchVolumeOptions() {
		try {
			const response = await storageClient.listPools({
				scopeUuid: $currentCeph?.scopeUuid,
				facilityName: $currentCeph?.name
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

		isMounted = true;
	});
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		Create
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create RADOS Block Device</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Image Name</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={requestManager.request.imageName}
						bind:invalid={isImageNameInvalid}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Pool Name</Form.Label>
					{#if isPoolsLoading}
						<Loading.Selection />
					{:else}
						<SingleSelect.Root
							required
							bind:options={poolOptions}
							bind:value={requestManager.request.poolName}
							bind:invalid={isPoolNameInvalid}
						>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.Input />
									<SingleSelect.List>
										<SingleSelect.Empty>No results found.</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $poolOptions as option}
												<SingleSelect.Item {option}>
													<Icon
														icon={option.icon ? option.icon : 'ph:empty'}
														class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
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
					<Form.Label>Quota Size</Form.Label>
					<SingleInput.Measurement
						bind:value={requestManager.request.quotaBytes}
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
						<Form.Legend>Striping</Form.Legend>

						<Form.Field>
							<Form.Label>Object Size</Form.Label>
							<SingleInput.Measurement
								bind:value={requestManager.request.objectSizeBytes}
								transformer={(value) => String(value)}
								units={[
									{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
									{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
								]}
							/>
						</Form.Field>

						<Form.Field>
							<Form.Label>Stripe Unit</Form.Label>
							<SingleInput.Measurement
								bind:value={requestManager.request.stripeUnitBytes}
								transformer={(value) => String(value)}
								units={[
									{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
									{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
								]}
							/>
						</Form.Field>

						<Form.Field>
							<Form.Label>Stripe Count</Form.Label>
							<SingleInput.General
								bind:value={requestManager.request.stripeCount}
								transformer={(value) => String(value)}
							/>
						</Form.Field>
					</Form.Fieldset>

					<Form.Fieldset>
						<Form.Legend>Features</Form.Legend>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => 'Allows the creation of snapshots and clones of an image.'}
								format="checkbox"
								bind:value={requestManager.request.layering}
							/>
						</Form.Field>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => 'Ensures that only one client can write to the image at a time.'}
								format="checkbox"
								bind:value={requestManager.request.exclusiveLock}
							/>
						</Form.Field>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => 'Tracks object existence to speed up image operations.'}
								format="checkbox"
								bind:value={requestManager.request.objectMap}
							/>
						</Form.Field>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => 'Speeds up the process of comparing two images.'}
								format="checkbox"
								bind:value={requestManager.request.fastDiff}
							/>
						</Form.Field>

						<Form.Field>
							<SingleInput.Boolean
								descriptor={() => 'Removes clone dependency on parent image for faster deletion.'}
								format="checkbox"
								bind:value={requestManager.request.deepFlatten}
							/>
						</Form.Field>
					</Form.Fieldset>
				</Collapsible.Content>
			</Collapsible.Root>
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
					disabled={isImageNameInvalid || isPoolNameInvalid}
					onclick={() => {
						toast.promise(() => storageClient.createImage(requestManager.request), {
							loading: `Creating ${requestManager.request.imageName}...`,
							success: (response) => {
								reloadManager.force();
								return `Create ${requestManager.request.imageName}`;
							},
							error: (error) => {
								let message = `Fail to create ${requestManager.request.imageName}`;
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
