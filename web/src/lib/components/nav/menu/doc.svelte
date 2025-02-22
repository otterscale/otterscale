<script lang="ts">
	import Icon from '@iconify/svelte';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
</script>

<DropdownMenu.Item
	class="space-x-2"
	on:click={() => window.open('https://openhdc.github.io', '_blank')}
>
	<Icon icon="ph:arrow-square-out" class="h-4 w-4" />
	<span>Documentation</span>
</DropdownMenu.Item>
