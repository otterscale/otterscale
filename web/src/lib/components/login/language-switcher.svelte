<script lang="ts">
	import Icon from '@iconify/svelte';

	import type { LanguageSwitcherProps } from './types';

	import { buttonVariants } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils.js';

	let {
		languages = [],
		value = $bindable(''),
		align = 'end',
		variant = 'ghost',
		onChange,
		class: className
	}: LanguageSwitcherProps = $props();

	// set default code if there isn't one selected
	if (value === '') {
		value = languages[0].code;
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={cn(buttonVariants({ variant, size: 'icon' }), className)}
		aria-label="Change language"
	>
		<Icon icon="ph:globe" class="size-5" />
		<span class="sr-only">Change language</span>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content {align}>
		<DropdownMenu.RadioGroup bind:value onValueChange={onChange}>
			{#each languages as language (language.code)}
				<DropdownMenu.RadioItem value={language.code}>
					{language.label}
				</DropdownMenu.RadioItem>
			{/each}
		</DropdownMenu.RadioGroup>
	</DropdownMenu.Content>
</DropdownMenu.Root>
