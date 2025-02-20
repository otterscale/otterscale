<script>
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger asChild let:builder>
		<Button builders={[builder]} variant="secondary" size="icon" class="h-8 w-8">
			<Avatar.Root class="h-8 w-8">
				<Avatar.Image src="https://github.com/iasthc.png" alt="@iasthc" />
				<Avatar.Fallback>PT</Avatar.Fallback>
				<span class="sr-only">Toggle user menu</span>
			</Avatar.Root>
		</Button>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end">
		<DropdownMenu.Label>
			<div class="flex flex-col font-medium leading-none">
				<span class="text-sm">Paul Tsai</span>
				<span class="text-xs text-muted-foreground">paul_tsai@phison.com</span>
			</div>
			<DropdownMenu.Separator />
		</DropdownMenu.Label>
		<DropdownMenu.Item on:click={() => goto('/settings/profile')}>
			<span>Profile</span>
			<DropdownMenu.Shortcut>⇧⌘P</DropdownMenu.Shortcut>
		</DropdownMenu.Item>
		<DropdownMenu.Separator />
		<DropdownMenu.Item on:click={() => goto('/account/logout')}>
			<span>Sign out</span>
			<DropdownMenu.Shortcut>⇧⌘Q</DropdownMenu.Shortcut>
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
