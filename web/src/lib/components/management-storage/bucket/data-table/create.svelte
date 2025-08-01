<script lang="ts" module>
	import type { Bucket, CreateBucketRequest } from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable, type Writable } from 'svelte/store';
	import { accessControlListOptions } from './utils.svelte';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		buckets: data = $bindable()
	}: { selectedScope: string; selectedFacility: string; buckets: Writable<Bucket[]> } = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		policy: '{}'
	} as CreateBucketRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let userOptions = $state(writable<SingleSelect.OptionType[]>([]));
	let isMounted = $state(false);
	onMount(() => {
		storageClient
			.listUsers({ scopeUuid: selectedScope, facilityName: selectedFacility })
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

	let invalid = $state(false);
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
					<div class="flex justify-end gap-2">
						<Button
							variant="outline"
							size="sm"
							href="https://awspolicygen.s3.amazonaws.com/policygen.html"
							target="_blank"
							class="flex items-center gap-1"
						>
							<Icon icon="ph:arrow-square-out" />
							Reference
						</Button>
						<Button
							variant="outline"
							size="sm"
							href="https://awspolicygen.s3.amazonaws.com/policygen.html"
							target="_blank"
							class="flex items-center gap-1"
						>
							<Icon icon="ph:arrow-square-out" />
							Generator
						</Button>
					</div>
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
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					disabled={invalid}
					onclick={() => {
						toast.info(`Creating ${request.bucketName}...`);
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
						stateController.close();
					}}
				>
					Create
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
