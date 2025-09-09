<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import Update from './update.svelte';

	import { ConfigurationService, type Configuration } from '$lib/api/configuration/v1/configuration_pb';
	import * as Layout from '$lib/components/settings/layout';
	import * as Card from '$lib/components/ui/card';
	import * as Tooltip from '$lib/components/ui/tooltip';
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
				<div class="grid grid-cols-4 gap-4">
					{#each $configuration.ntpServer.addresses as address}
						<Card.Root class="group">
							<Card.Content class="flex size-fit items-center gap-2">
								<div class="bg-muted group-hover:bg-muted-foreground size-fit rounded-full p-2">
									<Icon icon="ph:cloud" class="text-muted-foreground group-hover:text-muted size-6" />
								</div>
								<div>
									<p class="text-base font-medium">Address</p>
									<Tooltip.Root>
										<Tooltip.Trigger>
											<p class="text-muted-foreground max-w-3xs truncate text-sm">
												{address}
											</p>
										</Tooltip.Trigger>
										<Tooltip.Content>
											{address}
										</Tooltip.Content>
									</Tooltip.Root>
								</div>
							</Card.Content>
						</Card.Root>
					{/each}
				</div>
			{/if}
		</Layout.Viewer>
	</Layout.Root>
{/if}
