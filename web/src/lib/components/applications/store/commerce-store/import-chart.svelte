<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { type Chart, RegistryService } from '$lib/api/registry/v1/registry_pb';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		scope,
		charts
	}: {
		scope: string;
		charts: Writable<Chart[]>;
	} = $props();

	const transport: Transport = getContext('transport');
	const registryClient = createClient(RegistryService, transport);

	let open = $state(false);
	let chartRef = $state('');

	async function handleImport() {
		if (!chartRef) {
			toast.error('Please enter a chart link');
			return;
		}

		const promise = async () => {
			open = false;

			const getRegistryURLResponse = await registryClient.getRegistryURL({
				scope: scope
			});

			await registryClient.importChart({
				chartRef: chartRef,
				registryUrl: getRegistryURLResponse.registryUrl
			});

			const response = await registryClient.listCharts({
				scope: scope
			});

			charts.set(response.charts.sort((p, n) => p.name.localeCompare(n.name)));

			chartRef = '';
		};

		toast.promise(promise(), {
			loading: 'Importing chart...',
			success: 'Chart imported successfully!',
			error: (error) => {
				console.error('Failed to import chart:', error);
				const message = error instanceof ConnectError ? error.message : 'Unknown error';
				return `Failed to import chart: ${message}`;
			}
		});
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}>
		<Icon icon="ph:download-simple" />
		{m.import()}
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Import Chart</Dialog.Title>
			<Dialog.Description>
				Enter the direct link to the .tgz file of the chart you wish to import.
			</Dialog.Description>
		</Dialog.Header>
		<div class="flex items-center gap-2">
			<div class="grid flex-1 gap-2">
				<Label for="chart-ref" class="sr-only">Link</Label>
				<Input
					id="chart-ref"
					placeholder="https://charts.bitnami.com/bitnami/nginx-15.0.0.tgz"
					bind:value={chartRef}
					onkeydown={(e) => {
						if (e.key === 'Enter') handleImport();
					}}
				/>
			</div>
		</div>
		<Dialog.Footer>
			<Button type="submit" onclick={handleImport}>
				{m.confirm()}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
