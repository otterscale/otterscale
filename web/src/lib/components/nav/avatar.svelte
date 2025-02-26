<script lang="ts">
	import Icon from '@iconify/svelte';

	import { goto } from '$app/navigation';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as m from '$lib/paraglide/messages.js';
	import pb from '$lib/pb';
	import { i18n } from '$lib/i18n';
	import { cn } from '$lib/utils';

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
		}
		return 'NA';
	}
</script>

{#if isValid}
	<DropdownMenu.Root>
		<DropdownMenu.Trigger
			class={cn(buttonVariants({ variant: 'ghost', size: 'icon' }), 'h-8 w-8 rounded-full')}
		>
			<Avatar.Root class="h-8 w-8">
				<Avatar.Image src={record ? pb.files.getURL(record, record.avatar) : ''} />
				<Avatar.Fallback>{getFallback()}</Avatar.Fallback>
				<span class="sr-only">Toggle user menu</span>
			</Avatar.Root>
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="end" class="w-48 [&_svg]:size-4">
			<DropdownMenu.Group>
				<DropdownMenu.GroupHeading>
					<div class="flex flex-col font-medium leading-none">
						<span class="text-sm">{record?.name}</span>
						<span class="text-xs text-muted-foreground">{record?.email}</span>
					</div>
				</DropdownMenu.GroupHeading>
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
				<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/recents'))}>
					<Icon icon="ph:clock" />
					<span class="pl-2">{m.avatar_recents()}</span>
					<DropdownMenu.Shortcut>⇧⌘R</DropdownMenu.Shortcut>
				</DropdownMenu.Item>
				<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/favorites'))}>
					<Icon icon="ph:clover" />
					<span class="pl-2">{m.avatar_favorites()}</span>
					<DropdownMenu.Shortcut>⇧⌘F</DropdownMenu.Shortcut>
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
