<script lang="ts">
	import { page as pageState } from '$app/state';

	let { page, totalPages }: { page: number; totalPages: number } = $props();

	function pageUrl(target: number) {
		const url = new URL(pageState.url);
		if (target <= 1) url.searchParams.delete('page');
		else url.searchParams.set('page', String(target));
		return `${url.pathname}${url.search}`;
	}
</script>

{#if totalPages > 1}
	<nav class="mt-8 flex items-center justify-center gap-2" aria-label="Paginação">
		<a
			href={pageUrl(page - 1)}
			class="rounded-full border border-neutral-200 bg-white px-4 py-2 text-sm font-medium text-neutral-700 transition hover:border-brand-300 aria-disabled:pointer-events-none aria-disabled:opacity-40"
			aria-disabled={page <= 1}
		>
			Anterior
		</a>
		<span class="px-3 text-sm text-neutral-500">Página {page} de {totalPages}</span>
		<a
			href={pageUrl(page + 1)}
			class="rounded-full border border-neutral-200 bg-white px-4 py-2 text-sm font-medium text-neutral-700 transition hover:border-brand-300 aria-disabled:pointer-events-none aria-disabled:opacity-40"
			aria-disabled={page >= totalPages}
		>
			Seguinte
		</a>
	</nav>
{/if}
