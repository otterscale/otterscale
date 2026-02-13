<script lang="ts">
	import type { ConnectError } from '@connectrpc/connect';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/stores';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import { type GroupedFields, MultiStepSchemaForm } from '$lib/components/custom/schema-form';

	type Props = {
		schema?: any;
		onsuccess?: (simpleapp?: Record<string, any>) => void;
	};

	let { schema, onsuccess }: Props = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);
	const cluster = $derived($page.params.cluster ?? $page.params.scope ?? '');

	const initialData = {
		metadata: {
			name: '',
			namespace: $page.url.searchParams.get('namespace') ?? 'default'
		},
		spec: {
			deploymentSpec: {
				replicas: 1,
				selector: {
					matchLabels: {}
				},
				template: {
					metadata: {
						labels: {}
					},
					spec: {
						containers: [
							{
								name: 'app',
								image: '',
								command: ['bin/bash', 'sleep 3600s'],
								ports: [
									{
										containerPort: 8080,
										protocol: 'TCP'
									}
								],
								resources: {
									requests: {
										memory: '128Mi',
										cpu: '100m'
									},
									limits: {
										memory: '256Mi',
										cpu: '200m'
									}
								}
							}
						]
					}
				}
			},
			serviceSpec: {
				type: 'ClusterIP',
				selector: {},
				ports: [
					{
						port: 80,
						targetPort: '8080',
						protocol: 'TCP'
					}
				]
			},
			pvcSpec: {
				accessModes: ['ReadWriteOnce'],
				resources: {
					requests: {
						storage: '1Gi'
					}
				},
				storageClassName: ''
			}
		}
	};

	const groupedFields: GroupedFields = {
		General: {
			'metadata.name': { title: 'Name' },
			'metadata.namespace': { title: 'Namespace' },
			'spec.deploymentSpec.replicas': { title: 'Replicas' }
		},
		Container: {
			'spec.deploymentSpec.template.spec.containers': {
				title: 'Containers',
				uiSchema: {
					items: {
						'ui:title': 'Container',
						'ui:options': {
							addable: true,
							removable: true,
							orderable: true
						}
					}
				}
			},
			'spec.deploymentSpec.template.spec.containers.name': { title: 'Container Name' },
			'spec.deploymentSpec.template.spec.containers.image': { title: 'Image' },
			'spec.deploymentSpec.template.spec.containers.command': { title: 'Command' },
			'spec.deploymentSpec.template.spec.containers.ports': {
				title: 'Ports',
				uiSchema: {
					items: {
						'ui:title': 'Container Port'
					}
				}
			},
			'spec.deploymentSpec.template.spec.containers.ports.containerPort': {
				title: 'Container Port'
			},
			'spec.deploymentSpec.template.spec.containers.ports.protocol': { title: 'Protocol' },
			'spec.deploymentSpec.template.spec.containers.resources.requests.memory': {
				title: 'Memory Request'
			},
			'spec.deploymentSpec.template.spec.containers.resources.requests.cpu': {
				title: 'CPU Request'
			},
			'spec.deploymentSpec.template.spec.containers.resources.limits.memory': {
				title: 'Memory Limit'
			},
			'spec.deploymentSpec.template.spec.containers.resources.limits.cpu': {
				title: 'CPU Limit'
			}
		},
		Port: {
			'spec.serviceSpec.type': { title: 'Service Type' },
			'spec.serviceSpec.ports.port': { title: 'Service Port' },
			'spec.serviceSpec.ports.targetPort': {
				title: 'Target Port',
				uiSchema: {
					'ui:options': {
						inputType: 'number'
					}
				}
			}
		},
		Storage: {
			'spec.pvcSpec.accessModes': { title: 'Access Mode' },
			'spec.pvcSpec.resources.requests.storage': { title: 'Storage Size' },
			'spec.pvcSpec.storageClassName': { title: 'StorageClass Name' }
		}
	};

	function transformFormData(data: Record<string, unknown>): Record<string, unknown> {
		const transformed = { ...data };

		// Automatically set labels to match selector
		if (transformed.metadata && (transformed.metadata as any).name) {
			const name = (transformed.metadata as any).name;
			const labels = { app: name };

			// Set deployment selector and template labels
			if (transformed.spec && (transformed.spec as any).deploymentSpec) {
				const deploymentSpec = (transformed.spec as any).deploymentSpec;
				deploymentSpec.selector = { matchLabels: labels };
				if (deploymentSpec.template && deploymentSpec.template.metadata) {
					deploymentSpec.template.metadata.labels = labels;
				}
			}

			// Set service selector
			if (transformed.spec && (transformed.spec as any).serviceSpec) {
				(transformed.spec as any).serviceSpec.selector = labels;
			}
		}

		// Convert targetPort to appropriate type (Kubernetes intOrString)
		// - Number or numeric string → convert to integer
		// - Named port string like 'http' → keep as string
		if (transformed.spec && (transformed.spec as any).serviceSpec?.ports) {
			const ports = (transformed.spec as any).serviceSpec.ports;
			if (Array.isArray(ports)) {
				ports.forEach((port: any) => {
					if (port.targetPort !== undefined) {
						const numValue = parseInt(port.targetPort, 10);
						if (!isNaN(numValue) && String(numValue) === port.targetPort.trim()) {
							port.targetPort = numValue;
						}
						// If already a number, keep it as number
					}
				});
			}
		}

		return transformed;
	}

	let isSubmitting = $state(false);

	async function handleMultiStepSubmit(data: Record<string, unknown>) {
		if (isSubmitting) return;
		isSubmitting = true;

		const name = (data.metadata as any)?.name || 'simpleapp';

		const resourceObject = {
			apiVersion: 'apps.otterscale.io/v1alpha1',
			kind: 'SimpleApp',
			...data
		};

		toast.promise(
			async () => {
				const manifest = new TextEncoder().encode(JSON.stringify(resourceObject));
				await resourceClient.create({
					cluster,
					namespace: $page.url.searchParams.get('namespace') ?? '',
					group: 'apps.otterscale.io',
					version: 'v1alpha1',
					resource: 'simpleapps',
					manifest
				});
				return resourceObject;
			},
			{
				loading: `Creating simpleapp ${name}...`,
				success: (simpleapp) => {
					isSubmitting = false;
					onsuccess?.(simpleapp);
					return `Successfully created simpleapp ${name}`;
				},
				error: (err) => {
					isSubmitting = false;
					console.error('Failed to create simpleapp:', err);
					return `Failed to create simpleapp: ${(err as ConnectError).message}`;
				}
			}
		);
	}
</script>

<div class="h-full w-full">
	{#if schema}
		<MultiStepSchemaForm
			apiSchema={schema}
			fields={groupedFields}
			{initialData}
			title="Create SimpleApp"
			onSubmit={handleMultiStepSubmit}
			transformData={transformFormData}
			yamlEditable={true}
		/>
	{:else}
		<div class="flex h-32 items-center justify-center">
			<p class="text-muted-foreground">Loading schema...</p>
		</div>
	{/if}
</div>
