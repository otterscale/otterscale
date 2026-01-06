<script lang="ts">
	import Icon from '@iconify/svelte';
	import LogInIcon from '@lucide/svelte/icons/log-in';
	import type { HTMLAttributes } from 'svelte/elements';

	import { resolve } from '$app/paths';
	import LogoImage from '$lib/assets/logo.png';
	import { Button } from '$lib/components/ui/button';
	import { Field, FieldDescription, FieldGroup, FieldSeparator } from '$lib/components/ui/field';
	import { m } from '$lib/paraglide/messages';
	import { cn, type WithElementRef } from '$lib/utils';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<div class={cn('flex max-w-sm flex-col gap-6', className)} bind:this={ref} {...restProps}>
	<FieldGroup>
		<div class="flex flex-col items-center gap-2 text-center">
			<a
				target="_blank"
				href="https://otterscale.com"
				class="flex flex-col items-center gap-2 font-medium"
			>
				<div class="flex size-18 items-center justify-center rounded-md">
					<img src={LogoImage} alt="logo" class="size-18" />
				</div>
				<span class="sr-only">OtterScale</span>
			</a>
			<h1 class="text-2xl font-bold">{m.welcome_to({ name: 'OtterScale' })}</h1>
			<FieldDescription>
				{m.login_description()}
			</FieldDescription>
		</div>

		<Field>
			<Button href={resolve('/login/keycloak')}>
				<LogInIcon />
				{m.login()}
			</Button>
		</Field>

		<FieldSeparator class="[&_span]:bg-muted">
			<div class="flex items-center gap-2 text-xs text-muted-foreground/70">
				<Icon icon="simple-icons:keycloak" />
				{m.login_divider()}
			</div>
		</FieldSeparator>
	</FieldGroup>

	<FieldDescription class="px-6 text-center text-xs">
		<!-- eslint-disable-next-line svelte/no-at-html-tags -->
		{@html m.login_footer({
			terms_of_service: `<a href="${resolve('/terms-of-service')}">${m.terms_of_service()}</a>`,
			privacy_policy: `<a href="${resolve('/privacy-policy')}">${m.privacy_policy()}</a>`
		})}
	</FieldDescription>
</div>
