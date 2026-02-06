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

	import { type ActionsType, getActions } from './actions';
	import { type CreatorType, getCreator } from './creators';
	import { getColumnDefinitionsGetter, getFieldsGetter, getObjectGetter } from './viewers';

	let {
		apiResource,
		cluster,
		group,
		version,
		kind,
		resource,
		namespace
	}: {
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

	const getFields = getFieldsGetter(kind);
	const getObject = getObjectGetter(kind);
	const getColumnDefinitions = getColumnDefinitionsGetter(kind);

	// eslint-disable-next-line
	let schema: any = $state({});
	let fields: Record<string, { description: string; type: string; format?: string }> = $state({});
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
						namespace: apiResource.namespaced ? namespace : undefined,
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

				const newObjects = response.items.map((item) => getObject(item.object));
				objects = [...objects, ...newObjects];

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
					namespace: apiResource.namespaced ? namespace : undefined,
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
					const addedObject = getObject(response.resource?.object);

					const index = objects.findIndex(
						(object) =>
							object.Namespace === addedObject.Namespace && object.Name === addedObject.Name
					);

					if (index < 0) {
						objects = [...objects, addedObject];
					}
				} else if (response.type === WatchEvent_Type.MODIFIED) {
					const modifiedObject = getObject(response.resource?.object);

					objects = objects.map((object) =>
						object.Namespace === modifiedObject.Namespace && object.Name === modifiedObject.Name
							? modifiedObject
							: object
					);
				} else if (response.type === WatchEvent_Type.DELETED) {
					const deletedObject = getObject(response.resource?.object);

					objects = objects.filter(
						(object) =>
							object.Namespace === deletedObject.Namespace && object.Name !== deletedObject.Name
					);
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
		columnDefinitions = getColumnDefinitions(apiResource, fields);
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

	const Creator: CreatorType = $derived(getCreator(kind));
	const Actions: ActionsType = $derived(getActions(kind));
</script>

{#if isMounted}
	{#if columnDefinitions}
		<DynamicTable {objects} {fields} {columnDefinitions}>
			{#snippet create()}
				<Creator {schema} />
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
			{#snippet rowActions({ row, fields, objects })}
				<Actions {row} schema={fields.raw} object={objects[row.id].raw} />
			{/snippet}
		</DynamicTable>
	{/if}
{/if}
