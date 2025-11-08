<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import { type Configuration_BootImage } from '$lib/api/configuration/v1/configuration_pb';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		bootImage
	}: {
		bootImage: Configuration_BootImage;
	} = $props();
</script>

<HoverCard.Root>
	<HoverCard.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
		<Icon icon="ph:info" />
	</HoverCard.Trigger>
	<HoverCard.Content class="max-h-[50vh] w-fit overflow-y-auto">
		{#if Object.keys(bootImage.architectureStatusMap).length > 0}
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>{m.architecture()}</Table.Head>
						<Table.Cell>{m.status()}</Table.Cell>
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
			<p class="w-full p-2 text-center text-xs font-light text-muted-foreground">
				{m.no_data()}
			</p>
		{/if}
	</HoverCard.Content>
</HoverCard.Root>
