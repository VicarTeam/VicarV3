import { defineStore } from 'pinia'
import type { V5Character, V5CharacterCreateRequest, V5CharacterUpdateRequest } from '@/@types/v5'
import { getCharacters, getCharacter, createCharacter, updateCharacter, deleteCharacter } from '@/rest/api/characters'
import type { PaginationResult } from '@/rest'

export const useCharactersStore = defineStore('characters', {
  state: () => ({
    characters: [] as V5Character[],
    currentCharacter: null as V5Character | null,
    total: 0,
    loading: false,
  }),

  actions: {
    async fetchCharacters(search?: string, page?: number) {
      this.loading = true
      try {
        const result = await getCharacters({ search, page, itemsPerPage: 20 })
        this.characters = result.items
        this.total = result.total
      } finally {
        this.loading = false
      }
    },

    async fetchCharacter(id: string) {
      this.loading = true
      try {
        const char = await getCharacter(id)
        if (char) {
          this.currentCharacter = char
          const idx = this.characters.findIndex(c => c.id === id)
          if (idx !== -1) {
            this.characters[idx] = char
          }
        }
        return char
      } finally {
        this.loading = false
      }
    },

    async createCharacter(req: V5CharacterCreateRequest) {
      const result = await createCharacter(req)
      if (result) {
        const char = await getCharacter(result.id)
        if (char) {
          this.currentCharacter = char
          this.characters.unshift(char)
          this.total++
          return char
        }
      }
      return null
    },

    async updateCharacter(id: string, req: V5CharacterUpdateRequest) {
      const char = await updateCharacter(id, req)
      if (char) {
        this.currentCharacter = char
        const idx = this.characters.findIndex(c => c.id === id)
        if (idx !== -1) {
          this.characters[idx] = char
        }
      }
      return char
    },

    async deleteCharacter(id: string) {
      const success = await deleteCharacter(id)
      if (success) {
        this.characters = this.characters.filter(c => c.id !== id)
        this.total--
        if (this.currentCharacter?.id === id) {
          this.currentCharacter = null
        }
      }
      return success
    },
  },

  getters: {
    getCharacterById: (state) => (id: string) => {
      return state.characters.find(c => c.id === id)
    },
  }
})
