<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getAdminUsers, setUserRole } from '$lib/remote/users.remote';
	import { getCurrentUser } from '$lib/remote/auth.remote';
	import type { UserRole } from '$lib/types';

	const me = getCurrentUser();

	const search = $derived(page.url.searchParams.get('search') ?? '');
	const pageNumber = $derived(Number(page.url.searchParams.get('page')) || 1);

	const users = $derived(
		getAdminUsers({ search: search || undefined, page: pageNumber, limit: 20 })
	);

	let busyId = $state<string | null>(null);
	let error = $state<string | null>(null);

	async function changeRole(id: string, role: UserRole) {
		busyId = id;
		error = null;
		try {
			await setUserRole({ id, role }).updates(users);
		} catch (err) {
			error =
				err instanceof Error && err.message.includes('own role')
					? 'Não podes alterar o teu próprio papel.'
					: 'Não foi possível alterar o papel.';
		} finally {
			busyId = null;
		}
	}

	function handleSearch(event: SubmitEvent) {
		event.preventDefault();
		const form = event.currentTarget as HTMLFormElement;
		const value = new FormData(form).get('search')?.toString().trim() ?? '';
		goto(value ? `/admin?search=${encodeURIComponent(value)}` : '/admin', {
			keepFocus: true
		});
	}

	function pageUrl(target: number) {
		const params = new URLSearchParams();
		if (search) params.set('search', search);
		if (target > 1) params.set('page', String(target));
		const qs = params.toString();
		return qs ? `/admin?${qs}` : '/admin';
	}

	const ROLE_LABELS: Record<UserRole, string> = {
		user: 'Utilizador',
		moderator: 'Moderador',
		admin: 'Admin'
	};

	const roleStyles: Record<string, string> = {
		admin: 'bg-brand-100 text-brand-800',
		moderator: 'bg-sky-100 text-sky-800',
		user: 'bg-neutral-100 text-neutral-600'
	};

	function initials(name: string) {
		return name
			.split(/\s+/)
			.slice(0, 2)
			.map((part) => part[0]?.toUpperCase() ?? '')
			.join('');
	}
</script>

<svelte:head>
	<title>Painel admin | Despacha Aí</title>
</svelte:head>

<svelte:boundary>
	{#if (await me)?.role !== 'admin'}
		<div class="mx-auto mt-12 max-w-md rounded-2xl border border-red-200 bg-red-50 p-8 text-center">
			<p class="text-lg font-bold text-red-700">Acesso restrito</p>
			<p class="mt-1 text-sm text-red-600">Esta área é exclusiva para administradores.</p>
			<a
				href="/"
				class="mt-4 inline-block rounded-full bg-neutral-900 px-6 py-2.5 text-sm font-semibold text-white"
				>Voltar ao início</a
			>
		</div>
	{:else}
		{@const result = await users}
		{@const myId = (await me)?.id}

		<div class="flex flex-wrap items-end justify-between gap-4">
			<div>
				<h1 class="text-2xl font-bold text-neutral-900">Painel admin</h1>
				<p class="mt-1 text-sm text-neutral-500">
					{result.total}
					{result.total === 1 ? 'utilizador registado' : 'utilizadores registados'}
				</p>
			</div>

			<form onsubmit={handleSearch} class="flex gap-2" role="search">
				<input
					type="search"
					name="search"
					value={search}
					placeholder="Pesquisar por nome ou telefone"
					class="h-10 w-64 rounded-full border-neutral-200 bg-white px-4 text-sm focus:border-brand-500 focus:ring-brand-500"
				/>
				<button
					type="submit"
					class="h-10 rounded-full bg-brand-700 px-4 text-sm font-semibold text-white transition hover:bg-brand-800"
					>Pesquisar</button
				>
			</form>
		</div>

		{#if error}
			<p class="mt-4 rounded-xl bg-red-50 p-3 text-sm text-red-700">{error}</p>
		{/if}

		<div class="mt-6 overflow-hidden rounded-2xl border border-neutral-200 bg-white">
			{#each result.data as user (user.id)}
				<div
					class="flex flex-wrap items-center gap-3 border-b border-neutral-100 p-4 last:border-b-0"
				>
					{#if user.avatar_url}
						<img src={user.avatar_url} alt="" class="h-10 w-10 rounded-full object-cover" />
					{:else}
						<span
							class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-brand-100 text-sm font-bold text-brand-800"
							>{initials(user.name)}</span
						>
					{/if}

					<div class="min-w-0 flex-1">
						<p class="flex items-center gap-2 truncate font-semibold text-neutral-900">
							{user.name}
							{#if user.id === myId}
								<span
									class="rounded-full bg-neutral-100 px-2 py-0.5 text-[10px] font-bold text-neutral-500"
									>tu</span
								>
							{/if}
						</p>
						<p class="truncate text-sm text-neutral-500">
							{user.phone}{user.email ? ` · ${user.email}` : ''}
						</p>
					</div>

					<span
						class="rounded-full px-2.5 py-1 text-xs font-bold {roleStyles[user.role] ??
							roleStyles.user}"
					>
						{ROLE_LABELS[user.role] ?? user.role}
					</span>

					{#if user.status !== 'active'}
						<span class="rounded-full bg-red-100 px-2.5 py-1 text-xs font-bold text-red-700">
							{user.status}
						</span>
					{/if}

					<select
						value={user.role}
						disabled={busyId === user.id || user.id === myId}
						onchange={(event) =>
							changeRole(user.id, (event.currentTarget as HTMLSelectElement).value as UserRole)}
						class="h-9 rounded-lg border-neutral-200 text-sm focus:border-brand-500 focus:ring-brand-500 disabled:opacity-50"
					>
						<option value="user">Utilizador</option>
						<option value="moderator">Moderador</option>
						<option value="admin">Admin</option>
					</select>
				</div>
			{:else}
				<p class="p-8 text-center text-sm text-neutral-500">Nenhum utilizador encontrado.</p>
			{/each}
		</div>

		{#if result.total_pages > 1}
			<nav class="mt-6 flex items-center justify-center gap-2" aria-label="Paginação">
				<a
					href={pageUrl(result.page - 1)}
					class="rounded-full border border-neutral-200 bg-white px-4 py-2 text-sm font-medium text-neutral-700 transition hover:border-brand-300 aria-disabled:pointer-events-none aria-disabled:opacity-40"
					aria-disabled={result.page <= 1}>Anterior</a
				>
				<span class="px-3 text-sm text-neutral-500"
					>Página {result.page} de {result.total_pages}</span
				>
				<a
					href={pageUrl(result.page + 1)}
					class="rounded-full border border-neutral-200 bg-white px-4 py-2 text-sm font-medium text-neutral-700 transition hover:border-brand-300 aria-disabled:pointer-events-none aria-disabled:opacity-40"
					aria-disabled={result.page >= result.total_pages}>Seguinte</a
				>
			</nav>
		{/if}
	{/if}

	{#snippet pending()}
		<div class="h-8 w-48 animate-pulse rounded bg-neutral-200"></div>
		<div class="mt-6 space-y-px overflow-hidden rounded-2xl border border-neutral-200">
			{#each { length: 5 }, i (i)}
				<div class="h-18 animate-pulse bg-neutral-100"></div>
			{/each}
		</div>
	{/snippet}

	{#snippet failed()}
		<div class="mt-6 rounded-2xl border border-red-200 bg-red-50 p-6 text-sm text-red-700">
			Não foi possível carregar os utilizadores.
		</div>
	{/snippet}
</svelte:boundary>
