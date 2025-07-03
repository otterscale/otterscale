<script lang="ts" module>
	import type { CreateImageRequest, Image } from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		data = $bindable()
	}: { selectedScope: string; selectedFacility: string; data: Writable<Image[]> } = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		layering: false,
		exclusiveLock: false,
		objectMap: false,
		fastDiff: false,
		deepFlatten: false
	} as CreateImageRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const poolOptions: Writable<SingleSelect.OptionType[]> = writable([]);
	let isPoolsLoading = $state(true);
	async function fetchPools() {
		try {
			const response = await storageClient.listPools({
				scopeUuid: 'b62d195e-3905-4960-85ee-7673f71eb21e',
				facilityName: 'ceph-mon'
			});

			poolOptions.set(
				response.pools.map((pool) => ({
					value: pool.name,
					label: pool.name,
					icon: 'ph:cube'
				}))
			);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isPoolsLoading = false;
		}
	}

	onMount(async () => {
		try {
			await fetchPools();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
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
			Create RADOS Block Device
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Image Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.imageName} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Pool Name</Form.Label>

					<SingleSelect.Root required options={poolOptions} bind:value={request.poolName}>
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
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Striping</Form.Legend>

				<Form.Field>
					<Form.Label>Object Size</Form.Label>
					<SingleInput.General type="number" bind:value={request.objectSizeBytes} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Stripe Unit</Form.Label>
					<SingleInput.General type="number" bind:value={request.stripeUnitBytes} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Stripe Count</Form.Label>
					<SingleInput.General type="number" bind:value={request.stripeCount} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Quotas</Form.Legend>

				<Form.Field>
					<SingleInput.General type="number" bind:value={request.quotaBytes} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Features</Form.Legend>

				<Form.Field>
					<Form.Label>Layering</Form.Label>
					<SingleInput.Boolean required format="checkbox" bind:value={request.layering} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Exclusive Lock</Form.Label>
					<SingleInput.Boolean required format="checkbox" bind:value={request.exclusiveLock} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Object Map</Form.Label>
					<SingleInput.Boolean required format="checkbox" bind:value={request.objectMap} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Fast Diff</Form.Label>
					<SingleInput.Boolean required format="checkbox" bind:value={request.fastDiff} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Deep Flatten</Form.Label>
					<SingleInput.Boolean required format="checkbox" bind:value={request.deepFlatten} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						stateController.close();
						storageClient
							.createImage(request)
							.then((r) => {
								toast.success(`Create ${r.name}`);
								storageClient
									.listImages({ scopeUuid: selectedScope, facilityName: selectedFacility })
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
					}}
				>
					Create
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
