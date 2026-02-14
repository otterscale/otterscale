<script lang="ts" module>
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
</script>

<script lang="ts">
	let { ref = $bindable(null), ...restProps }: DropdownMenuPrimitive.GroupProps = $props();
</script>

<DropdownMenu.Group bind:ref data-slot="alert-dialog-action-group" {...restProps} />
