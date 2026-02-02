<script lang="ts">
	import { toJson } from '@bufbuild/protobuf';
	import { StructSchema } from '@bufbuild/protobuf/wkt';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import {
		type GroupedFields,
		type K8sOpenAPISchema,
		MultiStepSchemaForm
	} from '$lib/components/custom/schema-form';

	let {
		onsuccess
	}: {
		onsuccess?: (cronjob?: any) => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);
	const cluster = $derived(page.params.cluster ?? page.params.scope ?? '');

	let apiSchema: K8sOpenAPISchema | undefined = $state();
	let isSubmitting = $state(false);

	// Default values for CronJob
	const initialData = {
		spec: {
			schedule: '0 0 * * *',
			concurrencyPolicy: 'Allow',
			suspend: false,
			jobTemplate: {
				spec: {
					template: {
						spec: {
							restartPolicy: 'OnFailure',
							containers: [
								{
									name: 'hello',
									image: 'busybox',
									args: ['/bin/sh', '-c', 'date; echo Hello from the Kubernetes cluster']
								}
							]
						}
					}
				}
			}
		}
	};

	// Grouped fields for multi-step form
	const groupedFields: GroupedFields = {
		// Step 1: General Settings
		'General Settings': {
			'metadata.name': { title: 'Name' },
			'spec.namespace': { title: 'Namespace', showDescription: true },
			'spec.schedule': { title: 'Schedule', description: 'Cron schedule string, e.g. "0 0 * * *"' },
			'spec.concurrencyPolicy': { title: 'Concurrency Policy' },
			'spec.suspend': { title: 'Suspend execution' }
		},
		// Step 2: Job Settings
		'Job Settings': {
			'spec.jobTemplate.spec.template.spec.restartPolicy': { title: 'Restart Policy' },
			'spec.jobTemplate.spec.template.spec.containers': {
				title: 'Containers'
				// Allowing default array UI for containers
			}
		}
	};

	function transformFormData(data: Record<string, unknown>) {
		// Pass through data for CronJob
		return data;
	}

	async function handleMultiStepSubmit(data: Record<string, unknown>) {
		if (isSubmitting) return;
		isSubmitting = true;

		// Construct the full resource object
		const resourceObject = {
			apiVersion: 'batch/v1',
			kind: 'CronJob',
			...data
		};

		const name = (data.metadata as { name: string })?.name;

		toast.promise(
			async () => {
				const manifest = new TextEncoder().encode(JSON.stringify(resourceObject));

				await resourceClient.create({
					cluster,
					group: 'batch',
					version: 'v1',
					resource: 'cronjobs',
					manifest
				});

				return resourceObject;
			},
			{
				loading: `Creating cronjob ${name}...`,
				success: (cronjob) => {
					isSubmitting = false;
					onsuccess?.(cronjob);
					return `Successfully created cronjob ${name}`;
				},
				error: (err) => {
					isSubmitting = false;
					console.error('Failed to create cronjob:', err);
					return `Failed to create cronjob: ${(err as ConnectError).message}`;
				}
			}
		);
	}

	onMount(async () => {
		try {
			const res = await resourceClient.schema({
				cluster,
				group: 'batch',
				version: 'v1',
				kind: 'CronJob'
			});
			// Convert Protobuf Struct to plain JSON object
			apiSchema = toJson(StructSchema, res) as K8sOpenAPISchema;
		} catch (err) {
			console.error('Failed to fetch cronjob schema:', err);
			toast.error(`Failed to fetch cronjob schema: ${(err as ConnectError).message}`);
		}
	});
</script>

<div class="h-full w-full">
	{#if apiSchema}
		<MultiStepSchemaForm
			{apiSchema}
			fields={groupedFields}
			{initialData}
			title="Create CronJob"
			onSubmit={handleMultiStepSubmit}
			transformData={transformFormData}
		/>
	{:else}
		<div class="flex h-32 items-center justify-center">
			<p class="text-muted-foreground">Loading schema...</p>
		</div>
	{/if}
</div>
