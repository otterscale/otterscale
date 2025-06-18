<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import {
		LayeredMultiple as LayeredMultipleSelect,
		Multiple as MultipleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import {
		capacities,
		type CreateS3KeyRequest,
		maximumBucketsOptions,
		type Request,
		s3keys
	} from './create.svelte';
	import type { User } from './types';

	let { user }: { user: User } = $props();

	let isShowTenant = $state(false);
	let maximumBucketsOption = $state('');

	const DEFAULT_REQUEST = { userId: user.userId } as Request;
	let request: Request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const DEFAULT_CREATE_S3_KEY_REQUEST = {} as CreateS3KeyRequest;
	let createS3KeyRequest: CreateS3KeyRequest = $state(DEFAULT_CREATE_S3_KEY_REQUEST);
	function resetCreateS3KeyRequest() {
		createS3KeyRequest = DEFAULT_CREATE_S3_KEY_REQUEST;
	}

	const stateController = new DialogStateController(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			Edit User
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="filesystem-name">User ID</Form.Label>
					<SingleInput.General required type="text" bind:value={request.userId} />
				</Form.Field>

				<SingleInput.Boolean
					required
					class="border-none shadow-none ring-0"
					descriptor={() => {
						if (isShowTenant) {
							return 'Show Tenant';
						} else {
							return 'Not Show Tenant';
						}
					}}
					bind:value={isShowTenant}
				/>
				{#if isShowTenant}
					<Form.Field>
						<Form.Label for="filesystem-placement">Tenant</Form.Label>
						<SingleInput.General type="text" bind:value={request.tenant} />
					</Form.Field>
				{/if}

				<Form.Field>
					<Form.Label for="filesystem-placement">Fullname</Form.Label>
					<SingleInput.General required type="text" bind:value={request.userId} />
				</Form.Field>

				<Form.Field>
					<Form.Label for="filesystem-placement">Email Address</Form.Label>
					<SingleInput.General type="email" bind:value={request.emailAddress} />
				</Form.Field>

				<Form.Field>
					<Form.Label for="filesystem-placement">Maximum Bucket</Form.Label>
					<SingleSelect.Root options={maximumBucketsOptions} bind:value={maximumBucketsOption}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $maximumBucketsOptions as option}
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
					{#if maximumBucketsOption === 'custom'}
						<Form.Field>
							<SingleInput.General type="number" bind:value={request.maximumBuckets} />
						</Form.Field>
					{/if}
				</Form.Field>

				<Form.Field>
					<Form.Label for="filesystem-placement">Suspended</Form.Label>
					<SingleInput.Boolean required bind:value={request.suspended} />
				</Form.Field>
				<Form.Help>
					System user S3 key Auto-generate key User quota Enabled Bucket quota Enabled Suspending
					the user disables the user and subuser
				</Form.Help>

				<Form.Field>
					<Form.Label for="filesystem-placement">System User</Form.Label>
					<SingleInput.Boolean required bind:value={request.systemUser} />
				</Form.Field>
				<Form.Help>
					System users are distinct from regular users, they are used by the RGW service to perform
					administrative tasks, manage buckets and objects
				</Form.Help>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Capabilities</Form.Legend>
				<LayeredMultipleSelect.Root bind:value={request.capacities} options={capacities}>
					<LayeredMultipleSelect.Viewer />
					<LayeredMultipleSelect.Controller>
						<LayeredMultipleSelect.Trigger />
						<LayeredMultipleSelect.Content>
							<LayeredMultipleSelect.Group>
								{#each capacities as option}
									{#if option.subOptions && option.subOptions.length > 0}
										{#snippet child(
											options: LayeredMultipleSelect.OptionType[],
											option: LayeredMultipleSelect.OptionType,
											parents: LayeredMultipleSelect.OptionType[],
											level: number = 0
										)}
											<LayeredMultipleSelect.Sub>
												<LayeredMultipleSelect.SubTrigger>
													<Icon
														icon={option.icon && option.icon ? option.icon : 'ph:empty'}
														class={cn(
															'size-5',
															option.icon && option.icon ? 'visibale' : 'invisible'
														)}
													/>
													{option.label}
												</LayeredMultipleSelect.SubTrigger>
												<LayeredMultipleSelect.SubContent>
													{#each options as option}
														{#if option.subOptions && option.subOptions.length > 0}
															{@render child(
																option.subOptions,
																option,
																[...parents, option],
																level + 1
															)}
														{:else}
															<LayeredMultipleSelect.Item {option} {parents}>
																<Icon
																	icon={option.icon && option.icon ? option.icon : 'ph:empty'}
																	class={cn(
																		'size-5',
																		option.icon && option.icon ? 'visibale' : 'invisible'
																	)}
																/>
																{option.label}
																<LayeredMultipleSelect.Check {option} {parents} />
															</LayeredMultipleSelect.Item>
														{/if}
													{/each}
												</LayeredMultipleSelect.SubContent>
											</LayeredMultipleSelect.Sub>
										{/snippet}

										{@render child(option.subOptions, option, [option])}
									{:else}
										<LayeredMultipleSelect.Item {option}>
											<Icon
												icon={option.icon && option.icon ? option.icon : 'ph:empty'}
												class={cn('size-5', option.icon && option.icon ? 'visibale' : 'invisible')}
											/>
											{option.label}
											<LayeredMultipleSelect.Check {option} />
										</LayeredMultipleSelect.Item>
									{/if}
								{/each}
							</LayeredMultipleSelect.Group>
							<LayeredMultipleSelect.Actions>
								<LayeredMultipleSelect.ActionAll>All</LayeredMultipleSelect.ActionAll>
								<LayeredMultipleSelect.ActionClear>Clear</LayeredMultipleSelect.ActionClear>
							</LayeredMultipleSelect.Actions>
						</LayeredMultipleSelect.Content>
					</LayeredMultipleSelect.Controller>
				</LayeredMultipleSelect.Root>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>S3 Keys</Form.Legend>
				<MultipleSelect.Root bind:value={request.s3keyIds} options={s3keys}>
					<MultipleSelect.Viewer />
					<MultipleSelect.Controller>
						<MultipleSelect.Trigger />
						<MultipleSelect.Content>
							<MultipleSelect.Options>
								<MultipleSelect.Input>
									{#snippet addition({ accessor, manager })}
										<AlertDialog.Root>
											<AlertDialog.Trigger
												onclick={() => {
													createS3KeyRequest.accessKey = accessor.input;
												}}
												class={cn(buttonVariants({ variant: 'outline' }))}
											>
												<div class="flex items-center gap-1">
													<Icon icon="ph:plus" />
													<p class="text-xs">Add</p>
												</div>
											</AlertDialog.Trigger>
											<AlertDialog.Content>
												<AlertDialog.Header
													class="flex items-center justify-center text-xl font-bold"
												>
													Create S3 Key
												</AlertDialog.Header>
												<SingleInput.Boolean
													format="checkbox"
													required
													descriptor={() => {
														if (createS3KeyRequest.autoGenerateKey === true) {
															return 'Automatically Generate Key';
														} else if (createS3KeyRequest.autoGenerateKey === false) {
															return 'Not Automatically Generate Key';
														} else {
															return '';
														}
													}}
													bind:value={createS3KeyRequest.autoGenerateKey}
												/>
												{#if !createS3KeyRequest.autoGenerateKey}
													<Form.Field>
														<Form.Label>Access Key</Form.Label>
														<SingleInput.Password
															required
															bind:value={createS3KeyRequest.accessKey}
														/>
													</Form.Field>

													<Form.Field>
														<Form.Label>Secret Key</Form.Label>
														<SingleInput.Password
															required
															bind:value={createS3KeyRequest.secretKey}
														/>
													</Form.Field>
												{/if}
												<AlertDialog.Footer>
													<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
													<AlertDialog.Action
														onclick={() => {
															s3keys.set([
																...$s3keys,
																{
																	value: createS3KeyRequest.accessKey,
																	label: createS3KeyRequest.accessKey,
																	icon: 'ph:lock'
																} as MultipleSelect.OptionType
															]);
															manager.updateOptions($s3keys);
															resetCreateS3KeyRequest();
															accessor.input = '';
														}}
													>
														Confirm
													</AlertDialog.Action>
												</AlertDialog.Footer>
											</AlertDialog.Content>
										</AlertDialog.Root>
									{/snippet}
								</MultipleSelect.Input>
								<MultipleSelect.List>
									<MultipleSelect.Empty>No results found.</MultipleSelect.Empty>
									<MultipleSelect.Group>
										{#each $s3keys as option}
											<MultipleSelect.Item {option}>
												<Icon
													icon={option.icon ? option.icon : 'ph:empty'}
													class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
												/>
												{option.label}
												<MultipleSelect.Check {option} />
											</MultipleSelect.Item>
										{/each}
									</MultipleSelect.Group>
								</MultipleSelect.List>
								<MultipleSelect.Actions>
									<MultipleSelect.ActionAll>All</MultipleSelect.ActionAll>
									<MultipleSelect.ActionClear>Clear</MultipleSelect.ActionClear>
								</MultipleSelect.Actions>
							</MultipleSelect.Options>
						</MultipleSelect.Content>
					</MultipleSelect.Controller>
				</MultipleSelect.Root>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>User Quota</Form.Legend>
				<SingleInput.Boolean
					format="checkbox"
					required
					descriptor={() => {
						if (request.enableUserQuota === true) {
							return 'Enable User Quota';
						} else if (request.enableUserQuota === false) {
							return 'Not Enable User Quota';
						} else {
							return '';
						}
					}}
					bind:value={request.enableUserQuota}
				/>
				{#if request.enableUserQuota}
					<Form.Field>
						<Form.Label>Unlimited Size</Form.Label>
						<SingleInput.Boolean required bind:value={request.unlimitedUserSize} />
						{#if !request.unlimitedUserSize}
							<Form.Label>Maximum Size</Form.Label>
							<SingleInput.General required type="number" bind:value={request.maximumUserSize} />
						{/if}
					</Form.Field>

					<Form.Field>
						<Form.Label>Unlimited Objects</Form.Label>
						<SingleInput.Boolean required bind:value={request.unlimitedUserObjects} />
						{#if !request.unlimitedUserObjects}
							<Form.Label>Maximum Objects</Form.Label>
							<SingleInput.General required type="number" bind:value={request.maximumUserObjects} />
						{/if}
					</Form.Field>
				{/if}
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Bucket Quota</Form.Legend>
				<SingleInput.Boolean
					format="checkbox"
					required
					descriptor={() => {
						if (request.enableBucketQuota === true) {
							return 'Enable Bucket Quota';
						} else if (request.enableBucketQuota === false) {
							return 'Not Enable Bucket Quota';
						} else {
							return '';
						}
					}}
					bind:value={request.enableBucketQuota}
				/>
				{#if request.enableBucketQuota}
					<Form.Field>
						<Form.Label>Unlimited Size</Form.Label>
						<SingleInput.Boolean required bind:value={request.unlimitedBucketSize} />
						{#if !request.unlimitedBucketSize}
							<Form.Label>Maximum Size</Form.Label>
							<SingleInput.General required type="number" bind:value={request.maximumBucketSize} />
						{/if}
					</Form.Field>

					<Form.Field>
						<Form.Label>Unlimited Objects</Form.Label>
						<SingleInput.Boolean required bind:value={request.unlimitedBucketObjects} />
						{#if !request.unlimitedBucketObjects}
							<Form.Label>Maximum Objects</Form.Label>
							<SingleInput.General
								required
								type="number"
								bind:value={request.maximumBucketObjects}
							/>
						{/if}
					</Form.Field>
				{/if}
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						console.log(request);
					}}
				>
					Create
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
