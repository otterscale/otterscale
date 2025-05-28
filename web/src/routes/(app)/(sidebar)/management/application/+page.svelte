<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ScopeService, type Scope } from '$gen/api/scope/v1/scope_pb';
	import {
		Essential_Type,
		EssentialService,
		type Essential
	} from '$gen/api/essential/v1/essential_pb';
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import { ManagementApplicationController } from '$lib/components/otterscale/index';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const essentialClient = createClient(EssentialService, transport);

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopesStore.set(response.scopes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}

	const kubernetesestore = writable<Essential[]>();
	const kubernetesesLoading = writable(true);
	async function fetchKuberneteses(scopeUuid: string) {
		try {
			const response = await essentialClient.listEssentials({
				type: Essential_Type.KUBERNETES,
				scopeUuid: scopeUuid
			});
			kubernetesestore.set(response.essentials);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			kubernetesesLoading.set(false);
		}
	}

	let kuberneteses = [] as Essential[];

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
