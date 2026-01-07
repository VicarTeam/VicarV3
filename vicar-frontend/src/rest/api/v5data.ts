import { GET } from '@/rest'
import type { V5Book, V5Clan, V5PredatorType, V5Discipline, V5BloodRitual, V5OblivionCeremony, V5TraitPack } from '@/@types/v5'

export async function getV5Books(): Promise<V5Book[]> {
  return await GET<V5Book[]>('/v5data/books', undefined, [])
}

export async function getV5Clans(): Promise<V5Clan[]> {
  return await GET<V5Clan[]>('/v5data/clans', undefined, [])
}

export async function getV5PredatorTypes(): Promise<V5PredatorType[]> {
  return await GET<V5PredatorType[]>('/v5data/predator-types', undefined, [])
}

export async function getV5Disciplines(): Promise<V5Discipline[]> {
  return await GET<V5Discipline[]>('/v5data/disciplines', undefined, [])
}

export async function getV5BloodRituals(): Promise<V5BloodRitual[]> {
  return await GET<V5BloodRitual[]>('/v5data/blood-rituals', undefined, [])
}

export async function getV5OblivionCeremonies(): Promise<V5OblivionCeremony[]> {
  return await GET<V5OblivionCeremony[]>('/v5data/oblivion-ceremonies', undefined, [])
}

export async function getV5TraitPacks(): Promise<V5TraitPack[]> {
  return await GET<V5TraitPack[]>('/v5data/trait-packs', undefined, [])
}
