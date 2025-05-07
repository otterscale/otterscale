<script lang="ts">
	import { Button } from '$lib/components/ui/button';
</script>

<section id="cta" class="my-24 bg-muted/50 py-16 sm:my-32">
	<div class="container place-items-center lg:grid lg:grid-cols-2">
		<div class="lg:col-start-1">
			<h2 class="text-3xl font-bold md:text-4xl">
				All Your
				<span class="bg-gradient-to-b from-primary/60 to-primary bg-clip-text text-transparent">
					{' '}
					Ideas & Concepts{' '}
				</span>
				In One Interface
			</h2>
			<p class="mb-8 mt-4 text-xl text-muted-foreground lg:mb-0">
				Lorem ipsum dolor sit amet consectetur adipisicing elit. Eaque, beatae. Ipsa tempore ipsum
				iste quibusdam illum ducimus eos. Quasi, sed!
			</p>
		</div>

		<div class="space-y-4 lg:col-start-2">
			<Button class="w-full md:mr-4 md:w-auto">Request a Demo</Button>
			<Button variant="outline" class="w-full md:w-auto">View all features</Button>
		</div>
	</div>
</section>
