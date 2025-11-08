<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { ScaleApplicationRequest } from '$lib/api/application/v1/application_pb';
	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';

	import type { Application } from '../types';
</script>

<script lang="ts">
	// Component props - accepts an Application object
	let { application }: { application: Application } = $props();

	// Get required services from Svelte context
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	// Create gRPC client for application operations
	const applicationClient = createClient(ApplicationService, transport);

	// Form validation state
	let invalid = $state(false);

	// Default values for the scale application request
	const defaults = {
		scope: $currentKubernetes?.scope || '',
		facility: $currentKubernetes?.name || '',
		name: application.name,
		namespace: application.namespace,
		type: application.type,
		replicas: application.replicas
	} as ScaleApplicationRequest;

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
</script>

<!-- Modal component for application scaling -->
<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:resize" />
		{m.scale()}
	</Modal.Trigger>

	<!-- Modal content -->
	<Modal.Content>
		<Modal.Header>{m.scale_application()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.replica()}</Form.Label>
					<SingleInput.General
						type="number"
						required
						bind:value={request.replicas}
						bind:invalid
						min={0}
						max={100}
					/>
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
				<!-- Confirm button with extend operation -->
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						// Execute scale operation with toast notifications
						toast.promise(() => applicationClient.scaleApplication(request), {
							loading: `Scaling ${request.name} to ${request.replicas} replicas...`,
							success: () => {
								// Force reload to refresh data
								reloadManager.force();
								return `Successfully scaled ${request.name} to ${request.replicas} replicas`;
							},
							error: (error) => {
								// Handle and display error
								let message = `Failed to scale ${request.name}`;
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
