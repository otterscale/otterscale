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

	import type { Application } from '../types';
</script>

<script lang="ts">
	// Component props - accepts an Application object
	let {
		application,
		scope,
		reloadManager
	}: { application: Application; scope: string; reloadManager: ReloadManager } = $props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	let request = $state({} as ScaleApplicationRequest);
	let invalid = $state(false);
	let open = $state(false);

	function init() {
		request = {
			scope: scope,
			name: application.name,
			namespace: application.namespace,
			type: application.type,
			replicas: application.replicas
		} as ScaleApplicationRequest;
	}

	// Close modal function
	function close() {
		open = false;
	}
</script>

<!-- Modal component for application scaling -->
<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
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
			<Modal.Cancel>
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
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
