<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { Plus } from '@lucide/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import {
		type GroupedFields,
		MultiStepSchemaForm,
		UserSelectWidget
	} from '$lib/components/custom/schema-form';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Dialog from '$lib/components/ui/dialog';

	let {
		schema: apiSchema
	}: {
		schema?: any;
	} = $props();

	const cluster = $derived(page.params.cluster ?? page.params.scope ?? '');

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	const initialData = {
		spec: {
			resourceQuota: {
				hard: {
					'requests.cpu': '16',
					'requests.memory': '32Gi',
					'requests.otterscale.com/vgpu': '0',
					'requests.otterscale.com/vgpumem': '0',
					'requests.otterscale.com/vgpumem-percentage': '0',
					'limits.cpu': '16',
					'limits.memory': '32Gi'
				}
			},
			limitRange: {
				limits: [
					{
						type: 'Container',
						default: {
							cpu: '1',
							memory: '2Gi'
						},
						defaultRequest: {
							cpu: '1',
							memory: '2Gi'
						}
					}
				]
			},
			users: [
				{
					subject: page.data.user?.sub,
					name: `${page.data.user?.name} (${page.data.user?.email || page.data.user?.username})`,
					role: 'admin'
				}
			]
		}
	};

	const groupedFields: GroupedFields = {
		// Step 1: Workspace & Users
		'Workspace & Users': {
			'metadata.name': { title: 'Workspace Name' },
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
		// Step 2: Network Isolation
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
		// Step 3: Default Resource Settings (read-only with preset values)
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

	function transformFormData(data: Record<string, unknown>) {
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

		// Enforce fixed LimitRange
		if (spec) {
			spec.limitRange = {
				limits: [
					{
						type: 'Container',
						default: {
							cpu: '500m',
							memory: '512Mi'
						},
						defaultRequest: {
							cpu: '500m',
							memory: '512Mi'
						}
					}
				]
			};
		}

		return data;
	}

	let isSubmitting = $state(false);
	async function handleMultiStepSubmit(data: Record<string, unknown>) {
		if (isSubmitting) return;
		isSubmitting = true;

		// Construct the full resource object
		const resourceObject = {
			apiVersion: 'tenant.otterscale.io/v1alpha1',
			kind: 'Workspace',
			...data
		};

		const name = (data.metadata as { name: string })?.name;

		toast.promise(
			async () => {
				const manifest = new TextEncoder().encode(JSON.stringify(resourceObject));

				await resourceClient.create({
					cluster,
					group: 'tenant.otterscale.io',
					version: 'v1alpha1',
					resource: 'workspaces',
					manifest
				});

				return resourceObject;
			},
			{
				loading: `Creating workspace ${name}...`,
				success: () => {
					return `Successfully created workspace ${name}`;
				},
				error: (err) => {
					console.error('Failed to create workspace:', err);
					return `Failed to create workspace: ${(err as ConnectError).message}`;
				},
				finally() {
					isSubmitting = false;
					open = false;
				}
			}
		);
	}

	let open = $state(false);
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger>
		<Button variant="outline" size="icon">
			<Plus />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="min-h-[77vh] min-w-[50vw]">
		<MultiStepSchemaForm
			{apiSchema}
			fields={groupedFields}
			{initialData}
			title="Create Workspace"
			onSubmit={handleMultiStepSubmit}
			transformData={transformFormData}
		/>
	</Dialog.Content>
</Dialog.Root>
