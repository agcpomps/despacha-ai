<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getAdminUsers, setUserRole, resetUserPassword } from '$lib/remote/users.remote';
	import { getCurrentUser } from '$lib/remote/auth.remote';
	import type { User, UserRole } from '$lib/types';

	const me = getCurrentUser();

	const search = $derived(page.url.searchParams.get('search') ?? '');
	const pageNumber = $derived(Number(page.url.searchParams.get('page')) || 1);

	const users = $derived(
		getAdminUsers({ search: search || undefined, page: pageNumber, limit: 20 })
	);

	let busyId = $state<string | null>(null);
	let error = $state<string | null>(null);

	// resultado da reposição de password (mostrado em modal)
	let resetResult = $state<{ user: User; password: string } | null>(null);
	let copied = $state(false);

	async function resetPassword(user: User) {
		busyId = user.id;
		error = null;
		try {
			const { password } = await resetUserPassword(user.id);
			resetResult = { user, password };
			copied = false;
		} catch {
			error = 'Não foi possível repor a palavra-passe.';
		} finally {
			busyId = null;
		}
	}

	async function copyPassword() {
		if (!resetResult) return;
		try {
			await navigator.clipboard.writeText(resetResult.password);
			copied = true;
			setTimeout(() => (copied = false), 2000);
		} catch {
			// área de transferência indisponível — a password está visível na mesma
		}
	}

	function whatsappReset() {
		if (!resetResult) return '';
		const phone = resetResult.user.phone.replace(/\D/g, '');
		const text = encodeURIComponent(
			`Olá ${resetResult.user.name}! A tua nova palavra-passe no Despacha Aí é: ${resetResult.password}\n\nEntra e altera-a assim que puderes.`
		);
		return `https://wa.me/${phone}?text=${text}`;
	}

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

					<button
						type="button"
						onclick={() => resetPassword(user)}
						disabled={busyId === user.id}
						class="rounded-full border border-neutral-200 px-3 py-1.5 text-xs font-semibold text-neutral-600 transition hover:border-brand-400 hover:text-brand-700 disabled:opacity-50"
						>Repor senha</button
					>

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

{#if resetResult}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
		<button
			type="button"
			class="absolute inset-0 cursor-default bg-neutral-900/50 backdrop-blur-sm"
			aria-label="Fechar"
			onclick={() => (resetResult = null)}
		></button>

		<div
			role="dialog"
			aria-modal="true"
			class="relative w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl shadow-neutral-900/20"
		>
			<h2 class="text-lg font-bold text-neutral-900">Palavra-passe reposta</h2>
			<p class="mt-1 text-sm text-neutral-600">
				Nova palavra-passe temporária de <strong>{resetResult.user.name}</strong>. Envia-lha e pede
				para a alterar depois de entrar.
			</p>

			<div
				class="mt-4 flex items-center justify-between gap-3 rounded-xl border border-neutral-200 bg-neutral-50 px-4 py-3"
			>
				<code class="text-lg font-bold tracking-wider text-brand-800">{resetResult.password}</code>
				<button
					type="button"
					onclick={copyPassword}
					class="shrink-0 rounded-full border border-neutral-200 bg-white px-3 py-1.5 text-xs font-semibold text-neutral-600 transition hover:border-brand-400"
				>
					{copied ? 'Copiado!' : 'Copiar'}
				</button>
			</div>

			<div class="mt-5 flex justify-end gap-3">
				<button
					type="button"
					onclick={() => (resetResult = null)}
					class="h-10 rounded-full border border-neutral-200 px-5 text-sm font-medium text-neutral-700 transition hover:border-neutral-300"
					>Fechar</button
				>
				<a
					href={whatsappReset()}
					target="_blank"
					rel="noopener noreferrer"
					class="flex h-10 items-center gap-2 rounded-full bg-whatsapp px-5 text-sm font-bold text-white transition hover:brightness-95"
				>
					<svg class="h-4 w-4" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
						<path
							d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347Z"
						/>
					</svg>
					Enviar por WhatsApp
				</a>
			</div>
		</div>
	</div>
{/if}
