<script lang="ts">
	import Logo from '$lib/components/Logo.svelte';
</script>

<footer class="mt-16 border-t border-neutral-200 bg-white">
	<div
		class="mx-auto flex max-w-7xl flex-col items-center justify-between gap-4 px-4 py-8 text-sm text-neutral-500 sm:flex-row sm:px-6"
	>
		<Logo class="text-lg" />
		<nav class="flex flex-wrap items-center justify-center gap-x-6 gap-y-2">
			<a href="/" class="hover:text-brand-800">Início</a>
			<a href="/search" class="hover:text-brand-800">Explorar</a>
			<a href="/categories" class="hover:text-brand-800">Categorias</a>
			<a href="/anunciar" class="hover:text-brand-800">Anunciar</a>
		</nav>
		<p>© {new Date().getFullYear()} Despacha Aí · Benguela, Angola</p>
	</div>
</footer>
