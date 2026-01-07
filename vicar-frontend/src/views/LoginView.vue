<script setup lang="ts">
import { computed, reactive, ref } from "vue"
import VCard from "@/components/ui/VCard.vue"
import VInput from "@/components/ui/VInput.vue"
import VButton from "@/components/ui/VButton.vue"
import {useToast} from "@/composables/useToast.ts";
import {login, register} from "@/rest/api/auth.ts";
import {useConfirm} from "@/composables/useConfirm.ts";
import {useRoute} from "vue-router";
import router from "@/router.ts";

const toast = useToast();
const confirm = useConfirm();
const route = useRoute();

const username = ref<string | null>("")
const password = ref<string | null>("")
const loading = ref(false)

const payload = reactive<{
  username?: string
  password?: string
}>({})

const canSubmit = computed(() => {
  return Boolean((username.value ?? "").trim()) && Boolean((password.value ?? "").trim()) && !loading.value
})

function validate() {
  payload.username = undefined
  payload.password = undefined

  const b = (username.value ?? "").trim()
  const p = (password.value ?? "").trim()

  if (!b) payload.username = "Bitte Benutzernamen eingeben."
  if (!p) payload.password = "Bitte Passwort eingeben."

  return !payload.username && !payload.password
}

async function doLogin() {
  if (!validate()) return

  loading.value = true

  const trimmedUsername = (username.value ?? "").trim();
  const trimmedPassword = (password.value ?? "").trim();

  try {
    const res = await login(trimmedUsername, trimmedPassword);
    if (typeof res === 'string') {
      if (res === 'USER_NOT_FOUND') {
        const ok = await confirm({
          title: 'Benutzer nicht gefunden',
          message: 'Möchtest du einen neuen Account erstellen und dich damit anmelden?',
          yesText: 'Account erstellen',
          yesVariante: 'primary',
          noText: 'Abbrechen',
        });
        if (!ok) {
          return;
        }

        if (!await register(trimmedUsername, trimmedPassword)) {
          toast.error('Registrierung fehlgeschlagen!');
          return;
        }

        const newLoginRes = await login(trimmedUsername, trimmedPassword);
        if (newLoginRes !== 'SUCCESS') {
          toast.error('Anmeldung nach Registrierung fehlgeschlagen!');
          return;
        }

        await redirect();
        return;
      }

      if (res === 'SUCCESS') {
        await redirect();
        return;
      }

      if (res === 'INVALID_CREDS') {
        toast.error('Anmeldung fehlgeschlagen: Ungültiger Benutzername oder Passwort.');
        return;
      }
    }

    toast.error('Anmeldung fehlgeschlagen: Unbekannter Fehler.');
  } catch {
    toast.error('Anmeldung fehlgeschlagen: Ein unerwarteter Fehler ist aufgetreten.');
  } finally {
    loading.value = false
  }
}

async function redirect() {
  if (route.query.redirect) {
    await router.push(atob(route.query.redirect as string));
  } else {
    await router.push('/');
  }
}
</script>

<template>
  <main class="login">
    <div class="smoke" />

    <section class="login__wrap">
      <VCard class="login__card" elevated>
        <header class="login__header">
          <h1 class="login__title">Anmelden</h1>
          <p class="login__subtitle">Bitte melde dich an, um fortzufahren.</p>
        </header>

        <form class="login__form" @submit.prevent="doLogin">
          <VInput
            v-model="username"
            label="Benutzername"
            autocomplete="username"
            inputmode="text"
            spellcheck="false"
            placeholder="z. B. vicar"
            :error="payload.username"
            :disabled="loading"
          />

          <VInput
            v-model="password"
            label="Passwort"
            type="password"
            autocomplete="current-password"
            placeholder="••••••••"
            :error="payload.password"
            :disabled="loading"
          />

          <VButton
            class="login__submit"
            variant="primary"
            size="lg"
            :loading="loading"
            :disabled="!canSubmit"
            nativeType="submit"
            block
          >
            Anmelden
          </VButton>
        </form>
      </VCard>
    </section>
  </main>
</template>

<style scoped lang="scss">
.login {
  min-height: 100dvh;
  display: grid;
  place-items: center;
  padding: 1.25rem;
}

.login__wrap {
  width: min(520px, 100%);
  display: grid;
  gap: 1rem;
}

.login__card {
  padding: clamp(1rem, 3vw, 1.5rem);
}

.login__header {
  display: grid;
  gap: 0.35rem;
  margin-bottom: 1rem;
}

.login__title {
  margin: 0;
}

.login__subtitle {
  margin: 0;
  color: rgba(255, 255, 255, 0.72);
}

.login__form {
  display: grid;
  gap: 0.9rem;
}

.login__submit {
  margin-top: 0.25rem;
}

@media (max-width: 420px) {
  .login {
    padding: 0.9rem;
  }
}
</style>
