<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import Update from './update.svelte';

	import { ConfigurationService, type Configuration } from '$lib/api/configuration/v1/configuration_pb';
	import * as Layout from '$lib/components/settings/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const configurationClient = createClient(ConfigurationService, transport);

	const configuration = writable<Configuration>();
	let isConfigurationLoading = $state(true);

	onMount(async () => {
		try {
			await configurationClient.getConfiguration({}).then((response) => {
				configuration.set(response);
				isConfigurationLoading = false;
			});
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
			{#if $configuration.ntpServer && $configuration.ntpServer.addresses}
				<div class="grid gap-4 space-y-8 sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-5 w-full">
					{#each $configuration.ntpServer.addresses as address}
						<div class="flex flex-col items-center gap-4">
							<div
								class="bg-muted/50 m-2 rounded-full p-2 shadow-lg"
							>
								<Icon icon="ph:clock-user" class="m-2 size-20" />
							</div>
							<div class="flex flex-col items-center">
								<p class="text-sm font-bold">{address}</p>
								<p class="text-muted-foreground text-xs font-semibold">ntp server</p>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</Layout.Viewer>
	</Layout.Root>
{/if}
