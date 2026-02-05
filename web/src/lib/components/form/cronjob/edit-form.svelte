<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import {
		type GroupedFields,
		type K8sOpenAPISchema,
		MultiStepSchemaForm
	} from '$lib/components/custom/schema-form';

	let {
		name,
		schema,
		object,
		onsuccess
	}: {
		name: string;
		schema: K8sOpenAPISchema;
		object: Record<string, unknown>;
		onsuccess?: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);
	const cluster = $derived(page.params.cluster ?? page.params.scope ?? '');

	// Remove metadata.managedFields from object
	function getCleanedObject() {
		const copy = structuredClone($state.snapshot(object) as Record<string, unknown>);
		if (copy.metadata && typeof copy.metadata === 'object') {
			delete (copy.metadata as Record<string, unknown>).managedFields;
		}
		return copy;
	}

	let isSubmitting = $state(false);

	// Grouped fields for multi-step form
	const groupedFields: GroupedFields = {
		// Step 1: General Settings
		General: {
			'metadata.name': { title: 'Name' },
			'spec.schedule': { title: 'Schedule', showDescription: true },
			'spec.concurrencyPolicy': { title: 'Concurrency Policy' },
			'spec.suspend': {
				title: 'Suspend execution',
				uiSchema: {
					'ui:components': {
						checkboxWidget: 'switchWidget'
					}
				}
			}
		},
		// Step 2: Container Settings
		Container: {
			'spec.jobTemplate.spec.template.spec.containers.name': { title: 'Name' },
			'spec.jobTemplate.spec.template.spec.containers.image': { title: 'Image' },
			'spec.jobTemplate.spec.template.spec.containers.command': { title: 'Command' },
			'spec.jobTemplate.spec.template.spec.containers.args': { title: 'Arguments' },
			'spec.jobTemplate.spec.template.spec.containers.env': { title: 'Environment Variables' }
		},
		// Step 3: Resources
		Resources: {
			'spec.jobTemplate.spec.template.spec.containers.resources.requests.cpu': {
				title: 'Requests CPU'
			},
			'spec.jobTemplate.spec.template.spec.containers.resources.requests.memory': {
				title: 'Requests Memory'
			},
			'spec.jobTemplate.spec.template.spec.containers.resources.limits.cpu': {
				title: 'Limits CPU'
			},
			'spec.jobTemplate.spec.template.spec.containers.resources.limits.memory': {
				title: 'Limits Memory'
			}
		}
	};

	function transformFormData(data: Record<string, unknown>) {
		return data;
	}

	async function handleMultiStepSubmit(data: Record<string, unknown>) {
		if (isSubmitting) return;
		isSubmitting = true;

		// Construct the full resource object
		const resourceObject: Record<string, unknown> = {
			apiVersion: 'batch/v1',
			kind: 'CronJob',
			...data
		};

		// Ensure name is correct
		if (!resourceObject.metadata) resourceObject.metadata = {};
		(resourceObject.metadata as Record<string, unknown>).name = name;

		toast.promise(
			async () => {
				const manifest = new TextEncoder().encode(JSON.stringify(resourceObject));

				await resourceClient.apply({
					cluster,
					name,
					namespace: page.url.searchParams.get('namespace') ?? '',
					group: 'batch',
					version: 'v1',
					resource: 'cronjobs',
					manifest,
					fieldManager: 'otterscale-web-ui',
					force: true
				});
			},
			{
				loading: `Updating cronjob ${name}...`,
				success: () => {
					isSubmitting = false;
					onsuccess?.();
					return `Successfully updated cronjob ${name}`;
				},
				error: (err) => {
					isSubmitting = false;
					console.error('Failed to update cronjob:', err);
					return `Failed to update cronjob: ${(err as ConnectError).message}`;
				}
			}
		);
	}
</script>

<div class="h-full w-full">
	<MultiStepSchemaForm
		apiSchema={schema}
		fields={groupedFields}
		initialData={getCleanedObject()}
		title={`Edit CronJob: ${name}`}
		onSubmit={handleMultiStepSubmit}
		transformData={transformFormData}
		yamlEditable={true}
	/>
</div>
