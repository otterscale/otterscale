<script lang="ts">
	import { type JsonValue, toJson } from '@bufbuild/protobuf';
	import { StructSchema } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Cable, Unplug } from '@lucide/svelte';
	import type { ColumnDef } from '@tanstack/table-core';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type APIResource,
		type ListRequest,
		ResourceService,
		type SchemaRequest,
		WatchEvent_Type,
		type WatchRequest
	} from '$lib/api/resource/v1/resource_pb';
	import { DynamicTable } from '$lib/components/dynamic-table';
	import Button from '$lib/components/ui/button/button.svelte';

	import type { DataSchemaType, UISchemaType } from '../dynamic-table/utils';
	import type { ActionsType, CreateType } from './kind-viewer-actions';
	import { getActions, getCreate } from './kind-viewer-actions';
	import {
		getColumnDefinitions,
		getDataSchemas,
		getDataSet as getData,
		getUISchemas
	} from './kind-viewers';

	let {
		clustered,
		apiResource,
		cluster,
		group,
		version,
		kind,
		resource,
		namespace
	}: {
		clustered: boolean;
		apiResource: APIResource;
		cluster: string;
		group: string;
		version: string;
		kind: string;
		resource: string;
		namespace?: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	const uiSchemas: Record<string, UISchemaType> = $derived(getUISchemas(kind));
	const dataSchemas: Record<string, DataSchemaType> = $derived(getDataSchemas(kind));

	// eslint-disable-next-line
	let schema: any = $state({});

	async function fetchSchema() {
		try {
			const schemaResponse = await resourceClient.schema({
				cluster: cluster,
				group: group,
				version: version,
				kind: kind
			} as SchemaRequest);

			schema = toJson(StructSchema, schemaResponse);
		} catch (error) {
			console.error('Failed to fetch schema:', error);
			return null;
		}
	}

	let dataset: Record<string, JsonValue>[] = $state([]);
	let columnDefinitions: ColumnDef<Record<string, JsonValue>>[] | undefined = $state(undefined);

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
						namespace: clustered ? undefined : apiResource.namespaced ? namespace : undefined,
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

				const newData = response.items.map((item) => getData(kind, item.object));
				dataset = [...dataset, ...newData];

				if (listAbortController.signal.aborted) {
					break;
				}
			} while (continueToken);
		} catch (error) {
			if (listAbortController.signal.aborted) return;

			console.error('Failed to list resources:', error);

			return null;
		} finally {
			isListing = false;
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
					namespace: clustered ? undefined : apiResource.namespaced ? namespace : undefined,
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
					const addedData = getData(kind, response.resource?.object);

					const index = dataset.findIndex((object) => {
						if (apiResource.namespaced) {
							return object.Namespace === addedData.Namespace && object.Name === addedData.Name;
						} else {
							return object.Name === addedData.Name;
						}
					});

					if (index < 0) {
						dataset = [...dataset, addedData];
					}
				} else if (response.type === WatchEvent_Type.MODIFIED) {
					const modifiedData = getData(kind, response.resource?.object);

					dataset = dataset.map((object) => {
						if (
							apiResource.namespaced &&
							object.Namespace === modifiedData.Namespace &&
							object.Name === modifiedData.Name
						) {
							return modifiedData;
						} else if (!apiResource.namespaced && object.Name === modifiedData.Name) {
							return modifiedData;
						}

						return object;
					});
				} else if (response.type === WatchEvent_Type.DELETED) {
					const deletedData = getData(kind, response.resource?.object);
					dataset = dataset.filter((object) => {
						if (apiResource.namespaced) {
							return object.Namespace === deletedData.Namespace && object.Name !== deletedData.Name;
						} else {
							return object.Name !== deletedData.Name;
						}
					});
				} else {
					console.log('Unknown response type: ', response);
				}
			}
		} catch (error) {
			if (watchAbortController.signal.aborted) return;

			console.error('Failed to watch resources:', error);
		} finally {
			toast.info(`Watching resource ${namespace} ${resource} was cancelled.`);

			isWatching = false;
			watchAbortController = null;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		await fetchSchema();
		await listResources();
		columnDefinitions = getColumnDefinitions(kind, apiResource, uiSchemas, dataSchemas);
		watchResources();

		isMounted = true;
	});

	let isDestroyed = false;
	onDestroy(() => {
		isDestroyed = true;

		if (listAbortController) {
			listAbortController.abort();
		}
		if (watchAbortController) {
			watchAbortController.abort();
		}
	});

	function handleReload() {
		if (!isWatching) {
			watchResources();
		}
	}

	const Create: CreateType = $derived(getCreate(kind));
	const Actions: ActionsType = $derived(getActions(kind));
</script>

{#if isMounted}
	{#if columnDefinitions}
		<DynamicTable {dataset} {columnDefinitions} {uiSchemas} {dataSchemas}>
			{#snippet create()}
				<Create {schema} />
			{/snippet}
			{#snippet reload()}
				<Button onclick={handleReload} disabled={isWatching} variant="outline">
					{#if isWatching}
						<Cable class="opacity-60" size={16} />
					{:else}
						<Unplug class="text-destructive opacity-60" size={16} />
					{/if}
				</Button>
			{/snippet}
			{#snippet rowActions({ row })}
				<Actions {row} {schema} object={row.original.raw} />
			{/snippet}
		</DynamicTable>
	{/if}
{/if}
