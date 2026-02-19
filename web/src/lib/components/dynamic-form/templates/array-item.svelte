<script lang="ts">
	import { getComponent, getFormContext, type ComponentProps } from '@sjsf/form';

	const { children, buttons, config, errors }: ComponentProps['arrayItemTemplate'] = $props();

	const ctx = getFormContext();

	const Layout = $derived(getComponent(ctx, 'layout', config));
</script>

<div class="w-full border-l pl-4">
	<Layout type="array-item" {config} {errors}>
		<div class="flex w-full flex-col gap-3">
			<Layout type="array-item-content" {config} {errors}>
				{@render children()}
			</Layout>
			{#if buttons}
				<div class="ml-auto">
					<Layout type="array-item-controls" {config} {errors}>
						<div
							class="ml-auto [&_button]:size-7 [&_button]:border-none [&_button]:bg-transparent [&_button]:shadow-none"
						>
							{@render buttons()}
						</div>
					</Layout>
				</div>
			{/if}
		</div>
	</Layout>
</div>
