<script lang="ts">
	import { goto } from '$app/navigation';
	import { EssentialService } from '$gen/api/essential/v1/essential_pb';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Application from '$lib/components/dashboard/application/index.svelte';
	import PageLoading from '$lib/components/otterscale/ui/page-loading.svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onMount } from 'svelte';
	import type { Writable } from 'svelte/store';

	const transport: Transport = getContext('transport');
	const prometheusDriver: Writable<PrometheusDriver> = getContext('prometheusDriver');

	let scopes: Scope[] = $state([]);
	async function fetchJujuModelUuids() {
		const query = `group by (juju_model_uuid) (node_exporter_build_info{})`;

		try {
			const response = await $prometheusDriver.instantQuery(query);

			scopes = response.result
				.map((result) => {
					return {
						uuid: result.metric.labels.juju_model_uuid,
						name: result.metric.labels.juju_model_uuid
					} as Scope;
				})
				.reverse();
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let mounted = $state(false);
	onMount(async () => {
		try {
			const essentialClient = createClient(EssentialService, transport);
			const listEssentialsResponse = await essentialClient.listEssentials({});
			if (listEssentialsResponse.essentials?.length === 0) {
				goto('/dashboard/installation');
			}

			await fetchJujuModelUuids();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if !mounted}
	<PageLoading />
{:else}
	<Application client={$prometheusDriver} {scopes} />
{/if}
