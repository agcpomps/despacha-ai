<script lang="ts">
	import { changePassword } from '$lib/remote/auth.remote';
	import { getCurrentUser } from '$lib/remote/auth.remote';
	import Seo from '$lib/components/Seo.svelte';

	const me = getCurrentUser();
</script>

<Seo title="A minha conta | Despacha Aí" description="Gere a tua conta no Despacha Aí." />

<div class="mx-auto max-w-md">
	<h1 class="text-2xl font-bold text-neutral-900">A minha conta</h1>

	<svelte:boundary>
		{@const user = await me}
		{#if user}
			<div class="mt-4 rounded-2xl border border-neutral-200 bg-white p-6">
				<p class="font-semibold text-neutral-900">{user.name}</p>
				<p class="mt-0.5 text-sm text-neutral-500">
					{user.phone}{user.email ? ` · ${user.email}` : ''}
				</p>
			</div>
		{/if}
		{#snippet pending()}
			<div class="mt-4 h-20 animate-pulse rounded-2xl bg-neutral-200"></div>
		{/snippet}
		{#snippet failed()}{/snippet}
	</svelte:boundary>

	<section class="mt-6 rounded-2xl border border-neutral-200 bg-white p-6">
		<h2 class="text-base font-bold text-neutral-900">Alterar palavra-passe</h2>
		<p class="mt-1 text-sm text-neutral-500">
			Recomendamos alterar a palavra-passe se recebeste uma temporária do suporte.
		</p>

		{#if changePassword.result?.success}
			<p class="mt-4 rounded-xl bg-brand-50 p-3 text-sm font-medium text-brand-800">
				Palavra-passe alterada com sucesso. ✓
			</p>
		{/if}

		<form {...changePassword} class="mt-4 space-y-4">
			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Palavra-passe actual</span>
				<input
					{...changePassword.fields.current_password.as('password')}
					autocomplete="current-password"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each changePassword.fields.current_password.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Nova palavra-passe</span>
				<input
					{...changePassword.fields.new_password.as('password')}
					autocomplete="new-password"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each changePassword.fields.new_password.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<label class="block">
				<span class="text-sm font-semibold text-neutral-700">Confirmar nova palavra-passe</span>
				<input
					{...changePassword.fields.confirm_password.as('password')}
					autocomplete="new-password"
					class="mt-1 h-12 w-full rounded-xl border-neutral-200 focus:border-brand-500 focus:ring-brand-500"
				/>
				{#each changePassword.fields.confirm_password.issues() ?? [] as issue (issue.message)}
					<p class="mt-1 text-sm text-red-600">{issue.message}</p>
				{/each}
			</label>

			<button
				type="submit"
				disabled={changePassword.pending > 0}
				class="h-12 w-full rounded-full bg-brand-700 text-sm font-bold text-white transition hover:bg-brand-800 disabled:opacity-60"
			>
				{changePassword.pending > 0 ? 'A guardar…' : 'Alterar palavra-passe'}
			</button>
		</form>
	</section>
</div>
