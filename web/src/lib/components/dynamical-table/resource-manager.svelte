<script lang="ts">
	import { type JsonValue, toJson } from '@bufbuild/protobuf';
	import { StructSchema } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import CircleAlert from '@lucide/svelte/icons/circle-alert';
	import Download from '@lucide/svelte/icons/download';
	import Trash from '@lucide/svelte/icons/trash';
	import type { Table } from '@tanstack/table-core';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type ListRequest,
		ResourceService,
		type SchemaRequest,
		WatchEvent_Type,
		type WatchRequest
	} from '$lib/api/resource/v1/resource_pb';
	import DynamicalTable from '$lib/components/dynamical-table/dynamical-table.svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import Button from '$lib/components/ui/button/button.svelte';

	import Separator from '../ui/separator/separator.svelte';
	import ResourceActions from './resource-actions.svelte';
	import ResourceCreate from './resource-create.svelte';

	let {
		cluster,
		group,
		version,
		kind,
		resource,
		namespace
	}: {
		cluster: string;
		group: string;
		version: string;
		kind: string;
		resource: string;
		namespace?: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	// eslint-disable-next-line
	function getFields(schema: any): Record<string, JsonValue> {
		return {
			Name: schema?.properties?.metadata?.properties?.name ?? {},
			Namespace: schema?.properties?.metadata?.properties?.namespace ?? {},
			Labels: schema?.properties?.metadata?.properties?.labels ?? {},
			Annotations: schema?.properties?.metadata?.properties?.annotations ?? {},
			CreateTime: schema?.properties?.metadata?.properties?.creationTimestamp ?? {},
			Configuration: schema ?? {}
		};
	}
	// eslint-disable-next-line
	function getObject(object: any, fields: Record<string, JsonValue>): Record<string, JsonValue> {
		return {
			Name: object?.metadata?.name ?? null,
			Namespace: object?.metadata?.namespace ?? null,
			Labels: object?.metadata?.labels ?? null,
			Annotations: object?.metadata?.annotations ?? null,
			CreateTime: object?.metadata?.creationTimestamp ?? null,
			Configuration: object ?? null
		};
	}

	// eslint-disable-next-line
	let schema: any = $state({});
	let fields: Record<string, JsonValue> = $state({});
	async function fetchSchema() {
		try {
			const schemaResponse = await resourceClient.schema({
				cluster: cluster,
				group: group,
				version: version,
				kind: kind
			} as SchemaRequest);

			schema = toJson(StructSchema, schemaResponse);
			fields = getFields(schema);
		} catch (error) {
			console.error('Failed to fetch schema:', error);
			return null;
		}
	}

	let objects: Record<string, JsonValue>[] = $state([]);

	let listAbortController: AbortController | null = null;
	let watchAbortController: AbortController | null = null;

	let resourceVersion: string | undefined = $state(undefined);

	let isListing = $state(false);
	async function listResources() {
		if (isListing || isWatching || isDestroyed) return;

		isListing = true;
		listAbortController = new AbortController();
		try {
			let continueToken: string | undefined = undefined;
			do {
				const response = await resourceClient.list(
					{
						cluster: cluster,
						namespace: namespace,
						group: group,
						version: version,
						resource: resource,
						limit: BigInt(10),
						continue: continueToken
					} as ListRequest,
					{ signal: listAbortController.signal }
				);

				resourceVersion = response.resourceVersion;
				continueToken = response.continue;

				const newObjects = response.items.map((item) => getObject(item.object, fields));
				objects = [...objects, ...newObjects];

				if (listAbortController.signal.aborted) {
					break;
				}
			} while (continueToken);
		} catch (error) {
			if (error instanceof Error && error.name === 'ConnectError') {
				if (error.cause === 'Aborted due to component destroyed.') {
					return;
				}
			}

			console.error('Failed to list resources:', error);

			return null;
		} finally {
			isListing = false;
			if (listAbortController?.signal.aborted) {
				toast.info('Resource listing was cancelled.');
			} else {
				toast.warning('Resource listing has stopped unexpectedly. Please try again.');
			}
			listAbortController = null;
		}
	}

	let isWatching = $state(false);
	async function watchResources() {
		if (isListing || isWatching || isDestroyed) return;

		isWatching = true;
		watchAbortController = new AbortController();
		try {
			const watchResourcesStream = resourceClient.watch(
				{
					cluster: cluster,
					namespace: namespace,
					group: group,
					version: version,
					resource: resource,
					resourceVersion: resourceVersion
				} as WatchRequest,
				{ signal: watchAbortController.signal }
			);

			for await (const watchResourcesResponse of watchResourcesStream) {
				// eslint-disable-next-line
				const response: any = watchResourcesResponse;

				if (response.type === WatchEvent_Type.ERROR) {
					continue;
				}

				if (response.type === WatchEvent_Type.BOOKMARK) {
					resourceVersion = response.resourceVersion as string;
					continue;
				}

				resourceVersion = response.resource?.object?.metadata?.resourceVersion;

				if (response.type === WatchEvent_Type.ADDED) {
					const addedObject = getObject(response.resource?.object, fields);

					const index = objects.findIndex(
						(object) =>
							object.Namespace === addedObject.Namespace && object.Name === addedObject.Name
					);

					if (index < 0) {
						objects = [...objects, addedObject];
					}
				} else if (response.type === WatchEvent_Type.MODIFIED) {
					const modifiedObject = getObject(response.resource?.object, fields);

					objects = objects.map((object) =>
						object.Namespace === modifiedObject.Namespace && object.Name === modifiedObject.Name
							? modifiedObject
							: object
					);
				} else if (response.type === WatchEvent_Type.DELETED) {
					const deletedObject = getObject(response.resource?.object, fields);

					objects = objects.filter(
						(object) =>
							object.Namespace === deletedObject.Namespace && object.Name !== deletedObject.Name
					);
				} else {
					console.log('Unknown response type: ', response);
				}
			}
		} catch (error) {
			if (error instanceof Error && error.name === 'ConnectError') {
				if (error.cause === 'Aborted due to component destroyed.') {
					return;
				}
			}

			console.error('Failed to watch resources:', error);
		} finally {
			isWatching = false;
			if (watchAbortController?.signal.aborted) {
				toast.info('Resource watching was cancelled.');
			} else {
				toast.warning('Resource watching has stopped unexpectedly. Please try again.', {
					action: {
						label: 'Reload',
						onClick: handleReload
					},
					duration: 13 * 1000
				});
			}
			watchAbortController = null;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchSchema();
			await listResources();
			watchResources();

			isMounted = true;
		} catch (error) {
			console.error('Failed to initialize:', error);
		}
	});

	let isDestroyed = false;
	onDestroy(() => {
		isDestroyed = true;
		if (listAbortController) {
			listAbortController.abort('Aborted due to component destroyed.');
			listAbortController = null;
		}
		if (watchAbortController) {
			watchAbortController.abort('Aborted due to component destroyed.');
			watchAbortController = null;
		}
	});
	function handleDeleteRows(table: Table<Record<string, JsonValue>>) {
		const selectedRows = table.getSelectedRowModel().rows;
		objects = objects.filter(
			(object) =>
				!selectedRows.some((row) => row.original && object && row.original.id === object.id)
		);
		table.resetRowSelection();
	}
	function handleReload() {
		if (!isWatching) {
			watchResources();
		}
	}
