<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { type FormState, getValueSnapshot } from '@sjsf/form';
	import yaml from 'js-yaml';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import {
		type K8sOpenAPISchema,
		type PathOptions,
		SchemaForm
	} from '$lib/components/custom/schema-form';
	import Button from '$lib/components/ui/button/button.svelte';

	import vmSchema from './vm_api.json';

	let form = $state<FormState<Record<string, unknown>> | undefined>();
	let mode = $state<'basic' | 'advance'>('basic');

	const transport = getContext<Transport>('transport');
	const resourceClient = createClient(ResourceService, transport);

	async function handleCreate() {
		if (!formValues) {
			toast.error('Form is empty');
			return;
		}

		// Ensure we have the GVK
		const manifestObj = {
			apiVersion: 'kubevirt.io/v1',
			kind: 'VirtualMachine',
			...(formValues && typeof formValues === 'object' ? formValues : {})
		};

		// Construct the manifest bytes
		let manifestYaml: string;
		try {
			manifestYaml = yaml.dump(manifestObj);
		} catch (e) {
			console.error('YAML dump error:', e);
			toast.error('Failed to serialize manifest');
			return;
		}

		const encoder = new TextEncoder();
		const manifestBytes = encoder.encode(manifestYaml);

		// Get Name and Namespace from formValues (metadata)
		const metadata = (formValues.metadata as Record<string, any>) || {};
		const name = metadata.name as string;
		const namespace = (metadata.namespace as string) || 'default';

		if (!name) {
			toast.error('Resource name is required');
			return;
		}

		try {
			console.log('Request Body:', {
				cluster: 'aaa', // TODO: Make dynamic if needed
				group: 'kubevirt.io',
				version: 'v1',
				resource: 'virtualmachines',
				namespace: namespace,
				manifest: manifestYaml // Logging string yaml instead of bytes for readability
			});
			/*
			await resourceClient.create({
				cluster: 'aaa', // TODO: Make dynamic if needed
				group: 'kubevirt.io',
				version: 'v1',
				resource: 'virtualmachines',
				namespace: namespace,
				manifest: manifestBytes
			});
			*/
			toast.success(`VirtualMachine ${name} created successfully (Dry Run)`);
		} catch (error) {
			console.error('Create error:', error);
			toast.error(`Failed to create VirtualMachine: ${(error as any).message}`);
		}
	}

	const fields: Record<string, PathOptions> = {
		'metadata.name': { title: 'Name', required: true },
		'metadata.namespace': { title: 'Namespace' },
		'spec.runStrategy': { title: 'Run Strategy' },

		// Instance Type (Use predefined instance types)
		'spec.instancetype.name': { title: 'Instance Type', required: true },
		'spec.instancetype.kind': { title: 'Instance Type Kind' },

		// Preference
		'spec.preference.name': { title: 'Preference' },
		'spec.preference.kind': { title: 'Preference Kind' },

		// Simple Disks (Assume one boot disk + maybe data disk)
		'spec.template.spec.domain.devices.disks': { title: 'Disks' },
		'spec.template.spec.domain.devices.disks.name': { title: 'Name' },
		'spec.template.spec.domain.devices.disks.bootOrder': { title: 'Boot Order' },
		// 'spec.template.spec.domain.devices.disks.disk.bus': { title: 'Disk Bus' }, // Often inferred or standard

		// Volumes (Where the data comes from)
		'spec.template.spec.volumes': { title: 'Volumes' },
		'spec.template.spec.volumes.name': { title: 'Name' },
		// 'spec.template.spec.volumes.containerDisk.image': { title: 'Container Image (Boot)' }, // Disable containerDisk
		'spec.template.spec.volumes.persistentVolumeClaim.claimName': { title: 'PVC Name (Boot/Data)' }, // Only allow PVC
		'spec.template.spec.volumes.cloudInitNoCloud': { title: 'Cloud Init Config' },

		// Data Volume Templates (Dynamic Provisioning)
		'spec.dataVolumeTemplates': { title: 'New Data Volumes' },
		'spec.dataVolumeTemplates.metadata.name': { title: 'Name' },
		// 'spec.dataVolumeTemplates.spec.pvc.storageClassName': { title: 'Storage Class' },
		'spec.dataVolumeTemplates.spec.pvc.resources.requests.storage': { title: 'Size' },
		'spec.dataVolumeTemplates.spec.source.http.url': { title: 'HTTP Source URL' },
		'spec.dataVolumeTemplates.spec.source.registry.url': { title: 'Registry Source URL' },

		// Networking (Default to standard pod network usually)
		'spec.template.spec.networks': { title: 'Networks' },
		'spec.template.spec.networks.name': { title: 'Network Name' },
		// 'spec.template.spec.networks.pod': { title: 'Pod Network' }, // Default

		// Interfaces (Matched with Networks)
		'spec.template.spec.domain.devices.interfaces': { title: 'Interfaces' },
		'spec.template.spec.domain.devices.interfaces.name': { title: 'Interface Name' }
		// 'spec.template.spec.domain.devices.interfaces.masquerade': { title: 'Masquerade' }, // Default for KubeVirt

		// Access
		// 'spec.template.spec.accessCredentials': { title: 'Access Credentials' },
	};

	const formValues = $derived(form ? getValueSnapshot(form) : {});
	const fieldKeys = Object.keys(fields);
	// Default values to simplify the form (Hidden from Basic View but present in Data)
	const initialData = {
		spec: {
			runStrategy: 'Always', // Default to running
			instancetype: {
				kind: 'VirtualMachineClusterInstancetype' // Default kind
			},
			preference: {
				kind: 'VirtualMachineClusterPreference' // Default kind
			},
			template: {
				spec: {
					domain: {
						devices: {
							interfaces: [{ name: 'default', masquerade: {} }],
							disks: [{ name: 'boot', bootOrder: 1, disk: { bus: 'virtio' } }] // Default boot disk config
						}
					},
					networks: [{ name: 'default', pod: {} }] // Default pod network
				}
			}
		}
	};
</script>

<div class="container mx-auto py-10">
	<h1 class="mb-4 text-2xl font-bold">Schema Form Gen Experiment</h1>

	<div class="grid grid-cols-2 gap-8">
		<div class="rounded border bg-card p-4 text-card-foreground">
			<h2 class="mb-4 text-xl">Generated Form (Mode: {mode})</h2>
			<SchemaForm
				apiSchema={vmSchema as K8sOpenAPISchema}
				paths={fields}
				{initialData}
				bind:form
				bind:mode
				onSubmit={handleCreate}
			/>
		</div>

		<div class="rounded border bg-muted/50 p-4">
			<h2 class="mb-4 text-xl">Live Values</h2>
			<pre class="overflow-auto rounded bg-zinc-950 p-4 text-xs text-zinc-50 dark:bg-zinc-900">
{JSON.stringify(formValues, null, 2)}
            </pre>

			<h2 class="mt-4 mb-2 text-xl">Selected Paths</h2>
			<ul class="list-inside list-disc text-sm">
				{#each fieldKeys as field (field)}
					<li><code>{field}</code></li>
				{/each}
			</ul>
		</div>
	</div>
</div>
