<script lang="ts">
	import { toJson } from '@bufbuild/protobuf';
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
	import {
		getColumnDefinitionsDynamically,
		getFieldsDynamically,
		getValuesDynamically
	} from './viewers';
	import type { FieldsType, ValuesType } from './type';

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

	// eslint-disable-next-line
	let schema: any = $state({});
	let fields: FieldsType = $state({});
	async function fetchSchema() {
		try {
			const schemaResponse = await resourceClient.schema({
				cluster: cluster,
				group: group,
				version: version,
				kind: kind
			} as SchemaRequest);

			schema = toJson(StructSchema, schemaResponse);
			fields = getFieldsDynamically(kind, schema);
		} catch (error) {
			console.error('Failed to fetch schema:', error);
			return null;
		}
	}

	let dataset: ValuesType[] = $state([]);
	let columnDefinitions: ColumnDef<ValuesType>[] | undefined = $state(undefined);

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

				const newObjects = response.items.map((item) => getValuesDynamically(kind, item.object));
				dataset = [...dataset, ...newObjects];

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
					const addedData = getValuesDynamically(kind, response.resource?.object);

					const index = dataset.findIndex(
						(object) => object.Namespace === addedData.Namespace && object.Name === addedData.Name
					);

					if (index < 0) {
						dataset = [...dataset, addedData];
					}
				} else if (response.type === WatchEvent_Type.MODIFIED) {
					const modifiedData = getValuesDynamically(kind, response.resource?.object);

					dataset = dataset.map((object) =>
						object.Namespace === modifiedData.Namespace && object.Name === modifiedData.Name
							? modifiedData
							: object
					);
				} else if (response.type === WatchEvent_Type.DELETED) {
					const deletedData = getValuesDynamically(kind, response.resource?.object);
					dataset = dataset.filter(
						(object) =>
							object.Namespace === deletedData.Namespace && object.Name !== deletedData.Name
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
		columnDefinitions = getColumnDefinitionsDynamically(kind, apiResource, fields);
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
		<DynamicTable {dataset} {fields} {columnDefinitions}>
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
			{#snippet rowActions({ row, fields, dataset })}
				{@const objects: ValuesType = dataset[Number(row.id)]}
				<Actions {row} schema={fields['Configuration']} object={objects['Configuration']} />
			{/snippet}
		</DynamicTable>
	{/if}
{/if}
