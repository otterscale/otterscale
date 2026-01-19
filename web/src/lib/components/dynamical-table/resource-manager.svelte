<script lang="ts">
	import { type JsonValue, toJson } from '@bufbuild/protobuf';
	import { StructSchema } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Download from '@lucide/svelte/icons/download';
	import Plus from '@lucide/svelte/icons/plus';
	import { getContext, onDestroy, onMount } from 'svelte';

	import {
		type ListRequest,
		ResourceService,
		type SchemaRequest,
		WatchEvent_Type,
		type WatchRequest
	} from '$lib/api/resource/v1/resource_pb';
	import DynamicalTable from '$lib/components/dynamical-table/dynamical-table.svelte';
	import Button from '$lib/components/ui/button/button.svelte';

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
			if (error instanceof Error && error.name === 'AbortError') {
				console.log('List aborted');
				return;
			}

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
			if (error instanceof Error && error.name === 'AbortError') {
				console.log('Watch stream aborted');
				return;
			}

			console.error('Failed to watch resources:', error);
		} finally {
			isWatching = false;
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
			listAbortController.abort();
			listAbortController = null;
		}
		if (watchAbortController) {
			watchAbortController.abort();
			watchAbortController = null;
		}
	});

	function handleReload() {
		if (!isWatching) {
			watchResources();
		}
	}
</script>

{#if isMounted}
	<DynamicalTable {objects} {fields}>
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
