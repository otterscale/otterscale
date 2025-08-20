<script lang="ts">
	import {
		MachineService,
		type CreateMachineRequest,
		type Machine
	} from '$lib/api/machine/v1/machine_pb';
	import { TagService } from '$lib/api/tag/v1/tag_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as SingleStepModal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import {
		Multiple as MultipleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { activeScope } from '$lib/stores';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable } from 'svelte/store';

	let {
		machine
	}: {
		machine: Machine;
	} = $props();
	const reloadManager: ReloadManager = getContext('reloadManager');

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const tagClient = createClient(TagService, transport);

	const tagOptions = writable<SingleSelect.OptionType[]>([]);

	let isTagLoading = $state(true);
	let isMounted = $state(false);

	const defaults = {
		scopeUuid: $activeScope?.uuid,
		id: machine.id,
		enableSsh: true,
		skipBmcConfig: false,
		skipNetworking: false,
		skipStorage: false,
		tags: [] as string[]
	} as CreateMachineRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(async () => {
		try {
			tagClient
				.listTags({})
				.then((response) => {
					tagOptions.set(
						response.tags.flatMap((tag) => ({
							value: tag.name,
							label: tag.name,
							icon: 'ph:tag'
						}))
					);
				})
				.finally(() => {
					isTagLoading = false;
				});

			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<SingleStepModal.Root bind:open>
	<SingleStepModal.Trigger variant="creative">
		<Icon icon="ph:compass" />
		Add
	</SingleStepModal.Trigger>
	<SingleStepModal.Content>
		<SingleStepModal.Header>Add Machine</SingleStepModal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Legend>Features</Form.Legend>
				<Form.Field>
					<SingleInput.Boolean
						required
						descriptor={() => 'Enable SSH'}
						bind:value={request.enableSsh}
					/>
				</Form.Field>
				<Form.Field>
					<SingleInput.Boolean
						required
						descriptor={() => 'Skip BMC Configuration'}
						bind:value={request.skipBmcConfig}
					/>
				</Form.Field>
				<Form.Field>
					<SingleInput.Boolean
						required
						descriptor={() => 'Skip Networking'}
						bind:value={request.skipNetworking}
					/>
				</Form.Field>
				<Form.Field>
					<SingleInput.Boolean
						required
						descriptor={() => 'Skip Storage'}
						bind:value={request.skipStorage}
					/>
				</Form.Field>
			</Form.Fieldset>
			<Form.Fieldset>
				<Form.Legend>Tags</Form.Legend>
				<Form.Field>
					{#if isTagLoading}
						<Loading.Selection />
					{:else}
						<MultipleSelect.Root bind:value={request.tags} options={tagOptions}>
							<MultipleSelect.Viewer />
							<MultipleSelect.Controller>
								<MultipleSelect.Trigger />
								<MultipleSelect.Content>
									<MultipleSelect.Options>
										<MultipleSelect.Input />
										<MultipleSelect.List>
											<MultipleSelect.Empty>No results found.</MultipleSelect.Empty>
											<MultipleSelect.Group>
												{#each $tagOptions as option}
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
					{/if}
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<SingleStepModal.Footer>
			<SingleStepModal.Cancel
				onclick={() => {
					reset();
				}}
			>
				Cancel
			</SingleStepModal.Cancel>
			<SingleStepModal.ActionsGroup>
				<SingleStepModal.Action
					onclick={() => {
						toast.promise(() => machineClient.createMachine(request), {
							loading: 'Executing...',
							success: (response) => {
								reloadManager.force();
								return `Create ${response.fqdn} success`;
							},
							error: (error) => {
								let message = `Fail to create ${machine.fqdn}`;
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
					Create
				</SingleStepModal.Action>
			</SingleStepModal.ActionsGroup>
		</SingleStepModal.Footer>
	</SingleStepModal.Content>
</SingleStepModal.Root>
