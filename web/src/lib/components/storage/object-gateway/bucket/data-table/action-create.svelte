<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { CreateBucketRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';

	import { accessControlListOptions } from './utils.svelte';
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
	let userOptions = $state(writable<SingleSelect.OptionType[]>([]));
	let isMounted = $state(false);

	const defaults = {
		scope: scope,
		policy: '{}'
	} as CreateBucketRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let invalidity = $state({} as Booleanified<CreateBucketRequest>);
	const invalid = $derived(invalidity.bucketName || invalidity.owner);

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(async () => {
		try {
			const response = await storageClient.listUsers({ scope: scope });
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
		<Modal.Header>{m.create_bucket()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General
						id="name"
						required
						type="text"
						bind:value={request.bucketName}
						bind:invalid={invalidity.bucketName}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.owner()}</Form.Label>
					{#if isMounted}
						<SingleSelect.Root
							id="owner"
							bind:options={userOptions}
							bind:value={request.owner}
							required
							bind:invalid={invalidity.owner}
						>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.Input />
									<SingleSelect.List>
										<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $userOptions as option (option.value)}
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
					{:else}
						<Loading.Selection />
					{/if}
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>{m.policies()}</Form.Legend>

				<Form.Field>
					<Form.Label>{m.policy()}</Form.Label>
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
							{m.reference()}
						</Button>
						<Button
							variant="outline"
							size="sm"
							href="https://awspolicygen.s3.amazonaws.com/policygen.html"
							target="_blank"
							class="flex items-center gap-1"
						>
							<Icon icon="ph:arrow-square-out" />
							{m.generator()}
						</Button>
					</div>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.access_control_list()}</Form.Label>
					<SingleSelect.Root options={accessControlListOptions} bind:value={request.acl}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $accessControlListOptions as option (option.value)}
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
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.createBucket(request), {
							loading: `Creating ${request.bucketName}...`,
							success: () => {
								reloadManager.force();
								return `Create ${request.bucketName}`;
							},
							error: (error) => {
								let message = `Fail to create ${request.bucketName}`;
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
