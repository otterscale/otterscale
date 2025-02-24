<script lang="ts">
	import Icon from '@iconify/svelte';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as m from '$lib/paraglide/messages.js';
</script>

<DropdownMenu.Item
	class="space-x-2"
	on:click={() => window.open('https://openhdc.github.io', '_blank')}
>
	<Icon icon="ph:arrow-square-out" class="h-4 w-4" />
	<span>{m.nav_preferences_documentation()}</span>
</DropdownMenu.Item>
