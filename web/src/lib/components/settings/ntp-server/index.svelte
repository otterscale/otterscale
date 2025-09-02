<script lang="ts" module>
	import { ConfigurationService, type Configuration } from '$lib/api/configuration/v1/configuration_pb';
	import * as Layout from '$lib/components/settings/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import Update from './update.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const configurationClient = createClient(ConfigurationService, transport);

	const configuration = writable<Configuration>();
	let isConfigurationLoading = $state(true);

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await configurationClient.getConfiguration({}).then((response) => {
				configuration.set(response);
				isConfigurationLoading = false;
			});
			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !isConfigurationLoading}
	<Layout.Root>
		<Layout.Title>{m.address()}</Layout.Title>
		<Layout.Description>
			{m.setting_ntp_server_description()}
		</Layout.Description>
		<Layout.Controller>
			<Update {configuration} />
		</Layout.Controller>
		<Layout.Viewer>
			<Card.Root>
				<Card.Content>
					{#if $configuration.ntpServer && $configuration.ntpServer.addresses}
						{#each $configuration.ntpServer.addresses as address}
							<Badge>{address}</Badge>
						{/each}
					{/if}
				</Card.Content>
			</Card.Root>
		</Layout.Viewer>
	</Layout.Root>
{/if}
