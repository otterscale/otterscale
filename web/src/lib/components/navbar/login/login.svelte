<script lang="ts">
	import Icon from '@iconify/svelte';

	import { goto } from '$app/navigation';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { i18n } from '$lib/i18n';
	import { cn } from '$lib/utils';
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger
			class={cn(buttonVariants({ variant: 'outline', size: 'icon' }), 'bg-sidebar [&_svg]:size-5')}
			onclick={() => goto(i18n.resolveRoute('/login'))}
		>
			<Icon icon="ph:sign-in" />
		</Tooltip.Trigger>
		<Tooltip.Content>
			<p>Sign in</p>
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
