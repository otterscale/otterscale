<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import type { DoSMARTResponse_Output, OSD } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Code from '$lib/components/custom/code';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
	import { currentCeph } from '$lib/stores';
</script>

<script lang="ts">
	let { osd }: { osd: OSD } = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
	let open = $state(false);
	let smarts = $state(writable<Record<string, DoSMARTResponse_Output>>({}));
	let isSMARTsLoading = $state(true);

	async function fetchSMARTs() {
		try {
			const response = await storageClient.doSMART({
				scope: $currentCeph?.scope,
				facility: $currentCeph?.name,
				osdName: osd.name,
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

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:file" />
		{m.do_smart()}
	</Modal.Trigger>
	<Modal.Content class="min-w-[50vw]">
		<Modal.Header>
			{m.smart()}
		</Modal.Header>
		{#if !isMounted}
			<Loading.Report />
		{:else}
			{#each Object.entries($smarts) as [device, output]}
				{@const result = output.lines.join('\n')}
				<p class="text-lg font-bold">
					{device}
				</p>
				<Code.Root class="h-fit max-h-[77vh] w-full overflow-auto" code={result} hideLines>
					<Code.CopyButton />
				</Code.Root>
			{/each}
		{/if}
		<Modal.Footer>
			<Modal.Cancel>{m.cancel()}</Modal.Cancel>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
