<script lang="ts" module>
	import {
		ConfigurationService,
		type Configuration
	} from '$lib/api/configuration/v1/configuration_pb';
	import { TagService, type Tag } from '$lib/api/tag/v1/tag_pb';
	import * as Table from '$lib/components/custom/table';
	import * as Accordion from '$lib/components/ui/accordion';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import UpdateNTPServer from './update-ntp-server.svelte';
	import UpdatePackageRepository from './update-package-repository.svelte';
	import CreateTag from './create-tag.svelte';
	import DeleteTag from './delete-tag.svelte';
	import * as Alert from '$lib/components/ui/alert';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const tagClient = createClient(TagService, transport);
	const configurationClient = createClient(ConfigurationService, transport);

	let configuration = $state(writable<Configuration>());
	let isConfigurationLoading = $state(true);
	async function fetchConfiguration() {
		try {
			configurationClient.getConfiguration({}).then((response) => {
				configuration.set(response);
				isConfigurationLoading = false;
			});
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let tags = $state(writable<Tag[]>());
	let isTagLoading = $state(true);

	async function fetchTags() {
		try {
			tagClient.listTags({}).then((response) => {
				tags.set(response.tags);
				isTagLoading = false;
			});
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchConfiguration();
			await fetchTags();
			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<main class="mx-auto w-7xl space-y-6">
	<Alert.Root variant="destructive" class="bg-destructive/5 border-destructive/10">
		<Icon icon="ph:warning" />
		<Alert.Description>
			These settings will apply across all scopes in your infrastructure. Please review changes
			carefully before applying them.</Alert.Description
		>
	</Alert.Root>
	{#if isMounted}
		<Accordion.Root type="single" value="ntp_servers" class="space-y-4">
			<Accordion.Item value="ntp_servers">
				<Accordion.Trigger
					class="bg-muted rounded-none p-6 hover:no-underline data-[state=closed]:rounded-lg data-[state=open]:rounded-t-lg"
				>
					<div class="flex flex-1 items-center gap-4">
						<div class="relative hidden h-16 w-16 overflow-hidden rounded-md border sm:block">
							<Icon icon="ph:network" class="h-full w-full object-cover p-2" />
						</div>
						<div class="flex-1 space-y-1 text-left">
							<h3 class="text-xl font-semibold">NTP Servers</h3>
							<Badge variant="outline" class="text-xs">network</Badge>
						</div>
					</div>
				</Accordion.Trigger>

				<Accordion.Content class="space-y-2 rounded-b-lg border border-t-0 p-6">
					{#if !isConfigurationLoading}
						<span class="flex justify-between gap-2">
							<h1 class="text-lg font-bold">Address</h1>
							<UpdateNTPServer bind:configuration />
						</span>
						<p class="text-muted-foreground text-sm">
							NTP servers, specified as IP addresses or hostnames delimited by commas and/or spaces,
							to be used as time references for MAAS itself, the machines MAAS deploys, and devices
							that make use of MAAS's DHCP services.
						</p>
						<div class="bg-muted rounded-lg p-4">
							{#if $configuration.ntpServer && $configuration.ntpServer.addresses}
								{#each $configuration.ntpServer.addresses as address}
									<Badge>{address}</Badge>
								{/each}
							{/if}
						</div>
					{/if}
				</Accordion.Content>
			</Accordion.Item>

			<Accordion.Item value="package_repository" class="group">
				<Accordion.Trigger
					class="bg-muted rounded-none p-6 hover:no-underline data-[state=closed]:rounded-lg data-[state=open]:rounded-t-lg"
				>
					<div class="flex flex-1 items-center gap-4">
						<div class="relative hidden h-16 w-16 overflow-hidden rounded-md border sm:block">
							<Icon icon="ph:cloud" class="h-full w-full object-cover p-2" />
						</div>
						<div class="flex-1 space-y-1 text-left">
							<h3 class="text-xl font-semibold">Package Repository</h3>
							<Badge variant="outline" class="text-xs">network</Badge>
						</div>
					</div>
				</Accordion.Trigger>

				<Accordion.Content class="space-y-2 rounded-b-lg border border-t-0 p-6">
					{#if !isConfigurationLoading}
						<h1 class="text-lg font-bold">Repository</h1>
						<p class="text-muted-foreground text-sm">
							Package repositories contain software packages that can be installed on machines.
							These repositories can be configured with custom URLs and enabled/disabled as needed
							to control software sources and updates for your infrastructure.
						</p>
						<div>
							<Table.Root>
								<Table.Header>
									<Table.Row>
										<Table.Head>NAME</Table.Head>
										<Table.Head>URL</Table.Head>
										<Table.Head>ENABLED</Table.Head>
										<Table.Head></Table.Head>
									</Table.Row>
								</Table.Header>
								<Table.Body>
									{#each $configuration.packageRepositories as packageRepository}
										<Table.Row class="border-none *:text-xs">
											<Table.Cell>{packageRepository.name}</Table.Cell>
											<Table.Cell>
												<span class="flex items-center gap-1">
													{packageRepository.url}
													<Button
														variant="ghost"
														target="_blank"
														href={packageRepository.url}
														class="flex items-start gap-1"
													>
														<Icon icon="ph:arrow-square-out" />
													</Button>
												</span>
											</Table.Cell>
											<Table.Cell>
												{packageRepository.enabled}
											</Table.Cell>
											<Table.Cell>
												<UpdatePackageRepository bind:configuration {packageRepository} />
											</Table.Cell>
										</Table.Row>
									{/each}
								</Table.Body>
							</Table.Root>
						</div>
					{/if}
				</Accordion.Content>
			</Accordion.Item>

			<Accordion.Item value="tag" class="group">
				<Accordion.Trigger
					class="bg-muted rounded-none p-6 hover:no-underline data-[state=closed]:rounded-lg data-[state=open]:rounded-t-lg"
				>
					<div class="flex flex-1 items-center gap-4">
						<div class="relative hidden h-16 w-16 overflow-hidden rounded-md border sm:block">
							<Icon icon="ph:tag" class="h-full w-full object-cover p-2" />
						</div>
						<div class="flex-1 space-y-1 text-left">
							<h3 class="text-xl font-semibold">Tag</h3>
							<Badge variant="outline" class="text-xs">machine</Badge>
						</div>
					</div>
				</Accordion.Trigger>
				<Accordion.Content class="space-y-2 rounded-b-lg border border-t-0 p-6">
					{#if !isConfigurationLoading}
						<span class="flex justify-between gap-2">
							<h1 class="text-lg font-bold">Tag</h1>
							<CreateTag bind:tags />
						</span>
						<p class="text-muted-foreground text-sm">
							Tags are identifiable labels that can be assigned to machines for filtering and group
							management. These tags help in organizing and managing machines based on their
							characteristics or purposes.
						</p>
						<div>
							<Table.Root>
								<Table.Header>
									<Table.Row>
										<Table.Head>TAG</Table.Head>
										<Table.Head>COMMENT</Table.Head>
										<Table.Head></Table.Head>
									</Table.Row>
								</Table.Header>
								<Table.Body>
									{#each $tags as tag}
										<Table.Row>
											<Table.Cell>{tag.name}</Table.Cell>
											<Table.Cell>
												<p class={cn(tag.comment ? 'text-primary' : 'text-muted-foreground')}>
													{tag.comment ? tag.comment : 'No comments available.'}
												</p>
											</Table.Cell>
											<Table.Cell>
												<DeleteTag {tag} bind:tags />
											</Table.Cell>
										</Table.Row>
									{/each}
								</Table.Body>
							</Table.Root>
						</div>
					{/if}
				</Accordion.Content>
			</Accordion.Item>
		</Accordion.Root>
	{/if}
</main>
