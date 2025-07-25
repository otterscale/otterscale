<script lang="ts">
	import { MachineService } from '$gen/api/machine/v1/machine_pb';
	import Hardware from '$lib/components/dashboard/hardware/index.svelte';
	import PageLoading from '$lib/components/otterscale/ui/page-loading.svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import type { Writable } from 'svelte/store';
	import { getContext } from 'svelte';

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const prometheusDriver: Writable<PrometheusDriver> = getContext('prometheusDriver');
</script>

{#await machineClient.listMachines({})}
	<PageLoading />
{:then response}
	{@const machines = response.machines.filter(
		(result) =>
			result.workloadAnnotations['juju-machine-id'] &&
			result.workloadAnnotations['juju-machine-id'].includes('-machine-')
	)}
	<Hardware client={$prometheusDriver} machines={machines} />
{:catch e}
	<div class="flex w-fill items-center justify-center border">No Data</div>
{/await}
