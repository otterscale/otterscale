<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { writable, type Writable } from 'svelte/store';

	import { OrchestratorService, type Plugin } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { Single as Alert } from '$lib/components/custom/alert';
	import { installPlugins } from '$lib/components/settings/plugins/utils.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const instancePlugins: Writable<Plugin[]> = writable([]);

	orchestratorClient
		.listInstancePlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			instancePlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch plugins:', error);
		});

	const alert: Alert.AlertType = $derived({
		title: m.plugins_alert_title(),
		message: m.plugins_alert_description(),
		action: () => {
			installPlugins(['instance']);
		},
		variant: 'destructive',
	});
</script>

{#if $instancePlugins.filter((instancePlugin) => instancePlugin.current).length < $instancePlugins.length}
	<Alert.Root {alert}>
		<Alert.Icon />
		<Alert.Title>{alert.title}</Alert.Title>
		<Alert.Description>
			<p>{alert.message}</p>
			<p>
				{$instancePlugins
					.filter((instancePlugin) => !instancePlugin.current)
					.map((instancePlugin) => instancePlugin.name)
					.join(', ')}
			</p>
		</Alert.Description>
		<Alert.Action onclick={alert.action}>{m.install()}</Alert.Action>
	</Alert.Root>
{/if}
