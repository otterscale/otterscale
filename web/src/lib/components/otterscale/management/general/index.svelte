<script lang="ts">
	import UpdateNTPServer from './update-ntp-server.svelte';
	import UpdatePackageRepository from './update-package-repository.svelte';
	import CreateBootImage from './boot-image/create.svelte';
	import ImportBootImage from './boot-image/import.svelte';
	import SetBootImageAsDefault from './boot-image/set-default.svelte';
	import IsImportingBootImages from './boot-image/is-importing.svelte';
	import CreateTag from './tag/create.svelte';
	import DeleteTag from './tag/delete.svelte';
	import Icon from '@iconify/svelte';
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';

	import {
		type Configuration,
		type Configuration_BootImage,
		type Tag
	} from '$gen/api/nexus/v1/nexus_pb';

	let {
		configuration,
		tags
	}: {
		configuration: Configuration;
		tags: Tag[];
	} = $props();

	let isImportingBootImages = $state(false);
</script>

<div class="grid gap-2 p-4">
	<fieldset class="grid items-center gap-2 rounded-lg border-none">
		<legend class="w-full rounded-md bg-muted px-2 py-1 text-lg font-semibold">NTP Servers</legend>
		<div class="flex justify-between">
			<div class="flex flex-col justify-between gap-2 p-4">
				{#if configuration.ntpServer}
					<span class="flex items-center gap-1">
						{#if configuration.ntpServer.addresses}
							{#each configuration.ntpServer.addresses as address}
								<Badge variant="outline" class="text-xs">{address}</Badge>
							{/each}
						{/if}
					</span>
				{/if}
				<p class="text-xs font-light text-muted-foreground">
					NTP servers, specified as IP addresses or hostnames delimited by commas and/or spaces, to
					be used as time references for MAAS itself, the machines MAAS deploys, and devices that
					make use of MAAS's DHCP services.
				</p>
			</div>
			<UpdateNTPServer {configuration} />
		</div>
	</fieldset>

	<fieldset class="grid items-center gap-2 rounded-lg border-none">
		<legend class="w-full rounded-md bg-muted px-2 py-1 text-lg font-semibold"
			>Package Repositories</legend
		>
		<div class="p-4">
			<Table.Root>
				<Table.Header>
					<Table.Row class="*:text-xs *:font-light">
						<Table.Head>NAME</Table.Head>
						<Table.Head>URL</Table.Head>
						<Table.Head class="text-center">ENABLED</Table.Head>
						<Table.Head></Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each configuration.packageRepositories as packageRepository}
						<Table.Row class="border-none *:text-xs">
							<Table.Cell>{packageRepository.name}</Table.Cell>
							<Table.Cell>
								<a href={packageRepository.url} target="_blank" rel="noopener noreferrer">
									<span class="flex items-start gap-1">
										{packageRepository.url}
										<Icon icon="ph:arrow-square-out" class="h-4 w-4" />
									</span>
								</a>
							</Table.Cell>
							<Table.Cell>
								<span class="flex justify-center">
									{@render formatterBoolean(packageRepository.enabled)}
								</span>
							</Table.Cell>
							<Table.Cell class="text-right">
								<UpdatePackageRepository {packageRepository} />
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</fieldset>

	<fieldset class="grid items-center gap-2 rounded-lg border-none">
		<legend class="w-full rounded-md bg-muted px-2 py-1 text-lg font-semibold">Boot Image</legend>
		<span class="ml-auto flex p-2">
			<CreateBootImage />

			{#key isImportingBootImages}
				<ImportBootImage bind:isImportingBootImages />
			{/key}
		</span>
		<div class="px-4">
			<Table.Root>
				<Table.Header>
					<Table.Row class="*:text-xs *:font-light">
						<Table.Head>NAME</Table.Head>
						<Table.Head>SOURCE</Table.Head>
						<Table.Head>DISTRO SERIES</Table.Head>
						<Table.Head class="text-center">DEFAULT</Table.Head>
						<Table.Head class="text-right">ARCHITECTURES</Table.Head>
						<Table.Head></Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each configuration.bootImages as bootImage}
						<Table.Row class="border-none *:text-xs">
							<Table.Cell>{bootImage.name}</Table.Cell>
							<Table.Cell>{bootImage.source}</Table.Cell>
							<Table.Cell>
								<Badge variant="outline">{bootImage.distroSeries}</Badge>
							</Table.Cell>
							<Table.Cell>
								<span class="flex justify-center">
									{@render formatterBoolean(bootImage.default)}
								</span>
							</Table.Cell>
							<Table.Cell>
								<span class="flex justify-end">
									<span class="flex items-center gap-1">
										{Object.keys(bootImage.architectureStatusMap).length}
										{@render ReadArchitectures(bootImage)}
									</span>
								</span>
							</Table.Cell>
							<Table.Cell>
								<span class="flex justify-end">
									<SetBootImageAsDefault {bootImage} />
								</span>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</fieldset>

	<fieldset class="grid items-center gap-2 rounded-lg border-none">
		<legend class="w-full rounded-md bg-muted px-2 py-1 text-lg font-semibold">Tags</legend>
		<div class="ml-auto p-2">
			<CreateTag />
		</div>
		<div class="px-4">
			<Table.Root>
				<Table.Header>
					<Table.Row class="*:text-xs *:font-light">
						<Table.Head>TAG</Table.Head>
						<Table.Head>COMMENT</Table.Head>
						<Table.Head></Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each tags as tag}
						<Table.Row class="border-none *:text-xs">
							<Table.Cell>{tag.name}</Table.Cell>
							<Table.Cell>{tag.comment || 'No comments available'}</Table.Cell>
							<Table.Cell>
								<span class="flex justify-end">
									<DeleteTag {tag} />
								</span>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</fieldset>
</div>

{#snippet ReadArchitectures(bootImage: Configuration_BootImage)}
	<HoverCard.Root openDelay={13}>
		<HoverCard.Trigger>
			<Icon
				icon="ph:info"
				class="h-4 w-4 text-blue-700"
				data-tooltip-target="architectures-{bootImage.name}"
			/>
		</HoverCard.Trigger>
		<HoverCard.Content class="w-fit">
			<Table.Root>
				<Table.Body>
					<Table.Row>
						<Table.Head class="text-xs font-light">ARCHITECTURE</Table.Head>
						<Table.Cell class="text-xs font-light">STATUS</Table.Cell>
					</Table.Row>
					{#each Object.entries(bootImage.architectureStatusMap) as [architecture, status]}
						<Table.Row class="border-none">
							<Table.Head class="whitespace-nowrap text-xs">{architecture}</Table.Head>
							<Table.Cell class="text-right text-xs">{status}</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</HoverCard.Content>
	</HoverCard.Root>
{/snippet}

{#snippet formatterBoolean(b: boolean)}
	<Icon icon={b ? 'ph:check' : 'ph:x'} />
{/snippet}
