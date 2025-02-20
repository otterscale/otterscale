<script lang="ts">
	import Icon from '@iconify/svelte';
	import { mode, toggleMode } from 'mode-watcher';
	import { toast } from 'svelte-sonner';

	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	let favorited = false;
	function toggleFavorite() {
		if (!favorited) {
			favorited = true;
			toast.success('Added to favorites!');
			return;
		}
		favorited = false;
		toast.error('Removed from favorites!');
	}
</script>

<div class="flex justify-end space-x-2">
	<DropdownMenu.Root>
		<DropdownMenu.Trigger asChild let:builder>
			<Button builders={[builder]} class="w-14" variant="outline" size="icon">
				<Icon icon="material-symbols:add-2-rounded" class="h-5 w-5" />
				<Icon icon="material-symbols:arrow-drop-down-rounded" class="h-5 w-5" />
			</Button>
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="end">
			<DropdownMenu.Item class="space-x-2" on:click={toggleMode}>
				<span>{$mode === 'light' ? 'Use Dark Mode' : 'Use Light Mode'}</span>
			</DropdownMenu.Item>
			<DropdownMenu.Group>
				<DropdownMenu.Sub>
					<DropdownMenu.SubTrigger>Language</DropdownMenu.SubTrigger>
					<DropdownMenu.SubContent>
						<DropdownMenu.Item>Deutsch</DropdownMenu.Item>
						<DropdownMenu.Item>English</DropdownMenu.Item>
						<DropdownMenu.Item>Español</DropdownMenu.Item>
						<DropdownMenu.Item>Français</DropdownMenu.Item>
						<DropdownMenu.Item>Italiano</DropdownMenu.Item>
						<DropdownMenu.Item>日本語</DropdownMenu.Item>
						<DropdownMenu.Item>Português</DropdownMenu.Item>
						<DropdownMenu.Item>简体中文</DropdownMenu.Item>
						<DropdownMenu.Item>繁體中文</DropdownMenu.Item>
					</DropdownMenu.SubContent>
				</DropdownMenu.Sub>
			</DropdownMenu.Group>
			<DropdownMenu.Item
				class="space-x-2"
				on:click={() => window.open('https://openhdc.github.io', '_blank')}
			>
				<Icon icon="material-symbols:open-in-new" class="h-5 w-5" />
				<span>Documentation</span>
			</DropdownMenu.Item>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
	<Button variant="outline" size="icon" class="bg-header" on:click={toggleFavorite}>
		{#if favorited}
			<Icon icon="material-symbols:favorite" class="h-5 w-5" />
		{:else}
			<Icon icon="material-symbols:favorite-outline" class="h-5 w-5" />
		{/if}
	</Button>
	<Button variant="outline" size="icon" class="bg-header">
		<Icon icon="material-symbols:inbox-outline-rounded" class="h-5 w-5" />
	</Button>
</div>
