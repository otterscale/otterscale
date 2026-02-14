<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet';

	import CreateScope from './create-scope-form.svelte';
	import type { Plan } from './plans';

	let { open = $bindable(false), plan = $bindable({} as Plan) }: { open: boolean; plan: Plan } =
		$props();

	function handleClose() {
		open = false;
	}
</script>

<Sheet.Root bind:open onOpenChange={handleClose}>
	<Sheet.Content class="inset-y-auto bottom-0 h-9/10 rounded-tl-lg sm:max-w-4/5">
		<Sheet.Header class="h-full p-0">
			<div class="flex h-full">
				<div class="flex-1">
					<CreateScope {plan} oncancel={handleClose} onsuccess={handleClose} />
				</div>

				<!-- Plan Image -->
				<div class="relative hidden w-2/5 lg:block">
					<img src={plan.image} alt={plan.name} class="absolute inset-0 size-full object-cover" />
				</div>
			</div>
		</Sheet.Header>
	</Sheet.Content>
</Sheet.Root>
