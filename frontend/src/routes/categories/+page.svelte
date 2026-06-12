<script lang="ts">
	import { getCategories } from '$lib/remote/categories.remote';
	import Seo from '$lib/components/Seo.svelte';

	const categories = getCategories();
</script>

<Seo
	title="Categorias de anúncios | Despacha Aí"
	description="Explora anúncios por categoria em Benguela e toda Angola: electrónica, veículos, imóveis, moda, móveis e muito mais."
/>

<h1 class="text-2xl font-bold text-neutral-900">Categorias</h1>
<p class="mt-1 text-sm text-neutral-500">Explora os anúncios por categoria.</p>

<svelte:boundary>
	{@const cats = await categories}
	{#if cats.length}
		<div class="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
			{#each cats as category (category.id)}
				<section class="rounded-2xl border border-neutral-200 bg-white p-5">
					<a
						href={`/search?category_id=${category.id}`}
						class="text-base font-bold text-neutral-900 hover:text-brand-800"
					>
						{category.name}
					</a>
					{#if category.children?.length}
						<ul class="mt-3 space-y-1.5">
							{#each category.children as child (child.id)}
								<li>
									<a
										href={`/search?category_id=${child.id}`}
										class="text-sm text-neutral-600 transition hover:text-brand-700"
									>
										{child.name}
									</a>
								</li>
							{/each}
						</ul>
					{/if}
				</section>
			{/each}
		</div>
	{:else}
		<div class="mt-6 rounded-2xl border border-dashed border-neutral-300 bg-white p-10 text-center">
			<p class="text-neutral-600">Ainda não existem categorias.</p>
		</div>
	{/if}

	{#snippet pending()}
		<div class="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
			{#each { length: 6 }, i (i)}
				<div class="h-40 animate-pulse rounded-2xl bg-neutral-200"></div>
			{/each}
		</div>
	{/snippet}

	{#snippet failed()}
		<div class="mt-6 rounded-2xl border border-red-200 bg-red-50 p-6 text-sm text-red-700">
			Não foi possível carregar as categorias.
		</div>
	{/snippet}
</svelte:boundary>
