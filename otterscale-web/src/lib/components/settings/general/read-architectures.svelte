<script lang="ts" module>
	import { type Configuration_BootImage } from '$lib/api/configuration/v1/configuration_pb';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import Icon from '@iconify/svelte';
</script>

<script lang="ts">
	let {
		bootImage
	}: {
		bootImage: Configuration_BootImage;
	} = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger>
		<Icon icon="ph:info" data-tooltip-target="architectures-{bootImage.name}" />
	</HoverCard.Trigger>
	<HoverCard.Content class="w-fit">
		{#if Object.keys(bootImage.architectureStatusMap).length > 0}
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>ARCHITECTURE</Table.Head>
						<Table.Cell>STATUS</Table.Cell>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each Object.entries(bootImage.architectureStatusMap) as [architecture, status]}
						<Table.Row>
							<Table.Cell>{architecture}</Table.Cell>
							<Table.Cell>
								<Badge variant="outline">{status}</Badge>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		{:else}
			<p class="text-muted-foreground w-full p-2 text-center text-xs font-light">
				No architectures available at the moment.
			</p>
		{/if}
	</HoverCard.Content>
</HoverCard.Root>
