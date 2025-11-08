<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import {
		ConfigurationService,
		type Configuration_BootImage,
		type UpdateBootImageRequest
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Form from '$lib/components/custom/form';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import {
		Multiple as MultipleSelect,
		Single as SingleSelect
	} from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	// Component props - accepts a BootImage object
	let { bootImage }: { bootImage: Configuration_BootImage } = $props();

	// Get required services from Svelte context
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	// Create gRPC client for configuration operations
	const client = createClient(ConfigurationService, transport);

	// Architecture options for the current distro series
	let architecturesOptions = $state(writable<SingleSelect.OptionType[]>([]));

	// Default values for the update boot image request
	const defaults = {
		id: bootImage.id,
		distroSeries: bootImage.distroSeries,
		architectures: [...bootImage.architectures]
	} as UpdateBootImageRequest;

	// Current request state
	let request = $state({ ...defaults });

	// Reset form to default values
	function reset() {
		request = { ...defaults };
	}

	// Modal open/close state
	let open = $state(false);

	// Close modal function
	function close() {
		open = false;
	}

	// Load available architectures for the current distro series
	onMount(async () => {
		try {
			await client.listBootImageSelections({}).then((response) => {
				// Find the boot image selection for the current distro series
				const currentBootImageSelection = response.bootImageSelections.find(
					(bootImageSelection) => bootImageSelection.distroSeries === bootImage.distroSeries
				);

				if (currentBootImageSelection) {
					architecturesOptions.set(
						currentBootImageSelection.architectures.map((architecture) => ({
							value: architecture,
							label: architecture,
							icon: 'ph:empty'
						}))
					);
				}
			});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<!-- Modal component for boot image edit -->
<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		{m.edit()}
	</Modal.Trigger>

	<!-- Modal content -->
	<Modal.Content>
		<Modal.Header>{m.edit_boot_image()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.architecture()}</Form.Label>
					<MultipleSelect.Root bind:value={request.architectures} options={architecturesOptions}>
						<MultipleSelect.Viewer />
						<MultipleSelect.Controller>
							<MultipleSelect.Trigger />
							<MultipleSelect.Content>
								<MultipleSelect.Options>
									<MultipleSelect.Input />
									<MultipleSelect.List>
										<MultipleSelect.Empty>{m.no_result()}</MultipleSelect.Empty>
										<MultipleSelect.Group>
											{#each $architecturesOptions as option}
												<MultipleSelect.Item {option}>
													<Icon
														icon={option.icon ? option.icon : 'ph:empty'}
														class={cn('size-5', option.icon ? 'visible' : 'invisible')}
													/>
													{option.label}
													<MultipleSelect.Check {option} />
												</MultipleSelect.Item>
											{/each}
										</MultipleSelect.Group>
									</MultipleSelect.List>
									<MultipleSelect.Actions>
										<MultipleSelect.ActionAll>{m.all()}</MultipleSelect.ActionAll>
										<MultipleSelect.ActionClear>{m.clear()}</MultipleSelect.ActionClear>
									</MultipleSelect.Actions>
								</MultipleSelect.Options>
							</MultipleSelect.Content>
						</MultipleSelect.Controller>
					</MultipleSelect.Root>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>

		<!-- Modal footer with action buttons -->
		<Modal.Footer>
			<!-- Cancel button -->
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>

			<!-- Confirm action group -->
			<Modal.ActionsGroup>
				<!-- Confirm button with edit operation -->
				<Modal.Action
					onclick={() => {
						const architectures = `${request.architectures.join(', ')}`;
						// Execute edit operation with toast notifications
						toast.promise(() => client.updateBootImage(request), {
							loading: `Editing boot image ${request.distroSeries}...`,
							success: () => {
								// Force reload to refresh data
								reloadManager.force();
								return `Successfully edited boot image ${request.distroSeries}: ${architectures}`;
							},
							error: (error) => {
								// Handle and display error
								let message = `Failed to edit boot image ${request.distroSeries}: ${architectures}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						// Reset form and close modal
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
