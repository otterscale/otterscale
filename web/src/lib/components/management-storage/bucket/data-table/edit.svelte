<script lang="ts" module>
	import type { Bucket, Bucket_Grant, UpdateBucketRequest } from '$gen/api/storage/v1/storage_pb';
	import { Bucket_ACL, StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { type Writable } from 'svelte/store';
	import { accessControlListOptions } from './create.svelte';
</script>

<script lang="ts">
	function getAccessControlList(grants: Bucket_Grant[]): Bucket_ACL {
		if (grants.some((grant) => grant.uri.includes('AuthenticatedUsers'))) {
			return Bucket_ACL.AUTHENTICATED_READ;
		}

		if (
			grants.some((grant) => grant.uri.includes('AllUsers')) &&
			grants.some((grant) => grant.permission.includes('WRITE'))
		) {
			return Bucket_ACL.PUBLIC_READ_WRITE;
		}

		if (grants.some((grant) => grant.uri.includes('AllUsers'))) {
			return Bucket_ACL.PUBLIC_READ_WRITE;
		}

		return Bucket_ACL.PRIVATE;
	}

	let {
		selectedScope,
		selectedFacility,
		bucket,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		bucket: Bucket;
		data: Writable<Bucket[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		bucketName: bucket.name,
		owner: bucket.owner,
		policy: bucket.policy,
		acl: getAccessControlList(bucket.grants)
	} as UpdateBucketRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);
	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Edit Bucket</AlertDialog.Header>
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
					<SingleInput.Structure preview required bind:value={request.policy} language="json" />
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
							.updateBucket(request)
							.then((r) => {
								toast.success(`Update ${r.name}`);
								storageClient
									.listBuckets({ scopeUuid: selectedScope, facilityName: selectedFacility })
									.then((r) => {
										data.set(r.buckets);
									});
							})
							.catch((e) => {
								toast.error(`Fail to update bucket: ${e.toString()}`);
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
