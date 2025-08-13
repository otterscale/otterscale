<script lang="ts" module>
	import type { DoSMARTResponse_Output, OSD } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Loading from '$lib/components/custom/loading';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import * as Code from '$lib/components/custom/code';
	import { writable } from 'svelte/store';
</script>

<script lang="ts">
	let { osd }: { osd: OSD } = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
	const stateController = new StateController(false);

	let smarts = $state(writable<Record<string, DoSMARTResponse_Output>>({}));
	let isSMARTsLoading = $state(true);

	async function fetchSMARTs() {
		try {
			const response = await storageClient.doSMART({
				scopeUuid: 'b62d195e-3905-4960-85ee-7673f71eb21e',
				facilityName: 'ceph-mon',
				osdName: osd.name
			});
			smarts.set(response.deviceOutputMap);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isSMARTsLoading = false;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchSMARTs();
			if (!isSMARTsLoading) {
				isMounted = true;
			}
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:file" />
		Do SMART
	</AlertDialog.Trigger>
	<AlertDialog.Content class="min-w-[50vw]">
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold">
			S.M.A.R.T.
		</AlertDialog.Header>
		<Form.Root>
			{#if !isMounted}
				<Loading.Report />
			{:else}
				{#each Object.entries($smarts) as [device, output]}
					{@const result = output.lines.join('\n')}
					<Form.Fieldset>
						<Form.Legend>
							{device}
						</Form.Legend>
						<Form.Field class="gap-1">
							<Code.Root class="w-full" code={result} hideLines>
								<Code.CopyButton />
							</Code.Root>
						</Form.Field>
					</Form.Fieldset>
				{/each}
			{/if}
		</Form.Root>

		<AlertDialog.Footer>
			<AlertDialog.Cancel>Close</AlertDialog.Cancel>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
