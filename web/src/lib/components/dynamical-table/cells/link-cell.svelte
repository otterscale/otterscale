<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	let { display, hyperlink }: { display: string; hyperlink: string } = $props();
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger>
			<Button href={hyperlink} variant="link" size="sm">
				{display}
			</Button>
		</Tooltip.Trigger>
		<Tooltip.Content>
			{hyperlink}
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
