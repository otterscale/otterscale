<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Nexus, type Facility_Info, type Scope } from '$gen/api/nexus/v1/nexus_pb';
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import { ManagementApplicationController } from '$lib/components/otterscale/index';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.listScopes({});
			scopesStore.set(response.scopes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}

	const kubernetesestore = writable<Facility_Info[]>();
	const kubernetesesLoading = writable(true);
	async function fetchKuberneteses(scopeUuid: string) {
		try {
			const response = await client.listKuberneteses({
				scopeUuid: scopeUuid
			});
			kubernetesestore.set(response.kuberneteses);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			kubernetesesLoading.set(false);
		}
	}

	let kuberneteses = [] as Facility_Info[];

	let mounted = false;
	onMount(async () => {
		try {
			await fetchScopes();
			for (const scope of $scopesStore) {
				await fetchKuberneteses(scope.uuid);
				$kubernetesestore.forEach((k) => {
					kuberneteses.push(k);
				});
			}
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<ManagementApplicationController {kuberneteses} />
{:else}
	<PageLoading />
{/if}
