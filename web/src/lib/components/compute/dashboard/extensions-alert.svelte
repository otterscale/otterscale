<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';

	import { type Extension, OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { Single as Alert } from '$lib/components/custom/alert';
	import { installExtensions } from '$lib/components/settings/extensions/utils.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const generalExtensions: Writable<Extension[]> = writable([]);

	orchestratorClient
		.listGeneralExtensions({ scope: scope })
		.then((response) => {
			generalExtensions.set(response.Extensions);
		})
		.catch((error) => {
			console.error('Failed to fetch extensions:', error);
		});

	const alert: Alert.AlertType = $derived({
		title: m.extensions_alert_title(),
		message: m.extensions_alert_description(),
		action: () => {
			installExtensions(scope, ['general']);
		},
		variant: 'destructive'
	});
</script>

{#if $generalExtensions.filter((generalExtension) => generalExtension.current).length < $generalExtensions.length}
	<Alert.Root {alert}>
		<Alert.Icon />
		<Alert.Title>{alert.title}</Alert.Title>
		<Alert.Description>
			<div class="space-y-1">
				<p>{alert.message}</p>
				<div class="flex w-full flex-wrap gap-2">
					{#each $generalExtensions.filter((extension) => !extension.current) as extension}
						<Badge variant="outline" class="border-destructive/50 bg-destructive/5 text-destructive"
							>{extension.name}</Badge
						>
					{/each}
				</div>
			</div>
		</Alert.Description>
		<Alert.Action onclick={alert.action}>{m.install()}</Alert.Action>
	</Alert.Root>
{/if}
