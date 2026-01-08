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

export type AddViewerRequest = { userID: string }

export async function addCharacterViewer(characterId: string, userId: string): Promise<boolean> {
  const res = await POST(`/characters/${characterId}/viewers`, { userID: userId } satisfies AddViewerRequest)
  return !!res
}

export async function removeCharacterViewer(characterId: string, userId: string): Promise<boolean> {
  const res = await DELETE(`/characters/${characterId}/viewers/${userId}`)
  return !!res
}