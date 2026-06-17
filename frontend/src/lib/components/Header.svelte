<script lang="ts">
	import { page } from '$app/state';
	import { getCurrentUser, logout } from '$lib/remote/auth.remote';
	import Logo from '$lib/components/Logo.svelte';

	let menuOpen = $state(false);
	let menuContainer = $state<HTMLElement | null>(null);

	const user = $derived(getCurrentUser());

	function handleWindowClick(event: MouseEvent) {
		if (menuContainer && !menuContainer.contains(event.target as Node)) {
			menuOpen = false;
		}
	}

	function initials(name: string) {
		return name
			.split(/\s+/)
			.slice(0, 2)
			.map((part) => part[0]?.toUpperCase() ?? '')
			.join('');
	}
</script>

<svelte:window onclick={handleWindowClick} />

<header class="sticky top-0 z-30 border-b border-neutral-200 bg-white/90 backdrop-blur">
	<div class="mx-auto flex h-16 max-w-7xl items-center gap-6 px-4 sm:px-6">
		<a href="/" class="flex shrink-0 items-center" aria-label="Despacha Aí — início">
			<Logo class="text-2xl sm:text-3xl" />
		</a>

		<form action="/search" class="hidden flex-1 sm:block" role="search">
			<label class="relative block">
				<span class="sr-only">Pesquisar</span>
				<svg
					class="pointer-events-none absolute top-1/2 left-4 h-4 w-4 -translate-y-1/2 text-neutral-400"
					viewBox="0 0 20 20"
					fill="currentColor"
					aria-hidden="true"
				>
					<path
						fill-rule="evenodd"
						d="M9 3.5a5.5 5.5 0 1 0 0 11 5.5 5.5 0 0 0 0-11ZM2 9a7 7 0 1 1 12.452 4.391l3.328 3.329a.75.75 0 1 1-1.06 1.06l-3.329-3.328A7 7 0 0 1 2 9Z"
						clip-rule="evenodd"
					/>
				</svg>
				<input
					type="search"
					name="search"
					value={page.url.pathname === '/search' ? (page.url.searchParams.get('search') ?? '') : ''}
					placeholder="Pesquisar artigos, marcas e categorias"
					class="h-11 w-full rounded-full border-neutral-200 bg-neutral-100 pl-11 text-sm placeholder:text-neutral-400 focus:border-brand-500 focus:bg-white focus:ring-brand-500"
				/>
			</label>
		</form>

		<nav class="ml-auto flex shrink-0 items-center gap-2">
			<a
				href="/anunciar"
				class="hidden items-center gap-1.5 rounded-full bg-brand-700 px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-brand-800 sm:inline-flex"
			>
				<svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
					<path
						d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z"
					/>
				</svg>
				Anunciar
			</a>

			<svelte:boundary>
				{#if await user}
					{@const me = (await user)!}
					<div class="relative" bind:this={menuContainer}>
						<button
							type="button"
							class="flex items-center gap-2 rounded-full border border-neutral-200 bg-white py-1.5 pr-3 pl-1.5 text-sm font-medium text-neutral-700 transition hover:border-brand-300"
							onclick={() => (menuOpen = !menuOpen)}
							aria-expanded={menuOpen}
							aria-haspopup="menu"
						>
							{#if me.avatar_url}
								<img src={me.avatar_url} alt="" class="h-7 w-7 rounded-full object-cover" />
							{:else}
								<span
									class="flex h-7 w-7 items-center justify-center rounded-full bg-brand-100 text-xs font-bold text-brand-800"
									>{initials(me.name)}</span
								>
							{/if}
							<span class="hidden max-w-28 truncate md:block">{me.name.split(' ')[0]}</span>
							<svg class="h-4 w-4 text-neutral-400" viewBox="0 0 20 20" fill="currentColor">
								<path
									fill-rule="evenodd"
									d="M5.22 8.22a.75.75 0 0 1 1.06 0L10 11.94l3.72-3.72a.75.75 0 1 1 1.06 1.06l-4.25 4.25a.75.75 0 0 1-1.06 0L5.22 9.28a.75.75 0 0 1 0-1.06Z"
									clip-rule="evenodd"
								/>
							</svg>
						</button>

						{#if menuOpen}
							<div
								class="absolute right-0 mt-2 w-56 overflow-hidden rounded-xl border border-neutral-200 bg-white py-1 shadow-lg shadow-neutral-900/10"
								role="menu"
							>
								<div class="border-b border-neutral-100 px-4 py-3">
									<p class="truncate text-sm font-semibold text-neutral-900">{me.name}</p>
									<p class="truncate text-xs text-neutral-500">{me.phone}</p>
								</div>
								<a
									href="/dashboard"
									class="block px-4 py-2.5 text-sm text-neutral-700 hover:bg-neutral-50"
									role="menuitem">Os meus anúncios</a
								>
								<a
									href="/conta"
									class="block px-4 py-2.5 text-sm text-neutral-700 hover:bg-neutral-50"
									role="menuitem">A minha conta</a
								>
								<a
									href="/anunciar"
									class="block px-4 py-2.5 text-sm text-neutral-700 hover:bg-neutral-50 sm:hidden"
									role="menuitem">Anunciar</a
								>
								{#if me.role === 'admin'}
									<a
										href="/admin"
										class="block px-4 py-2.5 text-sm font-semibold text-brand-700 hover:bg-brand-50"
										role="menuitem">Painel admin</a
									>
								{/if}
								<form {...logout}>
									<button
										type="submit"
										class="block w-full px-4 py-2.5 text-left text-sm text-red-600 hover:bg-red-50"
										role="menuitem">Terminar sessão</button
									>
								</form>
							</div>
						{/if}
					</div>
				{:else}
					<a
						href="/login"
						class="rounded-full px-4 py-2.5 text-sm font-semibold text-neutral-700 transition hover:text-brand-800"
						>Entrar</a
					>
					<a
						href="/register"
						class="rounded-full border border-brand-700 px-4 py-2.5 text-sm font-semibold text-brand-800 transition hover:bg-brand-50"
						>Criar conta</a
					>
				{/if}

				{#snippet pending()}
					<div class="h-10 w-24 animate-pulse rounded-full bg-neutral-100"></div>
				{/snippet}

				{#snippet failed()}
					<a
						href="/login"
						class="rounded-full px-4 py-2.5 text-sm font-semibold text-neutral-700 transition hover:text-brand-800"
						>Entrar</a
					>
				{/snippet}
			</svelte:boundary>
		</nav>
	</div>

	<form action="/search" class="border-t border-neutral-100 px-4 py-2 sm:hidden" role="search">
		<label class="relative block">
			<span class="sr-only">Pesquisar</span>
			<input
				type="search"
				name="search"
				placeholder="Pesquisar no Despacha Aí"
				class="h-10 w-full rounded-full border-neutral-200 bg-neutral-100 px-4 text-sm placeholder:text-neutral-400 focus:border-brand-500 focus:bg-white focus:ring-brand-500"
			/>
		</label>
	</form>
</header>
