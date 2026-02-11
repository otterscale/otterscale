<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { Pencil } from '@lucide/svelte';
	import type { TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import {
		type GroupedFields,
		MultiStepSchemaForm,
		UserSelectWidget
	} from '$lib/components/custom/schema-form';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Item from '$lib/components/ui/item';

	let {
		schema: apiSchema,
		object,
		onOpenChangeComplete
	}: {
		schema: any;
		object: Record<string, unknown>;
		onOpenChangeComplete: () => void;
	} = $props();

	const cluster = $derived(page.params.cluster ?? page.params.scope ?? '');

	let open = $state(false);

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	function getCleanedObject() {
		const copy = structuredClone($state.snapshot(object) as Record<string, unknown>);
		if (copy.metadata && typeof copy.metadata === 'object') {
			delete (copy.metadata as Record<string, unknown>).managedFields;
		}
		return copy;
	}

	const fields: GroupedFields = {
		'Workspace & Users': {
			'metadata.name': { title: 'Workspace Name', disabled: true }, // Name is immutable
			'spec.users': {
				title: 'Users',
				uiSchema: {
					items: {
						'ui:components': {
							objectField: UserSelectWidget
						}
					}
				}
			}
		},
		'Network Isolation': {
			'spec.networkIsolation': { title: 'Network Isolation' },
			'spec.networkIsolation.enabled': {
				title: 'Enable Network Isolation',
				uiSchema: {
					'ui:components': {
						checkboxWidget: 'switchWidget'
					}
				}
			},
			'spec.networkIsolation.allowedNamespaces': { title: 'Allowed Namespaces' }
		},
		'Default Resource Settings': {
			'spec.resourceQuota.hard.requests.cpu': { title: 'Requests CPU' },
			'spec.resourceQuota.hard.requests.memory': { title: 'Requests Memory' },
			'spec.resourceQuota.hard.requests.otterscale.com/vgpu': {
				title: 'Requests GPU',
				disabled: true
			},
			'spec.resourceQuota.hard.requests.otterscale.com/vgpumem': {
				title: 'Requests GPU Memory',
				disabled: true
			},
			'spec.resourceQuota.hard.requests.otterscale.com/vgpumem-percentage': {
				title: 'Requests GPU Memory Percentage',
				disabled: true
			}
		}
	};

	function transfer(data: Record<string, unknown>) {
		const spec = data.spec as Record<string, any>;
		const metadata = data.metadata as Record<string, any>;

		// Set namespace to be the same as workspace name
		if (spec && metadata?.name) {
			spec.namespace = metadata.name;
		}

		// Handle Resource Quota Logic: limits align with requests, strict defaults
		if (spec?.resourceQuota?.hard) {
			const hard = spec.resourceQuota.hard;
			// Sync limits with requests
			if (hard['requests.cpu']) hard['limits.cpu'] = hard['requests.cpu'];
			if (hard['requests.memory']) hard['limits.memory'] = hard['requests.memory'];
		}

		return data;
	}

	let isSubmitting = $state(false);
	async function handleSubmit(data: Record<string, unknown>) {
		if (isSubmitting) return;
		isSubmitting = true;

		const resourceObject: Record<string, any> = {
			apiVersion: 'tenant.otterscale.io/v1alpha1',
			kind: 'Workspace',
			...data
		};

		if (!resourceObject.metadata) resourceObject.metadata = {};
		(resourceObject.metadata as any).name = (
			object as TenantOtterscaleIoV1Alpha1Workspace
		)?.metadata?.name;

		toast.promise(
			async () => {
				const manifest = new TextEncoder().encode(JSON.stringify(resourceObject));

				await resourceClient.apply({
					cluster,
					name: (object as TenantOtterscaleIoV1Alpha1Workspace)?.metadata?.name,
					group: 'tenant.otterscale.io',
					version: 'v1alpha1',
					resource: 'workspaces',
					manifest,
					fieldManager: 'otterscale-web-ui',
					force: true
				});
			},
			{
				loading: `Updating workspace ${
					(object as TenantOtterscaleIoV1Alpha1Workspace)?.metadata?.name
				}...`,
				success: () => {
					return `Successfully updated workspace ${
						(object as TenantOtterscaleIoV1Alpha1Workspace)?.metadata?.name
					}`;
				},
				error: (err) => {
					console.error('Failed to update workspace:', err);
					return `Failed to update workspace: ${(err as ConnectError).message}`;
				},
				finally: () => {
					isSubmitting = false;
					open = false;
				}
			}
		);
	}
</script>

<Dialog.Root bind:open {onOpenChangeComplete}>
	<Dialog.Trigger>
		<Item.Root class="p-0" size="sm">
			<Item.Media>
				<Pencil />
			</Item.Media>
			<Item.Content>
				<Item.Description>Update</Item.Description>
			</Item.Content>
		</Item.Root>
	</Dialog.Trigger>
	<Dialog.Content class="min-h-[77vh] min-w-[50vw]">
		<MultiStepSchemaForm
			{apiSchema}
			{fields}
			initialData={getCleanedObject()}
			title="Edit Workspace"
			onSubmit={handleSubmit}
			transformData={transfer}
		/>
	</Dialog.Content>
</Dialog.Root>
