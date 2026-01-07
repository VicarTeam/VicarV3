import { defineStore } from 'pinia'
import type { V5Book, V5Clan, V5PredatorType, V5Discipline, V5BloodRitual, V5OblivionCeremony, V5TraitPack } from '@/@types/v5'
import { getV5Books, getV5Clans, getV5PredatorTypes, getV5Disciplines, getV5BloodRituals, getV5OblivionCeremonies, getV5TraitPacks } from '@/rest/api/v5data'

export const useV5DataStore = defineStore('v5data', {
  state: () => ({
    books: null as V5Book[] | null,
    clans: null as V5Clan[] | null,
    predatorTypes: null as V5PredatorType[] | null,
    disciplines: null as V5Discipline[] | null,
    bloodRituals: null as V5BloodRitual[] | null,
    oblivionCeremonies: null as V5OblivionCeremony[] | null,
    traitPacks: null as V5TraitPack[] | null,
  }),

  actions: {
    async fetchBooks() {
      if (this.books) return this.books
      this.books = await getV5Books()
      return this.books
    },

    async fetchClans() {
      if (this.clans) return this.clans
      this.clans = await getV5Clans()
      return this.clans
    },

    async fetchPredatorTypes() {
      if (this.predatorTypes) return this.predatorTypes
      this.predatorTypes = await getV5PredatorTypes()
      return this.predatorTypes
    },

    async fetchDisciplines() {
      if (this.disciplines) return this.disciplines
      this.disciplines = await getV5Disciplines()
      return this.disciplines
    },

    async fetchBloodRituals() {
      if (this.bloodRituals) return this.bloodRituals
      this.bloodRituals = await getV5BloodRituals()
      return this.bloodRituals
    },

    async fetchOblivionCeremonies() {
      if (this.oblivionCeremonies) return this.oblivionCeremonies
      this.oblivionCeremonies = await getV5OblivionCeremonies()
      return this.oblivionCeremonies
    },

    async fetchTraitPacks() {
      if (this.traitPacks) return this.traitPacks
      this.traitPacks = await getV5TraitPacks()
      return this.traitPacks
    },
  },

  getters: {
    getClanById: (state) => (id: string) => {
      return state.clans?.find(c => c.id === id)
    },

    getPredatorTypeById: (state) => (id: string) => {
      return state.predatorTypes?.find(p => p.id === id)
    },

    getDisciplineById: (state) => (id: string) => {
      return state.disciplines?.find(d => d.id === id)
    },
  }
})
