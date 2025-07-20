import { queryOptions } from '@tanstack/react-query';
import {
  GetModArchivesByModId,
  GetModFilesByModId,
  GetModlistById,
  GetModlists,
  GetModsByProfileId,
  GetProfileFilesByProfileId,
  GetProfilesByModlistId,
} from '~/wailsjs/go/main/App';

export const modListsQueryOptions = queryOptions({
  queryKey: ['modlists'],
  queryFn: async () => {
    return await GetModlists();
  },
});

export const modListQueryOptions = (id: string) =>
  queryOptions({
    queryKey: ['modlists', id],
    queryFn: async () => {
      return await GetModlistById(id);
    },
  });

export const profilesQueryOptions = (modlistId: string) =>
  queryOptions({
    queryKey: ['modlists', modlistId, 'profiles'],
    queryFn: async () => {
      return await GetProfilesByModlistId(modlistId);
    },
  });

export const profileFilesQueryOptions = (profileId: string) =>
  queryOptions({
    queryKey: ['profiles', profileId, 'files'],
    queryFn: async () => {
      return await GetProfileFilesByProfileId(profileId);
    },
  });

export const profileModsQueryOptions = (profileId: string) =>
  queryOptions({
    queryKey: ['profiles', profileId, 'mods'],
    queryFn: async () => {
      return await GetModsByProfileId(profileId);
    },
  });

export const modArchivesQueryOptions = (modId: string) =>
  queryOptions({
    queryKey: ['mods', modId, 'archive'],
    queryFn: async () => {
      return await GetModArchivesByModId(modId);
    },
  });

export const modFilesQueryOptions = (modId: string) =>
  queryOptions({
    queryKey: ['mods', modId, 'files'],
    queryFn: async () => {
      return await GetModFilesByModId(modId);
    },
  });
