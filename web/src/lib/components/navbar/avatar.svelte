<script lang="ts">
	import type { User } from 'better-auth';
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as m from '$lib/paraglide/messages.js';
	import { i18n } from '$lib/i18n';
	import { cn } from '$lib/utils';

	export let user: User;
</script>

{#if user}
	<DropdownMenu.Root>
		<DropdownMenu.Trigger
			class={cn(buttonVariants({ variant: 'ghost', size: 'icon' }), 'h-8 w-8 rounded-full')}
		>
			<Avatar.Root class="h-8 w-8">
				<Avatar.Image src={user.image} />
				<Avatar.Fallback>{user.name[0]}</Avatar.Fallback>
				<span class="sr-only">Toggle user menu</span>
			</Avatar.Root>
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="end" class="w-48 [&_svg]:size-4">
			<DropdownMenu.Label class="p-0 font-normal">
				<div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
					<Avatar.Root class="h-8 w-8">
						<Avatar.Image src={user.image} />
						<Avatar.Fallback>{user.name[0]}</Avatar.Fallback>
						<span class="sr-only">Toggle user menu</span>
					</Avatar.Root>
					<div class="grid flex-1 text-left text-sm leading-tight">
						<span class="truncate font-semibold"> {user.name}</span>
						<span class="truncate text-xs"> {user.email}</span>
					</div>
				</div>
			</DropdownMenu.Label>
			<DropdownMenu.Group>
				<DropdownMenu.Separator />
				<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/settings#profile'))}>
					<Icon icon="ph:user" />
					<span class="pl-2">{m.avatar_profile()}</span>
					<DropdownMenu.Shortcut>⇧⌘P</DropdownMenu.Shortcut>
				</DropdownMenu.Item>
				<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/settings'))}>
					<Icon icon="ph:gear" />
					<span class="pl-2">{m.avatar_settings()}</span>
					<DropdownMenu.Shortcut>⇧⌘S</DropdownMenu.Shortcut>
				</DropdownMenu.Item>
				<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/settings#superuser'))}>
					<Icon icon="ph:key" />
					<span class="pl-2">{m.avatar_superuser()}</span>
					<DropdownMenu.Shortcut>⇧⌘U</DropdownMenu.Shortcut>
				</DropdownMenu.Item>
				<DropdownMenu.Separator />
				<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/logout'))}>
					<Icon icon="ph:sign-in" />
					<span class="pl-2">{m.avatar_sign_out()}</span>
					<DropdownMenu.Shortcut>⇧⌘Q</DropdownMenu.Shortcut>
				</DropdownMenu.Item>
			</DropdownMenu.Group>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
{/if}
