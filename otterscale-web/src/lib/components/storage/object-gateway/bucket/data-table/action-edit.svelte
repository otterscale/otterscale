<script lang="ts">
	import type { Bucket, UpdateBucketRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { currentCeph } from '$lib/stores';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable } from 'svelte/store';
	import { accessControlListOptions, getAccessControlList } from './utils.svelte';

	let {
		bucket
	}: {
		bucket: Bucket;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	let userOptions = $state(writable<SingleSelect.OptionType[]>([]));
	let isMounted = $state(false);
	let invalid = $state(false);

	const defaults = {
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		bucketName: bucket.name,
		owner: bucket.owner,
		policy: bucket.policy,
		acl: getAccessControlList(bucket.grants)
	} as UpdateBucketRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(() => {
		storageClient
			.listUsers({ scopeUuid: $currentCeph?.scopeUuid, facilityName: $currentCeph?.name })
			.then((response) => {
				userOptions.set(
					response.users.map(
						(user) =>
							({
								value: user.id,
								label: user.id,
								icon: 'ph:user'
							}) as SingleSelect.OptionType
					)
				);
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		Edit
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit Bucket</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General id="name" required type="text" bind:value={request.bucketName} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Owner</Form.Label>
					{#if isMounted}
						<SingleSelect.Root
							id="owner"
							bind:options={userOptions}
							bind:value={request.owner}
							required
						>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.Input />
									<SingleSelect.List>
										<SingleSelect.Empty>No results found.</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $userOptions as option}
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
					{:else}
						<Loading.Selection />
					{/if}
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Policies</Form.Legend>

				<Form.Field>
					<Form.Label>Policy</Form.Label>
					<SingleInput.Structure preview bind:value={request.policy} language="json" />
				</Form.Field>

				<Form.Field>
					<Form.Label>Access Control List</Form.Label>
					<SingleSelect.Root options={accessControlListOptions} bind:value={request.acl}>
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
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				Cancel
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.updateBucket(request), {
							loading: `Updating ${request.bucketName}...`,
							success: (response) => {
								reloadManager.force();
								return `Update ${request.bucketName}`;
							},
							error: (error) => {
								let message = `Fail to update ${request.bucketName}`;
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
					Update
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
