<script lang="ts">
    let {
        value,
        unit,
    }: {
        value?: string;
        unit?: string;
    } = $props();
</script>

<div class="aspect-auto h-[250px] w-full flex items-center justify-center">
    <p class="fill-foreground font-bold text-4xl!">{value || 'No Data'} {unit || ''}</p>
</div>