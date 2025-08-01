<script lang="ts" module>
	import * as Menubar from '$lib/components/ui/menubar/index';
	import { Menubar as MenubarPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import Separator from '$lib/components/ui/separator/separator.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: MenubarPrimitive.RootProps = $props();
</script>

<Menubar.Root
	bind:ref
	class={cn('border-none shadow-none [&_svg]:m-1', className)}
	{...restProps}
/>
<Separator />
