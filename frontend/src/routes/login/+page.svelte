<script lang="ts">
	import { page } from '$app/state';
	import { login } from '$lib/remote/auth.remote';

	const next = $derived(page.url.searchParams.get('next') ?? '/');
</script>

<svelte:head>
	<title>Entrar | Despacha Aí</title>
</svelte:head>

<div class="mx-auto mt-6 max-w-md">
	<div class="rounded-2xl border border-neutral-200 bg-white p-8">
		<h1 class="text-2xl font-bold text-neutral-900">Bem-vindo de volta</h1>
		<p class="mt-1 text-sm text-neutral-500">Inicia sessão para gerir os teus anúncios.</p>

		<form {...login} class="mt-6 space-y-4">
			<input type="hidden" name="next" value={next} />

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Telemóvel</span>
				<input
					{...login.fields.phone.as('tel')}
					placeholder="923 456 789"
					autocomplete="tel"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each login.fields.phone.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Palavra-passe</span>
				<input
					{...login.fields.password.as('password')}
					autocomplete="current-password"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each login.fields.password.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			{#each login.fields.allIssues() ?? [] as issue (issue.message)}
				{#if issue.path.length === 0}
					<p class="rounded-xl bg-red-50 p-3 text-sm text-red-700">{issue.message}</p>
				{/if}
			{/each}

			<div class="text-right">
				<a href="/recuperar" class="text-sm font-medium text-brand-700 hover:text-brand-800"
					>Esqueci-me da palavra-passe</a
				>
			</div>

			<button
				type="submit"
				disabled={login.pending > 0}
				class="h-12 w-full rounded-full bg-brand-700 text-sm font-bold text-white transition hover:bg-brand-800 disabled:opacity-60"
			>
				{login.pending > 0 ? 'A entrar…' : 'Entrar'}
			</button>
		</form>

		<p class="mt-6 text-center text-sm text-neutral-500">
			Ainda não tens conta?
			<a
				href={`/register${next !== '/' ? `?next=${encodeURIComponent(next)}` : ''}`}
				class="font-semibold text-brand-700 hover:text-brand-800">Criar conta</a
			>
		</p>
	</div>
</div>
