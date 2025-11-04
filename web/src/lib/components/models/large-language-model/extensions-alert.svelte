<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { writable, type Writable } from 'svelte/store';

	import { OrchestratorService, type Extension } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { Single as Alert } from '$lib/components/custom/alert';
	import { installExtensions } from '$lib/components/settings/extensions/utils.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const modelExtensions: Writable<Extension[]> = writable([]);

	orchestratorClient
		.listModelExtensions({ scope: scope, facility: facility })
		.then((response) => {
			modelExtensions.set(response.Extensions);
		})
		.catch((error) => {
			console.error('Failed to fetch extensions:', error);
		});

	const alert: Alert.AlertType = $derived({
		title: m.extensions_alert_title(),
		message: m.extensions_alert_description(),
		action: () => {
			installExtensions(['model']);
		},
		variant: 'destructive',
	});
</script>

{#if $modelExtensions.filter((modelExtension) => modelExtension.current).length < $modelExtensions.length}
	<Alert.Root {alert}>
		<Alert.Icon />
		<Alert.Title>{alert.title}</Alert.Title>
		<Alert.Description>
			<div class="space-y-1">
				<p>{alert.message}</p>
				<div class="flex w-full flex-wrap gap-2">
					{#each $modelExtensions.filter((extension) => !extension.current) as extension}
						<Badge variant="outline" class="border-destructive/50 text-destructive bg-destructive/5"
							>{extension.latest?.name}</Badge
						>
					{/each}
				</div>
			</div>
		</Alert.Description>
		<Alert.Action onclick={alert.action}>{m.install()}</Alert.Action>
	</Alert.Root>
{/if}
