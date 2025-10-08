<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { MachineService, type CreateMachineRequest, type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Multiple as MultipleSelect, Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { activeScope } from '$lib/stores';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		machine,
	}: {
		machine: Machine;
	} = $props();
	const reloadManager: ReloadManager = getContext('reloadManager');

	const transport: Transport = getContext('transport');
	const client = createClient(MachineService, transport);
	const tagOptions = writable<SingleSelect.OptionType[]>([]);

	let isTagLoading = $state(true);

	const defaults = {
		scopeUuid: $activeScope?.uuid,
		id: machine.id,
		enableSsh: true,
		skipBmcConfig: false,
		skipNetworking: false,
		skipStorage: false,
		tags: [] as string[],
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
			client
				.listTags({})
				.then((response) => {
					tagOptions.set(
						response.tags.flatMap((tag) => ({
							value: tag.name,
							label: tag.name,
							icon: 'ph:tag',
						})),
					);
				})
				.finally(() => {
					isTagLoading = false;
				});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:compass" />
		{m.add()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.add_machine()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Legend>{m.features()}</Form.Legend>
				<Form.Field>
					<SingleInput.Boolean required descriptor={() => m.enable_ssh()} bind:value={request.enableSsh} />
				</Form.Field>
				<Form.Field>
					<SingleInput.Boolean
						required
						descriptor={() => m.skip_bmc_config()}
						bind:value={request.skipBmcConfig}
					/>
				</Form.Field>
				<Form.Field>
					<SingleInput.Boolean
						required
						descriptor={() => m.skip_networking()}
						bind:value={request.skipNetworking}
					/>
				</Form.Field>
				<Form.Field>
					<SingleInput.Boolean
						required
						descriptor={() => m.skip_storage()}
						bind:value={request.skipStorage}
					/>
				</Form.Field>
			</Form.Fieldset>
			<Form.Fieldset>
				<Form.Legend>
					{m.tags()}
				</Form.Legend>
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
					onclick={() => {
						toast.promise(() => client.createMachine(request), {
							loading: 'Executing...',
							success: (response) => {
								reloadManager.force();
								return `Create ${response.fqdn} success`;
							},
							error: (error) => {
								let message = `Fail to create ${machine.fqdn}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
								});
								return message;
							},
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
