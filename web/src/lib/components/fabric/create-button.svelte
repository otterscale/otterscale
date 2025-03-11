<script lang="ts">
	import Icon from '@iconify/svelte';

	import { buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';

	import { FabricCreateOverview, FabricCreateConnector } from '$lib/components/fabric';
	import { listConnectors, type pbConnector } from '$lib/pb';
	import type { Connector } from '$lib/components/fabric/connector';
	import FabricCreatePipeline from '$lib/components/fabric/fabric-create-pipeline.svelte';
	import { onMount } from 'svelte';
	import { cn } from '$lib/utils';

	let sources: Connector[] = [
		{
			key: 'postgresql',
			name: 'PostgreSQL',
			icon: 'logos:postgresql',
			parameters: [
				{
					key: 'connection_string',
					name: 'Connection String',
					value: '',
					description: `connection string, such as 'postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10'`
				},
				{
					key: 'namespace',
					name: 'Namespace',
					value: '',
					description: `namespace of database, such as 'public'`
				}
			],
			templates: [
				{
					name: 'TEST',
					parameters: [
						{
							key: 'connection_string',
							value:
								'postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10'
						},
						{
							key: 'namespace',
							value: 'public'
						}
					]
				}
			]
		},
		{
			key: 'csv',
			name: 'CSV',
			icon: 'ph:file-csv',
			parameters: [
				{
					key: 'file_path',
					name: 'File Path',
					value: '',
					description: 'csv file path'
				},
				{
					key: 'table_name',
					name: 'Table Name',
					value: '',
					description: 'destination table name'
				},
				{
					key: 'inferring',
					name: 'Inferring',
					value: '',
					description: ''
				}
			]
		}
	];

	let destinations: Connector[] = [
		{
			key: 'postgresql',
			name: 'PostgreSQL',
			icon: 'logos:postgresql',
			parameters: [
				{
					key: 'connection_string',
					name: 'Connection String',
					value: '',
					description: `connection string, such as 'postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10'`
				},
				{
					key: 'namespace',
					name: 'Namespace',
					value: '',
					description: `namespace of database, such as 'public'`
				}
			],
			extraParameters: [
				{
					key: 'batch_size',
					name: 'Batch Size',
					value: '',
					description: 'default batch size of rows is 10,000 if not specified'
				},
				{
					key: 'batch_size_bytes',
					name: 'Batch Size Bytes',
					value: '',
					description: 'default batch size of bytes is 10,000,000 bytes if not specified'
				},
				{
					key: 'batch_timeout',
					name: 'Batch Timeout',
					value: '',
					description: 'default batch timeout is 60s if not specified'
				},
				{
					key: 'create_index',
					name: 'Create Index',
					value: 'true',
					description: 'create an index to improve performance'
				}
			]
		}
	];

	let items = $state([
		{
			name: 'Source',
			icon: 'ph:plug',
			active: false
		},
		{
			name: 'Destination',
			icon: 'ph:plugs',
			active: false
		},
		{
			name: 'Pipeline',
			icon: 'ph:plugs-connected',
			active: false
		}
	]);

	let open = $state(false);

	let pbSources: pbConnector[] = $state([]);
	let pbDestinations: pbConnector[] = $state([]);

	onMount(async () => {
		pbSources = await listConnectors(`kind='source'`);
		pbDestinations = await listConnectors(`kind='destination'`);
	});
	//
</script>

<Dialog.Root
	bind:open
	onOpenChange={(open) => {
		if (!open) {
			setTimeout(() => {
				items = items.map((item) => ({ ...item, active: false }));
			}, 100);
		}
	}}
>
	<Dialog.Trigger class={cn(buttonVariants({ size: 'icon' }), '[&_svg]:size-5')}>
		<Icon icon="ph:plus" />
	</Dialog.Trigger>
	<Dialog.Content class="max-w-2xl">
		<Dialog.Header class="flex-col space-y-8 py-4">
			<Dialog.Title class="flex">
				{#each items as item}
					{#if item.active}
						<div class="flex items-center pl-2">
							<Icon icon={item.icon} class="size-8" />
							<div class="space-x-2 pl-2">{item.name}</div>
						</div>
					{/if}
				{/each}
				{#if items.filter((item) => item.active).length === 0}
					<div class="flex w-full items-center justify-center">What do you need?</div>
				{/if}
			</Dialog.Title>
			<Dialog.Description class="flex w-full justify-center px-2">
				{#if items[0].active}
					<FabricCreateConnector bind:parent={open} items={sources} />
				{:else if items[1].active}
					<FabricCreateConnector bind:parent={open} items={destinations} />
				{:else if items[2].active}
					<FabricCreatePipeline
						bind:parent={open}
						sources={pbSources}
						destinations={pbDestinations}
					/>
				{:else}
					<FabricCreateOverview bind:items />
				{/if}
			</Dialog.Description>
		</Dialog.Header>
	</Dialog.Content>
</Dialog.Root>
