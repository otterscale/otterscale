<script lang="ts">
	import Icon from '@iconify/svelte';

	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import pb from '$lib/pb';

	let isValid = pb.authStore.isValid;
	let record = pb.authStore.record;

	pb.authStore.onChange((_, r) => {
		isValid = pb.authStore.isValid;
		record = r;
	});

	function getFallback(): string {
		if (record) {
			const names = record.name.split(' ');
			if (names.length >= 2) {
				return (names[0][0] + names[1][0]).toUpperCase();
			}
			return names[0][0].toUpperCase();
		}
		return '';
	}
</script>

{#if isValid}
	<DropdownMenu.Root>
		<DropdownMenu.Trigger asChild let:builder>
			<Button builders={[builder]} variant="ghost" size="icon" class="h-8 w-8 rounded-full">
				<Avatar.Root class="h-8 w-8">
					<Avatar.Image src={record ? pb.files.getURL(record, record.avatar) : ''} />
					<Avatar.Fallback>{getFallback()}</Avatar.Fallback>
					<span class="sr-only">Toggle user menu</span>
				</Avatar.Root>
			</Button>
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="end" class="w-48">
			<DropdownMenu.Label>
				<div class="flex flex-col font-medium leading-none">
					<span class="text-sm">{record?.name}</span>
					<span class="text-xs text-muted-foreground">{record?.email}</span>
				</div>
				<DropdownMenu.Separator />
			</DropdownMenu.Label>
			<DropdownMenu.Item on:click={() => goto('/settings/profile')}>
				<Icon icon="ph:user" class="h-4 w-4" />
				<span class="pl-2">Profile</span>
				<DropdownMenu.Shortcut>⇧⌘P</DropdownMenu.Shortcut>
			</DropdownMenu.Item>
			<DropdownMenu.Item on:click={() => goto('/settings/profile')}>
				<Icon icon="ph:gear" class="h-4 w-4" />
				<span class="pl-2">Settings</span>
				<DropdownMenu.Shortcut>⇧⌘S</DropdownMenu.Shortcut>
			</DropdownMenu.Item>
			<DropdownMenu.Item on:click={() => goto('/settings/profile')}>
				<Icon icon="ph:key" class="h-4 w-4" />
				<span class="pl-2">Superuser</span>
				<DropdownMenu.Shortcut>⇧⌘U</DropdownMenu.Shortcut>
			</DropdownMenu.Item>
			<DropdownMenu.Separator />
			<DropdownMenu.Item on:click={() => goto('/logout')}>
				<Icon icon="ph:sign-in" class="h-4 w-4" />
				<span class="pl-2">Sign out</span>
				<DropdownMenu.Shortcut>⇧⌘Q</DropdownMenu.Shortcut>
			</DropdownMenu.Item>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
{/if}
