import { GET, GET_AND_THROW, GET_PAGINATED, POST, POST_AND_THROW, PATCH, DELETE, type PaginationOptions, type PaginationResult } from '@/rest'
import type { V5Character, V5CharacterCreateRequest, V5CharacterUpdateRequest } from '@/@types/v5'

export async function getCharacters(opts?: PaginationOptions): Promise<PaginationResult<V5Character>> {
  return await GET_PAGINATED<V5Character>('/characters', opts)
}

export async function getCharacter(id: string): Promise<V5Character | null> {
  return await GET<V5Character>(`/characters/${id}`)
}

export async function createCharacter(req: V5CharacterCreateRequest): Promise<V5Character | null> {
  return await POST_AND_THROW<V5Character>('/characters', req)
}

export async function updateCharacter(id: string, req: V5CharacterUpdateRequest): Promise<V5Character | null> {
  return await PATCH<V5Character>(`/characters/${id}`, req)
}

export async function deleteCharacter(id: string): Promise<boolean> {
  const res = await DELETE(`/characters/${id}`)
  return !!res
}

export async function updateCharacterAttribute(characterId: string, attrId: string, value: number): Promise<boolean> {
  const res = await PATCH(`/characters/${characterId}/attributes/${attrId}`, { value })
  return !!res
}

export async function updateCharacterSkill(characterId: string, skillId: string, value: number, specializations?: string[]): Promise<boolean> {
  const res = await PATCH(`/characters/${characterId}/skills/${skillId}`, { value, specializations })
  return !!res
}

export async function addCharacterDiscipline(characterId: string, disciplineID: string): Promise<boolean> {
  const res = await POST(`/characters/${characterId}/disciplines`, { disciplineID })
  return !!res
}

export async function updateCharacterDiscipline(characterId: string, disciplineId: string, data: { points?: number; currentLevel?: number }): Promise<boolean> {
  const res = await PATCH(`/characters/${characterId}/disciplines/${disciplineId}`, data)
  return !!res
}

export async function deleteCharacterDiscipline(characterId: string, disciplineId: string): Promise<boolean> {
  const res = await DELETE(`/characters/${characterId}/disciplines/${disciplineId}`)
  return !!res
}
