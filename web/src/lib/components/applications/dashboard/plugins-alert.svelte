<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { writable, type Writable } from 'svelte/store';

	import { OrchestratorService, type Plugin } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { Single as Alert } from '$lib/components/custom/alert';
	import { installPlugins } from '$lib/components/settings/plugins/utils.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const generalPlugins: Writable<Plugin[]> = writable([]);

	orchestratorClient
		.listGeneralPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			generalPlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch plugins:', error);
		});

	const alert: Alert.AlertType = $derived({
		title: m.plugins_alert_title(),
		message: m.plugins_alert_description(),
		action: () => {
			installPlugins(['general']);
		},
		variant: 'destructive',
	});
</script>

{#if $generalPlugins.filter((generalPlugin) => generalPlugin.current).length < $generalPlugins.length}
	<Alert.Root {alert}>
		<Alert.Icon />
		<Alert.Title>{alert.title}</Alert.Title>
		<Alert.Description>
			<div class="space-y-1">
				<p>{alert.message}</p>
				<div class="flex w-full flex-wrap gap-2">
					{#each $generalPlugins.filter((plugin) => !plugin.current) as plugin}
						<Badge variant="outline" class="border-destructive/50 text-destructive bg-destructive/5"
							>{plugin.latest?.name}</Badge
						>
					{/each}
				</div>
			</div>
		</Alert.Description>
		<Alert.Action onclick={alert.action}>{m.install()}</Alert.Action>
	</Alert.Root>
{/if}
