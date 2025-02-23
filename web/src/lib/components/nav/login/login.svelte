<script lang="ts">
	import Icon from '@iconify/svelte';

	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		<Button variant="outline" size="icon" class="bg-header" on:click={() => goto('/login')}>
			<Icon icon="ph:sign-in" class="h-5 w-5" />
		</Button>
	</Tooltip.Trigger>
	<Tooltip.Content>
		<p>Sign in</p>
	</Tooltip.Content>
</Tooltip.Root>
