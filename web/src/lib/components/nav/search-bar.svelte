<script>
	import Icon from '@iconify/svelte';

	import { Input } from '$lib/components/ui/input';
	import * as m from '$lib/paraglide/messages.js';
</script>

<form class="ml-auto flex-1 sm:flex-initial">
	<div class="relative">
		<Icon icon="ph:magnifying-glass" class="absolute left-2 top-2 h-5 w-5 text-muted-foreground" />
		<Input
			type="search"
			placeholder={m.nav_search()}
			class="pl-8 sm:w-[300px] md:w-[200px] lg:w-[300px]"
		/>
	</div>
</form>
