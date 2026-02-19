<script lang="ts">
	import { getComponent, getFormContext, type ComponentProps } from '@sjsf/form';

	const {
		children,
		keyInput,
		removeButton,
		config,
		errors
	}: ComponentProps['objectPropertyTemplate'] = $props();

	const ctx = getFormContext();

	const Layout = $derived(getComponent(ctx, 'layout', config));
</script>

<Layout type="object-property" {config} {errors}>
	<div class="flex w-full flex-col gap-3">
		{#if keyInput}
			<Layout type="object-property-key-input" {config} {errors}>
				{@render keyInput()}
			</Layout>
		{/if}
		<div class="*:data-[layout=object-property-content]:space-y-3">
			<Layout type="object-property-content" {config} {errors}>
				{@render children()}
			</Layout>
		</div>
		{#if removeButton}
			<div class="ml-auto">
				<Layout type="object-property-controls" {config} {errors}>
					<div
						class="[&_button]:size-7 [&_button]:border-none [&_button]:bg-transparent [&_button]:shadow-none"
					>
						{@render removeButton()}
					</div>
				</Layout>
			</div>
		{/if}
	</div>
</Layout>
