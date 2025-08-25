<script lang="ts" module>
	import {
		ConfigurationService,
		type Configuration
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Table from '$lib/components/custom/table';
	import * as Layout from '$lib/components/settings/layout';
	import { m } from '$lib/paraglide/messages';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import UpdatePackageRepository from './update-package-repository.svelte';
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
	<Layout.Title>{m.repository()}</Layout.Title>
	<Layout.Description>{m.setting_package_repository_description()}</Layout.Description>
	<Layout.Controller>
		<div class="rounded-lg border shadow-sm">
			<Table.Root>
				<Table.Header>
					<Table.Row
						class="*:bg-muted *:rounded-t-lg *:px-4 *:first:rounded-tl-lg *:last:rounded-tr-lg"
					>
						<Table.Head>{m.name()}</Table.Head>
						<Table.Head>{m.url()}</Table.Head>
						<Table.Head>{m.enabled()}</Table.Head>
						<Table.Head></Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each $configuration.packageRepositories as packageRepository}
						<Table.Row class="*:px-4">
							<Table.Cell>{packageRepository.name}</Table.Cell>
							<Table.Cell>
								<a
									href={packageRepository.url}
									class="flex items-start gap-1 underline hover:no-underline"
								>
									{packageRepository.url}
								</a>
							</Table.Cell>
							<Table.Cell>
								<Icon icon={packageRepository.enabled ? 'ph:circle' : 'ph:x'} />
							</Table.Cell>
							<Table.Cell>
								<div class="flex items-center justify-end">
									<UpdatePackageRepository {configuration} {packageRepository} />
								</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</Layout.Controller>
{/if}
