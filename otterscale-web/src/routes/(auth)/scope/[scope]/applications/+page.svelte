<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { Dashboard } from '$lib/components/applications/dashboard';
	import Loading from '$lib/components/custom/loading/report.svelte';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onMount } from 'svelte';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.applications(page.params.scope) });

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let prometheusDriver: PrometheusDriver | null = null;
	let mounted = false;

	async function initializePrometheusDriver(): Promise<PrometheusDriver> {
		const response = await environmentService.getPrometheus({});
		return new PrometheusDriver({
			endpoint: `${env.PUBLIC_API_URL}/prometheus`,
			baseURL: response.baseUrl
		});
	}

	onMount(async () => {
		try {
			prometheusDriver = await initializePrometheusDriver();
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		} finally {
			mounted = true;
		}
	});
</script>

{#if mounted && prometheusDriver && $activeScope}
	<Dashboard client={prometheusDriver} scope={$activeScope} />
{:else}
	<!-- <div class="flex items-center justify-center p-8">
		<Icon icon="mdi:loading" class="animate-spin text-2xl" />
		<span class="ml-2">Loading Storage...</span>
	</div> -->
	<Loading />
{/if}
