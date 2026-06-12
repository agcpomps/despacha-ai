<script lang="ts">
	import { page } from '$app/state';

	const items = [
		{
			href: '/',
			label: 'Início',
			d: 'M11.47 3.841a.75.75 0 0 1 1.06 0l8.69 8.69a.75.75 0 1 0 1.06-1.061l-8.689-8.69a2.25 2.25 0 0 0-3.182 0l-8.69 8.69a.75.75 0 1 0 1.061 1.06l8.69-8.689Z M12 5.432l8.159 8.159c.03.03.06.058.091.086v6.198c0 1.035-.84 1.875-1.875 1.875H15a.75.75 0 0 1-.75-.75v-4.5a.75.75 0 0 0-.75-.75h-3a.75.75 0 0 0-.75.75V21a.75.75 0 0 1-.75.75H5.625a1.875 1.875 0 0 1-1.875-1.875v-6.198a2.29 2.29 0 0 0 .091-.086L12 5.432Z'
		},
		{
			href: '/search',
			label: 'Pesquisar',
			d: 'M10.5 3.75a6.75 6.75 0 1 0 0 13.5 6.75 6.75 0 0 0 0-13.5ZM2.25 10.5a8.25 8.25 0 1 1 14.59 5.28l4.69 4.69a.75.75 0 1 1-1.06 1.06l-4.69-4.69A8.25 8.25 0 0 1 2.25 10.5Z'
		},
		null, // espaço para o botão central de anunciar
		{
			href: '/dashboard',
			label: 'Anúncios',
			d: 'M11.644 1.59a.75.75 0 0 1 .712 0l9.75 5.25a.75.75 0 0 1 0 1.32l-9.75 5.25a.75.75 0 0 1-.712 0l-9.75-5.25a.75.75 0 0 1 0-1.32l9.75-5.25Z M3.265 10.602l7.668 4.129a2.25 2.25 0 0 0 2.134 0l7.668-4.13 1.37.739a.75.75 0 0 1 0 1.32l-9.75 5.25a.75.75 0 0 1-.71 0l-9.75-5.25a.75.75 0 0 1 0-1.32l1.37-.738Z M10.933 19.231l-7.668-4.13-1.37.739a.75.75 0 0 0 0 1.32l9.75 5.25c.221.12.489.12.71 0l9.75-5.25a.75.75 0 0 0 0-1.32l-1.37-.738-7.668 4.13a2.25 2.25 0 0 1-2.134-.001Z'
		},
		{
			href: '/dashboard',
			label: 'Conta',
			d: 'M18.685 19.097A9.723 9.723 0 0 0 21.75 12c0-5.385-4.365-9.75-9.75-9.75S2.25 6.615 2.25 12a9.723 9.723 0 0 0 3.065 7.097A9.716 9.716 0 0 0 12 21.75a9.716 9.716 0 0 0 6.685-2.653Zm-12.54-1.285A7.486 7.486 0 0 1 12 15a7.486 7.486 0 0 1 5.855 2.812A8.224 8.224 0 0 1 12 20.25a8.224 8.224 0 0 1-5.855-2.438ZM15.75 9a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0Z'
		}
	] as const;

	function isActive(href: string) {
		if (href === '/') return page.url.pathname === '/';
		return page.url.pathname.startsWith(href);
	}
</script>

<nav
	class="fixed inset-x-0 bottom-0 z-40 border-t border-neutral-200 bg-white/95 backdrop-blur sm:hidden"
	aria-label="Navegação principal"
>
	<div class="grid grid-cols-5 items-end pb-[env(safe-area-inset-bottom)]">
		{#each items as item, index (index)}
			{#if item}
				<a
					href={item.href}
					class="flex flex-col items-center gap-0.5 py-2 text-[10px] font-semibold transition active:scale-95 {isActive(
						item.href
					) && item.label !== 'Conta'
						? 'text-brand-600'
						: 'text-neutral-400'}"
				>
					<svg class="h-6 w-6" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
						<path fill-rule="evenodd" d={item.d} clip-rule="evenodd" />
					</svg>
					{item.label}
				</a>
			{:else}
				<a
					href="/anunciar"
					class="flex flex-col items-center gap-0.5 py-2 text-[10px] font-bold text-brand-700"
					aria-label="Publicar anúncio"
				>
					<span
						class="-mt-6 flex h-12 w-12 items-center justify-center rounded-full bg-brand-600 text-white shadow-lg shadow-brand-600/40 transition active:scale-95"
					>
						<svg class="h-6 w-6" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
							<path
								d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z"
							/>
						</svg>
					</span>
					Anunciar
				</a>
			{/if}
		{/each}
	</div>
</nav>