</script>

{#if isMounted}
	<DynamicalTable {objects} {fields}>
		{#snippet bulkDelete({ table })}
			{#if table.getSelectedRowModel().rows.length > 0}
				<AlertDialog.Root>
					<AlertDialog.Trigger>
						{#snippet child({ props })}
							<Button class="ml-auto text-destructive" variant="outline" {...props}>
								<Trash class="-ms-1 opacity-60" size={16} aria-hidden="true" />
								Delete
								<Separator orientation="vertical" />
								{table.getSelectedRowModel().rows.length}
							</Button>
						{/snippet}
					</AlertDialog.Trigger>
					<AlertDialog.Content>
						<div class="flex flex-col gap-2 max-sm:items-center sm:flex-row sm:gap-4">
							<div
								class="flex size-9 shrink-0 items-center justify-center rounded-full border"
								aria-hidden="true"
							>
								<CircleAlert class="opacity-80" size={16} />
							</div>
							<AlertDialog.Header>
								<AlertDialog.Title>Are you absolutely sure?</AlertDialog.Title>
								<AlertDialog.Description>
									This action cannot be undone. This will permanently delete
									{table.getSelectedRowModel().rows.length} selected
									{table.getSelectedRowModel().rows.length === 1 ? 'row' : 'rows'}.
								</AlertDialog.Description>
							</AlertDialog.Header>
						</div>
						<AlertDialog.Footer>
							<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
							<AlertDialog.Action onclick={() => handleDeleteRows(table)}>Delete</AlertDialog.Action
							>
						</AlertDialog.Footer>
					</AlertDialog.Content>
				</AlertDialog.Root>
			{/if}
		{/snippet}
		{#snippet create()}
			<ResourceCreate {resource} />
		{/snippet}
		{#snippet reload()}
			<Button onclick={handleReload} disabled={isWatching} variant="outline">
				<Download class="opacity-60" size={16} />
				Reload
			</Button>
		{/snippet}
		{#snippet rowActions({ row })}
			<ResourceActions {row} />
		{/snippet}
	</DynamicalTable>
{/if}
