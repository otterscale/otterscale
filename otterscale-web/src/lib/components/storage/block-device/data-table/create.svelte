<script lang="ts" module>
	import type { CreateImageRequest, Image } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		images: data = $bindable()
	}: { selectedScopeUuid: string; selectedFacility: string; images: Writable<Image[]> } = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScopeUuid,
		facilityName: selectedFacility,
		layering: true,
		exclusiveLock: true,
		objectMap: true,
		fastDiff: true,
		deepFlatten: true
	} as CreateImageRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new StateController(false);

	let advance = $state(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let poolOptions = $state(writable<SingleSelect.OptionType[]>([]));
	let isPoolsLoading = $state(true);
	async function fetchVolumeOptions() {
		try {
			const response = await storageClient.listPools({
				scopeUuid: selectedScopeUuid,
				facilityName: selectedFacility
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
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isPoolsLoading = false;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchVolumeOptions();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		isMounted = true;
	});

	let isImageNameInvalid = $state(false);
	let isPoolNameInvalid = $state(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="default">
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
						bind:value={request.imageName}
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
							bind:value={request.poolName}
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
						bind:value={request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>
			</Form.Fieldset>

			<Collapsible.Root bind:open={advance}>
				<div class="flex items-center justify-between gap-2">
					<p class={cn('text-base font-bold', advance ? 'invisible' : 'visible')}>Advance</p>
					<Collapsible.Trigger class="bg-muted rounded-full p-1 ">
						<Icon
							icon="ph:caret-left"
							class={cn('transition-all duration-300', advance ? '-rotate-90' : 'rotate-0')}
						/>
					</Collapsible.Trigger>
				</div>

				<Collapsible.Content>
					<Form.Fieldset>
						<Form.Legend>Striping</Form.Legend>

						<Form.Field>
							<Form.Label>Object Size</Form.Label>
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
							<Form.Label>Stripe Unit</Form.Label>
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
							<Form.Label>Stripe Count</Form.Label>
							<SingleInput.General
								bind:value={request.stripeCount}
								transformer={(value) => String(value)}
							/>
						</Form.Field>
					</Form.Fieldset>

					<Form.Fieldset>
						<Form.Legend>Features</Form.Legend>

						<Form.Field>
							<Form.Label>Layering</Form.Label>
							<SingleInput.Boolean
								descriptor={() => 'Allows the creation of snapshots and clones of an image.'}
								format="checkbox"
								bind:value={request.layering}
							/>
						</Form.Field>

						<Form.Field>
							<Form.Label>Exclusive Lock</Form.Label>
							<SingleInput.Boolean
								descriptor={() => 'Ensures that only one client can write to the image at a time.'}
								format="checkbox"
								bind:value={request.exclusiveLock}
							/>
						</Form.Field>

						<Form.Field>
							<Form.Label>Object Map</Form.Label>
							<SingleInput.Boolean
								descriptor={() =>
									'Tracks object existence to speed up I/O operations for cloning, importing/exporting sparse images, and deletion.'}
								format="checkbox"
								bind:value={request.objectMap}
							/>
						</Form.Field>

						<Form.Field>
							<Form.Label>Fast Diff</Form.Label>
							<SingleInput.Boolean
								descriptor={() => 'Speeds up the process of comparing two images.'}
								format="checkbox"
								bind:value={request.fastDiff}
							/>
						</Form.Field>

						<Form.Field>
							<Form.Label>Deep Flatten</Form.Label>
							<SingleInput.Boolean
								descriptor={() =>
									'Speeds up the process of deleting a clone by removing the dependency on the parent image.'}
								format="checkbox"
								bind:value={request.deepFlatten}
							/>
						</Form.Field>
					</Form.Fieldset>
				</Collapsible.Content>
			</Collapsible.Root>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={isImageNameInvalid || isPoolNameInvalid}
					onclick={() => {
						toast.info(`Creating ${request.imageName}...`);
						storageClient
							.createImage(request)
							.then((r) => {
								toast.success(`Create ${r.name}`);
								storageClient
									.listImages({ scopeUuid: selectedScopeUuid, facilityName: selectedFacility })
									.then((r) => {
										data.set(r.images);
									});
							})
							.catch((e) => {
								toast.error(`Fail to create image: ${e.toString()}`);
							})
							.finally(() => {
								reset();
							});
						stateController.close();
					}}
				>
					Create
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
