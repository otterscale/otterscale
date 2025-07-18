<script lang="ts">
	import { DataTable } from './data-table';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { BISTService } from '$gen/api/bist/v1/bist_pb';
	import { getContext } from 'svelte';
	import { PageLoading } from '$lib/components/otterscale/ui';

	// grpc
	const transport: Transport = getContext('transport');
	const bistClient = createClient(BISTService, transport);
</script>

{#await bistClient.listTestResults({})}
	<PageLoading />
{:then response}
	{@const testResults = response.testResults.filter((result) => result.kind.case === 'fio' )}
	<DataTable testResults={testResults} />
{:catch e}
	<div class="flex w-fill items-center justify-center border">No Data</div>
{/await}