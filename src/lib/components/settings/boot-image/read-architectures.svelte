<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import { type Configuration_BootImage } from '$lib/api/configuration/v1/configuration_pb';
	import { Badge } from '$lib/components/ui/badge';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		bootImage
	}: {
		bootImage: Configuration_BootImage;
	} = $props();
</script>

{#if Object.keys(bootImage.architectureStatusMap).length > 0}
	<HoverCard.Root>
		<HoverCard.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
			<Icon icon="ph:info" />
		</HoverCard.Trigger>
		<HoverCard.Content class="max-h-[50vh] w-fit overflow-y-auto">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>{m.architecture()}</Table.Head>
						<Table.Cell>{m.status()}</Table.Cell>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each Object.entries(bootImage.architectureStatusMap) as [architecture, status] (architecture)}
						<Table.Row>
							<Table.Cell>{architecture}</Table.Cell>
							<Table.Cell>
								<Badge variant="outline">{status}</Badge>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</HoverCard.Content>
	</HoverCard.Root>
{/if}
