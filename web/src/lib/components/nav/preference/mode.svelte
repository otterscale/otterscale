<script lang="ts">
	import Icon from '@iconify/svelte';
	import { mode, toggleMode } from 'mode-watcher';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
</script>

<DropdownMenu.Item class="space-x-2" on:click={toggleMode}>
	<Icon icon={$mode === 'light' ? 'ph:sun' : 'ph:moon'} class="h-4 w-4" />
	<span>{$mode === 'light' ? 'Light Mode' : 'Dark Mode'}</span>
</DropdownMenu.Item>
