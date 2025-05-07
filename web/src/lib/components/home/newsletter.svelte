<script>
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
</script>

<section id="newsletter">
	<hr class="mx-auto w-11/12" />

	<div class="container py-24 sm:py-32">
		<h3 class="text-center text-4xl font-bold md:text-5xl">
			Join Our Daily{' '}
			<span class="bg-gradient-to-b from-primary/60 to-primary bg-clip-text text-transparent">
				Newsletter
			</span>
		</h3>
		<p class="mb-8 mt-4 text-center text-xl text-muted-foreground">
			Lorem ipsum dolor sit amet consectetur.
		</p>

		<form class="mx-auto flex w-full flex-col gap-4 md:w-6/12 md:flex-row md:gap-2 lg:w-4/12">
			<Input
				placeholder="leomirandadev@gmail.com"
				class="bg-muted/50 dark:bg-muted/80 "
				aria-label="email"
			/>
			<Button onclick={() => toast.success('Subscribed!')}>Subscribe</Button>
		</form>
	</div>

	<hr class="mx-auto w-11/12" />
</section>
