<script>
	import Icon from '@iconify/svelte';

	import { buttonVariants } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as m from '$lib/paraglide/messages.js';

	import Mode from './mode.svelte';
	import I18n from './i18n.svelte';
	import Doc from './doc.svelte';
	import { cn } from '$lib/utils';
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={cn(buttonVariants({ variant: 'outline', size: 'icon' }), 'bg-sidebar [&_svg]:size-5')}
	>
		<Icon icon="ph:circles-three-plus" />
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end">
		<DropdownMenu.Group>
			<DropdownMenu.GroupHeading>{m.nav_preferences()}</DropdownMenu.GroupHeading>
			<DropdownMenu.Separator />
			<Mode />
			<I18n />
			<DropdownMenu.Separator />
			<Doc />
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
