<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import {
		LayeredSingle as LayeredSingleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { acls, modes, placementTargets, type Request } from './create.svelte';
	import type { Bucket } from './types';

	let { bucket }: { bucket: Bucket } = $props();

	const DEFAULT_REQUEST = { name: bucket.name } as Request;
	let request: Request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
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
			Edit Bucket
		</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label for="filesystem-name">Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.name} />
				</Form.Field>

				<Form.Field>
					<Form.Label for="filesystem-name">Owner</Form.Label>
					<SingleInput.General required type="text" bind:value={request.owner} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Object Locking</Form.Legend>
				<Form.Help>
					Store objects using a write-once-read-many (WORM) model to prevent objects from being
					deleted or overwritten for a fixed amount of time or indefinitely. Object Locking works
					only in versioned buckets.
				</Form.Help>
				<Form.Field>
					<SingleInput.Boolean
						required
						format="checkbox"
						descriptor={() =>
							'Enables locking for the objects in the bucket. Locking can only be enabled while creating a bucket.'}
						bind:value={request.objectLocking}
					/>
					{#if request.objectLocking}
						<Form.Field>
							<Form.Label>Mode</Form.Label>
							<SingleSelect.Root options={modes} bind:value={request.objectLockingMode}>
								<SingleSelect.Trigger />
								<SingleSelect.Content>
									<SingleSelect.Options>
										<SingleSelect.Input />
										<SingleSelect.List>
											<SingleSelect.Empty>No results found.</SingleSelect.Empty>
											<SingleSelect.Group>
												{#each $modes as option}
													<SingleSelect.Item {option}>
														<SingleSelect.ItemInformation>
															{option.information}
														</SingleSelect.ItemInformation>
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

						<Form.Field>
							<Form.Label>Days</Form.Label>
							<SingleInput.General type="number" bind:value={request.objectLockingDays} />
							<Form.Help>
								The number of days that you want to specify for the default retention period that
								will be applied to new objects placed in this bucket.
							</Form.Help>
						</Form.Field>
					{/if}
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Encryption</Form.Label>
					<SingleInput.Boolean
						required
						format="checkbox"
						descriptor={() =>
							'Enables encryption for the objects in the bucket. To enable encryption on a bucket you need to set the configuration values for SSE-S3 or SSE-KMS.'}
						bind:value={request.encryption}
					/>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Replication</Form.Legend>
				<Form.Field>
					<SingleInput.Boolean
						required
						format="checkbox"
						descriptor={() => 'Enables replication for the objects in the bucket.'}
						bind:value={request.replication}
					/>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Help>Tagging provides a way to categorize storage</Form.Help>
				<Form.Legend>Tags</Form.Legend>
				<Form.Field>
					<MultipleInput.Root
						type="text"
						bind:values={request.tags}
						contextData={{ icon: 'ph:tag' }}
					>
						<MultipleInput.Viewer />
						<MultipleInput.Controller>
							<MultipleInput.Input />
							<MultipleInput.Add />
							<MultipleInput.Clear />
						</MultipleInput.Controller>
					</MultipleInput.Root>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Policies</Form.Legend>
				<Form.Field>
					<Form.Label>Access Control List</Form.Label>
					<LayeredSingleSelect.Root bind:value={request.acl} options={acls} required>
						<LayeredSingleSelect.Trigger />
						<LayeredSingleSelect.Content>
							<LayeredSingleSelect.Group>
								{#each acls as option}
									{#if option.subOptions && option.subOptions.length > 0}
										{#snippet Branch(
											options: LayeredSingleSelect.OptionType[],
											option: LayeredSingleSelect.OptionType,
											parents: LayeredSingleSelect.OptionType[],
											level: number = 0
										)}
											<LayeredSingleSelect.Sub>
												<LayeredSingleSelect.SubTrigger>
													<Icon
														icon={option.icon && option.icon ? option.icon : 'ph:empty'}
														class={cn(
															'size-5',
															option.icon && option.icon ? 'visibale' : 'invisible'
														)}
													/>
													{option.label}
												</LayeredSingleSelect.SubTrigger>
												<LayeredSingleSelect.SubContent>
													{#each options as option}
														{#if option.subOptions && option.subOptions.length > 0}
															{@render Branch(
																option.subOptions,
																option,
																[...parents, option],
																level + 1
															)}
														{:else}
															<LayeredSingleSelect.Item {option} {parents}>
																<Icon
																	icon={option.icon && option.icon ? option.icon : 'ph:empty'}
																	class={cn(
																		'size-5',
																		option.icon && option.icon ? 'visibale' : 'invisible'
																	)}
																/>
																{option.label}
																<LayeredSingleSelect.Check {option} {parents} />
															</LayeredSingleSelect.Item>
														{/if}
													{/each}
												</LayeredSingleSelect.SubContent>
											</LayeredSingleSelect.Sub>
										{/snippet}
										{@render Branch(option.subOptions, option, [option])}
									{:else}
										<LayeredSingleSelect.Item {option}>
											<Icon
												icon={option.icon && option.icon ? option.icon : 'ph:empty'}
												class={cn('size-5', option.icon && option.icon ? 'visibale' : 'invisible')}
											/>
											{option.label}
											<LayeredSingleSelect.Check {option} />
										</LayeredSingleSelect.Item>
									{/if}
								{/each}
							</LayeredSingleSelect.Group>
						</LayeredSingleSelect.Content>
					</LayeredSingleSelect.Root>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Placement Target</Form.Label>
					<SingleSelect.Root options={placementTargets} bind:value={request.placementTarget}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>No results found.</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $placementTargets as option}
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
					<Form.Help>
						When creating a bucket, a placement target can be provided as part of the
						LocationConstraint to override the default placement targets from the user and
						zonegroup.
					</Form.Help>
				</Form.Field>
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
					Confirm
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
