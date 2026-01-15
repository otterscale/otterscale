<script lang="ts">
	import { type JsonValue, toJson } from '@bufbuild/protobuf';
	import { StructSchema } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Download from '@lucide/svelte/icons/download';
	import { getContext, onDestroy, onMount } from 'svelte';

	import {
		type ListRequest,
		ResourceService,
		type SchemaRequest,
		WatchEvent_Type,
		type WatchRequest
	} from '$lib/api/resource/v1/resource_pb';
	import { Button } from '$lib/components/ui/button/index.js';

	import DynamicalTable from './components/dynamical-table.svelte';

	const configuration = {
		cluster: 'gpu',
		group: 'batch',
		namespace: 'llm-d',
		version: 'v1',
		kind: 'Job',
		resource: 'jobs'
	};
	// eslint-disable-next-line
	function getFields(schema: any): Record<string, JsonValue> {
		return {
			Name: schema?.properties?.metadata?.properties?.name ?? {},
			Namespace: schema?.properties?.metadata?.properties?.namespace ?? {},
			CreationTimestamp: schema?.properties?.metadata?.properties?.creationTimestamp ?? {},
			Succeeded: schema?.properties?.status?.properties?.succeeded ?? {},
			Terminating: schema?.properties?.status?.properties?.terminating ?? {},
			StartTime: schema?.properties?.status?.properties?.startTime ?? {},
			CompletionTime: schema?.properties?.status?.properties?.completionTime ?? {},
			Suspend: schema?.properties?.spec?.properties?.suspend ?? {},
			Object: schema ?? {}
		};
	}
	// eslint-disable-next-line
	function getObject(object: any): Record<string, JsonValue> {
		return {
			Name: object?.metadata?.name ?? null,
			Namespace: object?.metadata?.namespace ?? null,
			CreationTimestamp: object?.metadata?.creationTimestamp ?? null,
			Succeeded: object?.status?.succeeded ?? null,
			Terminating: object?.status?.terminating ?? null,
			StartTime: object?.status?.startTime ?? null,
			CompletionTime: object?.status?.completionTime ?? null,
			Suspend: object?.spec?.suspend ?? null,
			Object: object ?? null
		};
	}

	const transport: Transport = getContext('transport');
	const resourceService = createClient(ResourceService, transport);

	// eslint-disable-next-line
	let schema: any = $state({});
	let fields: Record<string, JsonValue> = $state({});
	async function fetchSchema() {
		try {
			const schemaResponse = await resourceService.schema({
				version: configuration.version,
				group: configuration.group,
				cluster: configuration.cluster,
				kind: configuration.kind
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
				const response = await resourceService.list(
					{
						version: configuration.version,
						group: configuration.group,
						cluster: configuration.cluster,
						namespace: configuration.namespace,
						resource: configuration.resource,
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
			const watchResourcesStream = resourceService.watch(
				{
					version: configuration.version,
					group: configuration.group,
					cluster: configuration.cluster,
					namespace: configuration.namespace,
					resource: configuration.resource,
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
	<DynamicalTable {configuration} {schema} {objects} {fields}>
		{#snippet reloader()}
			<Button onclick={handleReload} disabled={isWatching} variant="outline">
				<Download class="opacity-60" size={16} />
				Reload
			</Button>
		{/snippet}
	</DynamicalTable>
{/if}
