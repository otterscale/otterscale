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
		'metadata.annotations': { title: 'Annotations' },
		'spec.runStrategy': { title: 'Run Strategy' },
		'spec.instancetype.name': { title: 'Instance Type Name' },
		'spec.instancetype.kind': { title: 'Instance Type Kind' },
		'spec.preference.name': { title: 'Preference Name' },
		'spec.template.spec.domain.cpu.cores': { title: 'CPU Cores' },

		// Disks
		'spec.template.spec.domain.devices.disks': { title: 'Disks' },
		'spec.template.spec.domain.devices.disks.name': { title: 'Name' },
		'spec.template.spec.domain.devices.disks.bootOrder': { title: 'Boot Order' },
		'spec.template.spec.domain.devices.disks.disk.bus': { title: 'Disk Bus' },
		'spec.template.spec.domain.devices.disks.cdrom.bus': { title: 'CD-ROM Bus' },

		// Volumes
		'spec.template.spec.volumes': { title: 'Volumes' },
		'spec.template.spec.volumes.name': { title: 'Name' },
		'spec.template.spec.volumes.containerDisk.image': { title: 'Container Disk Image' },
		'spec.template.spec.volumes.persistentVolumeClaim.claimName': { title: 'PVC Name' },
		'spec.template.spec.volumes.dataVolume.name': { title: 'Data Volume Name' },
		'spec.template.spec.volumes.cloudInitNoCloud': { title: 'Cloud Init NoCloud' },

		// Networks
		'spec.template.spec.networks': { title: 'Networks' },
		'spec.template.spec.networks.name': { title: 'Name' },
		'spec.template.spec.networks.multus.networkName': { title: 'Multus Network Name' },
		'spec.template.spec.networks.pod': { title: 'Pod Network' },

		// Interfaces
		'spec.template.spec.domain.devices.interfaces': { title: 'Interfaces' },
		'spec.template.spec.domain.devices.interfaces.name': { title: 'Name' },
		'spec.template.spec.domain.devices.interfaces.model': { title: 'Model' },
		'spec.template.spec.domain.devices.interfaces.masquerade': { title: 'Masquerade' },
		'spec.template.spec.domain.devices.interfaces.bridge': { title: 'Bridge' },

		// Data Volume Templates
		'spec.dataVolumeTemplates': { title: 'Data Volume Templates' },
		'spec.dataVolumeTemplates.metadata.name': { title: 'Name' },
		'spec.dataVolumeTemplates.spec.pvc.storageClassName': { title: 'Storage Class' },
		'spec.dataVolumeTemplates.spec.pvc.resources.requests.storage': { title: 'Size (e.g. 10Gi)' },
		'spec.dataVolumeTemplates.spec.source.http.url': { title: 'HTTP Source URL' },
		'spec.dataVolumeTemplates.spec.source.registry.url': { title: 'Registry Source URL' },

		// Access Credentials
		'spec.template.spec.accessCredentials': { title: 'Access Credentials' },
		'spec.template.spec.accessCredentials.sshPublicKey.source.secret.secretName': {
			title: 'SSH Key Secret'
		},
		'spec.template.spec.accessCredentials.sshPublicKey.propagationMethod.noCloud': {
			title: 'Propagate via NoCloud'
		},
		'spec.template.spec.accessCredentials.sshPublicKey.propagationMethod.qemuGuestAgent': {
			title: 'Propagate via QEMU Agent'
		}
	};

	// Using $derived to reactively get values from the form store
	const formValues = $derived(form ? getValueSnapshot(form) : {});
	// const initialData = { spec: { runStrategy: '123' } };
	const initialData = {};
	const fieldKeys = Object.keys(fields);
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
