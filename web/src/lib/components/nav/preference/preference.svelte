<script>
	import Icon from '@iconify/svelte';

	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Tooltip from '$lib/components/ui/tooltip';

	import Mode from './mode.svelte';
	import I18n from './i18n.svelte';
	import Doc from './doc.svelte';
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger asChild let:builder>
		<Tooltip.Root>
			<Tooltip.Trigger asChild>
				<Button builders={[builder]} variant="outline" size="icon" class="bg-header">
					<Icon icon="ph:circles-three-plus" class="h-5 w-5" />
				</Button>
			</Tooltip.Trigger>
			<Tooltip.Content>
				<p>Preference</p>
			</Tooltip.Content>
		</Tooltip.Root>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end" class="w-40">
		<DropdownMenu.Label>Preference</DropdownMenu.Label>
		<DropdownMenu.Separator />
		<Mode />
		<I18n />
		<DropdownMenu.Separator />
		<Doc />
	</DropdownMenu.Content>
</DropdownMenu.Root>
