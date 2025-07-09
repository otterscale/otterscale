<script lang="ts" module>
	import type { Bucket, CreateBucketRequest } from '$gen/api/storage/v1/storage_pb';
	import { Bucket_ACL, StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';

	export const accessControlListOptions = writable([
		{
			value: Bucket_ACL.PRIVATE,
			label: 'PRIVATE',
			icon: 'ph:user'
		},
		{
			value: Bucket_ACL.PUBLIC_READ,
			label: 'PUBLIC_READ',
			icon: 'ph:user'
		},
		{
			value: Bucket_ACL.PUBLIC_READ_WRITE,
			label: 'PUBLIC_READ_WRITE',
			icon: 'ph:user'
		},
		{
			value: Bucket_ACL.AUTHENTICATED_READ,
			label: 'AUTHENTICATED_READ',
			icon: 'ph:user'
		}
	]);
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		data = $bindable()
	}: { selectedScope: string; selectedFacility: string; data: Writable<Bucket[]> } = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility
	} as CreateBucketRequest;
	let request = $state(DEFAULT_REQUEST);
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
		<AlertDialog.Header>Create Bucket</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.bucketName} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Owner</Form.Label>
					<SingleInput.General required type="text" bind:value={request.owner} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Policies</Form.Legend>

				<Form.Field>
					<Form.Label>Policy</Form.Label>
					<SingleInput.Structure bind:value={request.policy} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Access Control List</Form.Label>
					<SingleSelect.Root required options={accessControlListOptions} bind:value={request.acl}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $accessControlListOptions as option}
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
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						stateController.close();
						storageClient
							.createBucket(request)
							.then((r) => {
								toast.success(`Create ${r.name}`);
								storageClient
									.listBuckets({ scopeUuid: selectedScope, facilityName: selectedFacility })
									.then((r) => {
										data.set(r.buckets);
									});
							})
							.catch((e) => {
								toast.error(`Fail to create bucket: ${e.toString()}`);
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
